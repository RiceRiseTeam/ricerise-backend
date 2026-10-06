package service

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AgentService struct {
	agent *adk.ChatModelAgent
}

func (a AgentService) Chat(ctx *gin.Context, sessionId uint64, message string) {

}

func (a AgentService) GetMessageHistory(ctx *gin.Context, sessionId uint64) {

}

func NewAgentService(injector do.Injector) (*AgentService, error) {
	agentService := &AgentService{}
	chatModel := do.MustInvoke[*openai.ChatModel](injector)
	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "ricerise-agent",
		Description: "ricerise 饭来 Agent",
		Instruction: "请牢记 你是一个",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{},
			},
		},
	})

	if err != nil {
		panic("failed to create chat model agent: " + err.Error())
	}

	agentService.agent = agent

	return agentService, nil
}
