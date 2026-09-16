package request

type UserRegisterRequest struct {
	UserName string `json:"username" binding:"required,min=8,max=64"`
	NickName string `json:"nickname" binding:"required,min=8,max=64"`
	Password string `json:"password" binding:"required,min=16,max=32"`
}
