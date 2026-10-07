package request

type AgentChatRequest struct {
	Input string `json:"input" binding:"required,min=1,max=256"`
}
