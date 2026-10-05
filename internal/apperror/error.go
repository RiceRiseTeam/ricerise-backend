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

var ValidationError = &AppError{
	Code:    40000,
	Message: "参数校验失败",
}

// 401 鉴权问题

var NoAccessTokenError = &AppError{
	Code:    40100,
	Message: "未登录或登录已过期",
}

var NoRefreshTokenError = &AppError{
	Code:    40101,
	Message: "refresh token 无效或已过期",
}

var NoPermissionError = &AppError{
	Code:    40102,
	Message: "没有权限",
}

var AccountPasswordError = &AppError{
	Code:    40103,
	Message: "账户或密码错误",
}

var AccessNoFoundError = &AppError{
	Code:    40400,
	Message: "记录不存在",
}

var InviteCodeNotFoundError = &AppError{
	Code:    40401,
	Message: "邀请码不存在或已失效",
}

// 409XX 冲突问题

var NameConflictError = &AppError{
	Code:    40900,
	Message: "名称已存在",
}

var UserNameConflictError = &AppError{
	Code:    40901,
	Message: "用户名已存在",
}

var DinnerConflictError = &AppError{
	Code:    40902,
	Message: "已有进行中的约饭!",
}

var ParticipantConflictError = &AppError{
	Code:    40903,
	Message: "你已经加入饭局了!",
}
