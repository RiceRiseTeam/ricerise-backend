package dto

// JoinSSE join事件结构体
// @Description 用户加入饭局事件结构体
type JoinSSE struct {
	User UserDto `json:"user"`
}

// LeftSSE left事件结构体
// @Description 用户离开饭局事件结构体
type LeftSSE struct {
	User UserDto `json:"user"`
}

// ChatSSE chat事件结构体
// @Description 用户局内聊天结构体
type ChatSSE struct {
	User    UserDto `json:"user"`
	Message string  `json:"message"`
}
