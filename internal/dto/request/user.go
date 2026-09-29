package request

// UserRegisterRequest 用户注册请求结构
// @Description 用户注册时提交的数据结构
type UserRegisterRequest struct {
	UserName string `json:"username" binding:"required,alphanum,min=8,max=64"`
	NickName string `json:"nickname" binding:"required,min=8,max=64"`
	Password string `json:"password" binding:"required,min=12,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

// UserLoginRequest 用户登录请求结构
// @Description 用户登录时提交的数据结构
type UserLoginRequest struct {
	UserID   string `json:"user_id" binding:"required,min=1"` // 用户标识符 可以是UserName 或是邮箱
	Password string `json:"password" binding:"required,min=12,max=32"`
}
