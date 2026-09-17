package apperror

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var InternalServerError = &AppError{
	Code:    50000,
	Message: "Internal Server Error",
}

// 401 鉴权问题

var NoRefreshTokenError = &AppError{
	Code:    40101,
	Message: "refresh token 无效或已过期",
}

var NoPermissionError = &AppError{
	Code:    40402,
	Message: "没有权限",
}

var AccountPasswordError = &AppError{
	Code:    40403,
	Message: "账户或密码错误",
}

// 409XX 冲突问题

var UserNameConflictError = &AppError{
	Code:    40901,
	Message: "用户名已存在",
}
