package middleware

import (
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do/v2"
)

type TokenInfo struct {
	UserId          uint64
	Username        string
	PermissionLevel int
	TokenType       TokenType
}

type PermissionLevel int

type AuthMiddleware struct {
	appConfig *config.AppConfig
	secret    []byte

	/*  实际上 AuthMiddleware 更像是 中间件和Token Manager 的结合体
	这里这样写是因为我不想让Level 的调用是middleware.XXXLevel
	搞个auth 包会更好 但这样就又得拆中间件的功能 就这样吧
	*/
	UserLevel  PermissionLevel
	AdminLevel PermissionLevel
	OwnerLevel PermissionLevel
}

type TokenType string

type TokenClaims struct {
	jwt.RegisteredClaims

	UserId          uint64    `json:"userid"`
	Username        string    `json:"username"`
	PermissionLevel int       `json:"level"`
	TokenType       TokenType `json:"type"`
}

const (
	accessType  TokenType = "access"
	refreshType TokenType = "refresh"
)

func (a AuthMiddleware) CreateHandler(level PermissionLevel) gin.HandlerFunc {
	return func(context *gin.Context) {
		a.handle(context, level)
	}
}

func (a AuthMiddleware) handle(context *gin.Context, expectPermission PermissionLevel) {
	header := context.GetHeader("Authorization")
	if header == "" {
		context.AbortWithStatusJSON(200, dto.Error(apperror.NoAccessTokenError))
		return
	}

	userInfo := a.parseToken(strings.TrimPrefix(header, "Bearer "), accessType)
	if userInfo == nil {
		context.AbortWithStatusJSON(200, dto.Error(apperror.NoAccessTokenError))
		return
	}

	if expectPermission >= 0 && userInfo.PermissionLevel < int(expectPermission) {
		context.AbortWithStatusJSON(200, dto.Error(apperror.NoPermissionError))
		return
	}

	context.Set("auth", userInfo)
}

func (a AuthMiddleware) ParseRefreshToken(token string) *TokenInfo {
	return a.parseToken(token, refreshType)
}

func (a AuthMiddleware) ParseAccessToken(token string) *TokenInfo {
	return a.parseToken(token, accessType)
}

func (a AuthMiddleware) CreateRefreshToken(userId uint64, userName string, permissionLevel int) (string, error) {
	return a.createToken(userId, userName, permissionLevel, refreshType)
}

func (a AuthMiddleware) CreateAccessToken(userId uint64, userName string, permissionLevel int) (string, error) {
	return a.createToken(userId, userName, permissionLevel, accessType)
}

func (a AuthMiddleware) createToken(userId uint64, userName string, permissionLevel int, tokenType TokenType) (string, error) {
	var expiration int
	if tokenType == accessType {
		expiration = a.appConfig.AccessTokenExpire
	} else if tokenType == refreshType {
		expiration = a.appConfig.RefreshTokenExpire
	}
	tokenClaims := TokenClaims{
		UserId:          userId,
		Username:        userName,
		PermissionLevel: permissionLevel,
		TokenType:       tokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ricerise",
			Subject:   userName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiration) * time.Second)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString(a.secret)
}

func (a AuthMiddleware) parseToken(tokenString string, expectType TokenType) *TokenInfo {
	tokenClaims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return a.secret, nil
	}, jwt.WithIssuer("ricerise"), jwt.WithExpirationRequired())

	if err != nil || !token.Valid {
		return nil
	}

	return &TokenInfo{
		UserId:          tokenClaims.UserId,
		Username:        tokenClaims.Username,
		PermissionLevel: tokenClaims.PermissionLevel,
		TokenType:       expectType,
	}
}

// GetUserInfo 实际上只是一个工具方法 传入gin.Context 可能显得权责不分明 但是因为GO 并没有static 方法的概念 就这么写吧
func (a AuthMiddleware) GetUserInfo(context *gin.Context) *TokenInfo {
	authInfo, exists := context.Get("auth")
	if !exists {
		return nil
	}

	return authInfo.(*TokenInfo)
}

func NewAuthMiddleware(injector do.Injector) (*AuthMiddleware, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AuthMiddleware{
		appConfig:  appConfig,
		secret:     []byte(appConfig.JWTSecret),
		UserLevel:  0,
		AdminLevel: 1,
		OwnerLevel: 2,
	}, nil
}
