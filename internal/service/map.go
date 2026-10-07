package service

import (
	"context"
	"errors"
	"fmt"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto"
	"ricerise/internal/dto/querydto"
	"ricerise/internal/dto/request"
	"ricerise/internal/logger"
	"ricerise/internal/middleware"
	"ricerise/internal/model"
	"ricerise/internal/repository"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/Wood-Q/Eino-pgvector/indexer"
	"github.com/Wood-Q/Eino-pgvector/retriever"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
	"github.com/restayway/gogis"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type MapService struct {
	appConfig             *config.AppConfig
	locationRepository    *repository.LocationRepository
	commentRepository     *repository.CommentRepository
	participantRepository *repository.ParticipantRepository
	authMiddleware        *middleware.AuthMiddleware
	indexer               *indexer.Indexer
	retriever             *retriever.Retriever
	chatModel             *openai.ChatModel
}

func (m MapService) DeleteComment(ctx *gin.Context, id uint64) error {
	goContext := ctx.Request.Context()
	comment, err := m.fetchComment(goContext, id)
	if err != nil {
		return err
	}

	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return apperror.NoAccessTokenError
	}
	if userInfo.PermissionLevel < m.authMiddleware.OwnerLevel && userInfo.Username == comment.User.Username {
		return apperror.NoPermissionError
	}
	_, err = m.commentRepository.Where(query.CommentModel.ID.Eq(id)).Delete(goContext)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) DeleteLocation(ctx *gin.Context, id uint64) error {
	goContext := ctx.Request.Context()
	_, err := m.fetchLocation(goContext, id)
	if err != nil {
		return err
	}
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return apperror.NoAccessTokenError
	}
	if userInfo.PermissionLevel < m.authMiddleware.OwnerLevel {
		return apperror.NoPermissionError
	}

	_, err = m.locationRepository.Where(query.LocationModel.ID.Eq(id)).Delete(goContext)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) GetLocationDetail(ctx *gin.Context, id uint64) (*dto.LocationDto, error) {
	location, err := m.locationRepository.Where(query.LocationModel.ID.Eq(id)).Preload(query.LocationModel.User.Name(), nil).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return dto.NewLocationDto(&location), nil
}

func (m MapService) GetLocationsInRange(ctx *gin.Context, request request.GetLocationsRequest) ([]*dto.LocationDto, error) {
	result, err := m.locationRepository.FindAllInRange(ctx.Request.Context(), request.MinLng, request.MinLat, request.MaxLng, request.MaxLat)
	if err != nil {
		return nil, err
	}

	return dto.Map(result, func(t model.LocationModel) *dto.LocationDto {
		return dto.NewLocationDto(&t)
	}), nil
}

func (m MapService) fetchLocation(ctx context.Context, id uint64) (*model.LocationModel, error) {
	location, err := m.locationRepository.Where(query.LocationModel.ID.Eq(id)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return &location, nil
}

func (m MapService) fetchComment(ctx context.Context, id uint64) (*model.CommentModel, error) {
	comment, err := m.commentRepository.Where(query.LocationModel.ID.Eq(id)).Preload(query.CommentModel.User.Name(), nil).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return &comment, nil
}

func (m MapService) SearchLocation(ctx *gin.Context, input string) ([]model.LocationModel, error) {
	goContext := ctx.Request.Context()
	results, err := m.retriever.Retrieve(goContext, input, &retriever.SearchOptions{Limit: 10})
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(results))
	for _, document := range results {
		id, err := strconv.ParseUint(document.ID, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []model.LocationModel{}, nil
	}

	locations, err := m.locationRepository.Where("id IN (?)", ids).Find(goContext)
	if err != nil {
		return nil, err
	}

	indexMap := make(map[uint64]model.LocationModel, len(locations))
	for _, loc := range locations {
		indexMap[loc.ID] = loc
	}
	orderedLocations := make([]model.LocationModel, 0, len(ids))
	for _, id := range ids {
		if d, ok := indexMap[id]; ok {
			orderedLocations = append(orderedLocations, d)
		}
	}

	return orderedLocations, nil
}

func (m MapService) updateDocument(locationId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("error while updating document: panic=%v\\nstack:\\n%s\\n ", r, string(debug.Stack()))
		}
	}()

	goContext := context.Background()
	location, err := m.locationRepository.Where(query.LocationModel.ID.Eq(locationId)).First(goContext)
	if err != nil {
		logger.Error("error when updating document" + err.Error())
		return
	}

	comments, err := m.commentRepository.Where(query.CommentModel.ID.Eq(locationId)).Limit(10).Find(goContext)
	if err != nil {
		logger.Error("error when updating document" + err.Error())
		return
	}

	var summary string
	if len(comments) > 0 {
		summary, err = m.generateSummary(goContext, comments)
		if err != nil {
			logger.Error("error when updating document" + err.Error())
			return
		}
	}

	err = m.saveEmbeddingDocument(goContext, &location, summary)
	if err != nil {
		logger.Error("error when updating document" + err.Error())
	}
}

