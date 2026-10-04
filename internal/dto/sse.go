package dto

// ToastSSE toast事件结构体
// @Description 服务端创建Toast事件结构体
type ToastSSE struct {
	Message string `json:"message"`
}

// ChatSSE chat事件结构体
// @Description 用户局内聊天结构体
type ChatSSE struct {
	User    *UserDto `json:"user"`
	Message string   `json:"message"`
}
