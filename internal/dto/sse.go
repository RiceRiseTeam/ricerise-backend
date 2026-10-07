package dto

// ToastSSE toast事件结构体
// @Description 服务端创建Toast事件结构体
type ToastSSE struct {
	Message string `json:"message"`
}

// ChatSSE chat事件结构体
// @Description 用户局内聊天结构体
type ChatSSE struct {
	DinnerId uint64   `json:"dinnerId"`
	User     *UserDto `json:"user"`
	Message  string   `json:"message"`
}

type AgentSSE struct {
	SessionId uint64 `json:"sessionId"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func AgentDone(sessionId uint64) *AgentSSE {
	return &AgentSSE{
		SessionId: sessionId,
		Type:      "done",
	}
}

func AgentText(sessionId uint64, message string) *AgentSSE {
	return &AgentSSE{
		SessionId: sessionId,
		Type:      "text",
		Message:   message,
	}
}

func AgentError(sessionId uint64, message string) *AgentSSE {
	return &AgentSSE{
		SessionId: sessionId,
		Type:      "error",
		Message:   message,
	}
}