func (m MapService) saveEmbeddingDocument(ctx context.Context, location *model.LocationModel, summary string) error {
	var builder strings.Builder
	builder.WriteString("[饭店名称]: ")
	builder.WriteString(location.Name)
	builder.WriteString("\n[饭店地址]: ")
	builder.WriteString(location.Address)
	builder.WriteString("\n[饭店描述]: ")
	builder.WriteString(location.Description)
	builder.WriteString("\n[评论总结]: ")
	if summary != "" {
		builder.WriteString(summary)
	} else {
		builder.WriteString("暂无评论")
	}

	content := builder.String()
	docs := []*schema.Document{
		{
			ID:      strconv.FormatUint(location.ID, 10),
			Content: content,
		},
	}

	_, err := m.indexer.Store(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

func (m MapService) generateSummary(ctx context.Context, comments []model.CommentModel) (string, error) {
	systemPrompt := `你是一个专业的餐饮评论分析助手。你的任务是把顾客评论总结成一段结构化的中文摘要。
		输出要求：
		- 长度控制在 150 字以内
		- 重点覆盖：口味、环境、服务、性价比、适合场景（如约会/聚餐/家庭）
		- 如果评论中反复提到某道菜，要特别指出
		- 只输出摘要正文，不要任何额外解释或前缀`

	var builder strings.Builder
	for _, comment := range comments {
		builder.WriteString(comment.Content)
	}

	userPrompt := fmt.Sprintf(`请把下面这家饭店的顾客评论总结成一段中文描述：%s`, builder.String())

	msg, err := m.chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msg.Content), nil
}

func (m MapService) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) (*dto.CommentDto, error) {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	id, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, apperror.AccessNoFoundError
	}
	goContext := ctx.Request.Context()
	location, err := m.fetchLocation(ctx, id)
	if location == nil {
		return nil, err
	}

	participant, err := m.participantRepository.
		Where(query.ParticipantModel.DinnerId.Eq(request.DinnerId)).
		Where(query.ParticipantModel.UserId.Eq(userInfo.UserId)).Preload(query.ParticipantModel.Dinner.Name(), nil).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}

	if participant.Dinner.ID != request.DinnerId {
		return nil, apperror.NoPermissionError
	}

	newComment := &model.CommentModel{
		Rating:     request.Rating,
		Content:    request.Content,
		UserId:     userInfo.UserId,
		LocationId: location.ID,
		DinnerId:   &request.DinnerId,
	}
	err = m.commentRepository.Create(goContext, newComment)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.NameConflictError
		}
		return nil, err
	}

	go m.updateDocument(location.ID)

	return dto.NewCommentDto(newComment), nil
}

func (m MapService) UploadLocation(ctx *gin.Context, request request.UploadLocationRequest) (*dto.LocationDto, error) {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}

	newLocation := &model.LocationModel{
		Location: gogis.Point{
			Lng: request.Longitude,
			Lat: request.Latitude,
		},
		Name:        request.Name,
		Description: request.Description,
		UserId:      userInfo.UserId,
		Address:     request.Address,
	}

	err := m.locationRepository.Create(ctx.Request.Context(), newLocation)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.NameConflictError
		}
		return nil, err
	}
	go m.updateDocument(newLocation.ID)
	return dto.NewLocationDto(newLocation), nil
}

func (m MapService) GetComments(ctx *gin.Context, id uint64, request querydto.MapLocationCommentsQuery) ([]*dto.CommentDto, error) {
	id, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, apperror.AccessNoFoundError
	}

	sql := m.commentRepository.Where(query.CommentModel.LocationId.Eq(id)).Preload(query.CommentModel.User.Name(), nil)

	if request.OrderedBy == "time" {
		if request.StartId != nil {
			sql = sql.Where(query.CommentModel.ID.Lt(*request.StartId), query.CommentModel.CreatedAt.Lt(*request.StartTime))
		}
		sql = sql.Order(query.CommentModel.CreatedAt.Desc()).Order(query.CommentModel.ID.Desc())
	} else if request.OrderedBy == "rank" {
		if request.StartId != nil {
			sql = sql.Where(query.CommentModel.ID.Lt(*request.StartId), query.CommentModel.Rating.Lt(*request.StartRank))
		}
		sql = sql.Order(query.CommentModel.Rating.Desc()).Order(query.CommentModel.ID.Desc())
	}

	result, err := sql.Limit(request.PageSize).Find(ctx.Request.Context())
	return dto.Map(result, func(t model.CommentModel) *dto.CommentDto {
		return dto.NewCommentDto(&t)
	}), err
}

func NewMapService(injector do.Injector) (*MapService, error) {
	return &MapService{
		appConfig:             do.MustInvoke[*config.AppConfig](injector),
		authMiddleware:        do.MustInvoke[*middleware.AuthMiddleware](injector),
		locationRepository:    do.MustInvoke[*repository.LocationRepository](injector),
		commentRepository:     do.MustInvoke[*repository.CommentRepository](injector),
		participantRepository: do.MustInvoke[*repository.ParticipantRepository](injector),
		indexer:               do.MustInvoke[*indexer.Indexer](injector),
		retriever:             do.MustInvoke[*retriever.Retriever](injector),
		chatModel:             do.MustInvoke[*openai.ChatModel](injector),
	}, nil
}
