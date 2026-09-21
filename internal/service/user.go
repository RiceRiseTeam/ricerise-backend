package service

import (
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto/request"
	"ricerise/internal/middleware"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo      *repository.UserRepository
	auth      *middleware.AuthMiddleware
	appConfig *config.AppConfig
}

func (h UserService) Register(ctx *gin.Context, request *request.UserRegisterRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return err
	}

	newUser := &model.UserModel{
		Username: request.UserName,
		Nickname: request.NickName,
		Email:    request.Email,
		Password: string(hash),
	}
	if err = h.repo.Create(ctx.Request.Context(), newUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperror.UserNameConflictError
		}

		return err
	}

	return nil
}

func (h UserService) Login(ctx *gin.Context, request *request.UserLoginRequest) (string, error) {
	user, err := h.repo.Where(query.UserModel.Username.Eq(request.UserID)).First(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.AccountPasswordError
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", apperror.AccountPasswordError
	}

	accessToken, err := h.auth.CreateAccessToken(user.ID, user.Username, middleware.PermissionLevel(user.PermissionLevel))
	if err != nil {
		return "", err
	}

	refreshToken, err := h.auth.CreateRefreshToken(user.ID, user.Username, middleware.PermissionLevel(user.PermissionLevel))
	if err != nil {
		return "", err
	}

	h.setRefreshToken(ctx, refreshToken)

	return accessToken, nil
}

func (h UserService) Refresh(ctx *gin.Context, oldToken string) error {
	userInfo := h.auth.ParseRefreshToken(oldToken)
	if userInfo == nil {
		return apperror.NoRefreshTokenError
	}

	// TODO(linstarowo): 从数据库验证

	refreshToken, err := h.auth.CreateRefreshToken(userInfo.UserId, userInfo.Username, userInfo.PermissionLevel)
	if err != nil {
		return apperror.InternalServerError
	}

	h.setRefreshToken(ctx, refreshToken)

	// TODO(linstarowo): 使oldToken失效

	return nil
}

func (h UserService) Logout(ctx *gin.Context) {
	h.setRefreshToken(ctx, "")
	h.setAccessToken(ctx, "")
}

func (h UserService) setRefreshToken(ctx *gin.Context, refreshToken string) {
	var expire int
	if refreshToken == "" {
		expire = -1
	} else {
		expire = h.appConfig.RefreshTokenExpire
	}
	ctx.SetCookie("refresh_token",
		refreshToken,
		expire,
		"/api/v1/auth/refresh",
		"", // 默认行为 不包括子域名 只有当前域名
		true,
		true,
	)
}

func (h UserService) setAccessToken(ctx *gin.Context, accessToken string) {
	var expire int
	if accessToken == "" {
		expire = -1
	} else {
		expire = h.appConfig.AccessTokenExpire
	}
	ctx.SetCookie("access_token",
		accessToken,
		expire,
		"/api/v1/sse",
		"", // 默认行为 不包括子域名 只有当前域名
		true,
		true,
	)
}

func NewUserService(injector do.Injector) (*UserService, error) {
	repo := do.MustInvoke[*repository.UserRepository](injector)
	auth := do.MustInvoke[*middleware.AuthMiddleware](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &UserService{
		repo:      repo,
		auth:      auth,
		appConfig: appConfig,
	}, nil
}
