package response

// UserLoginResponse 用户登录响应结构
// @Description 用户登录成功后的返回数据结构
type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}
