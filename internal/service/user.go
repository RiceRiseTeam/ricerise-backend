package service

import (
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto/request"
	"ricerise/internal/middleware"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserService struct {
	repo      *repository.UserRepository
	auth      *middleware.AuthMiddleware
	appConfig *config.AppConfig
}

func (h UserService) Register(context *gin.Context, request *request.UserRegisterRequest) error {
	newUser := &model.UserModel{
		Username: request.UserName,
		Password: request.Password,
	}
	if err := h.repo.Create(context.Request.Context(), newUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			_ = context.Error(apperror.UserNameConflictError)
			return err
		}

		_ = context.Error(err)
	}

	return nil
}

func (h UserService) Login(context *gin.Context, request *request.UserLoginRequest) (string, error) {
	user, err := h.repo.FindBy(context.Request.Context(), &model.UserModel{
		Username: request.UserID,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = context.Error(apperror.AccountPasswordError)
			return "", err
		}

		_ = context.Error(err)
		return "", err
	}

	accessToken, err := h.auth.CreateAccessToken(user.ID, user.Username, user.PermissionLevel)
	if err != nil {
		_ = context.Error(err)
		return "", err
	}

	refreshToken, err := h.auth.CreateRefreshToken(user.ID, user.Username, user.PermissionLevel)
	if err != nil {
		_ = context.Error(err)
		return "", err
	}

	h.setRefreshToken(context, refreshToken)

	return accessToken, nil
}

func (h UserService) Refresh(context *gin.Context, oldToken string) error {
	userInfo := h.auth.ParseRefreshToken(oldToken)
	if userInfo == nil {
		_ = context.Error(apperror.NoRefreshTokenError)
		return apperror.NoRefreshTokenError
	}

	// TODO(linstarowo): 从数据库验证

	refreshToken, err := h.auth.CreateRefreshToken(userInfo.UserId, userInfo.Username, userInfo.PermissionLevel)
	if err != nil {
		_ = context.Error(apperror.InternalServerError)
		return apperror.InternalServerError
	}

	h.setRefreshToken(context, refreshToken)

	// TODO(linstarowo): 使oldToken失效

	return nil
}

func (h UserService) Logout(context *gin.Context) {

}

func (h UserService) setRefreshToken(context *gin.Context, refreshToken string) {
	var expire int
	if refreshToken == "" {
		expire = -1
	} else {
		expire = h.appConfig.RefreshTokenExpire
	}
	context.SetCookie("refresh_token",
		refreshToken,
		expire,
		"/api/v1/auth/refresh",
		"", // 默认行为 不包括子域名 只有当前域名
		true,
		true,
	)
}

func (h UserService) setAccessToken(context *gin.Context, accessToken string) {
	var expire int
	if accessToken == "" {
		expire = -1
	} else {
		expire = h.appConfig.AccessTokenExpire
	}
	context.SetCookie("access_token",
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
