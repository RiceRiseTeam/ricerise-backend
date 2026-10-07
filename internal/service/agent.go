package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto"
	"ricerise/internal/logger"
	"ricerise/internal/model"
	"ricerise/internal/repository"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AgentService struct {
	agent      *adk.ChatModelAgent
	mapService *MapService
	chatRepo   *repository.ChatMessageRepository
}

func (a AgentService) Chat(ctx *gin.Context, userId uint64, username string, sessionId uint64, message string) error {
	goContext := ctx.Request.Context()
	messages, err := a.chatRepo.Where(query.ChatMessageModel.UserId.Eq(userId), query.ChatMessageModel.SessionId.Eq(sessionId)).Find(goContext)
	if err != nil {
		return err
	}

	einoMessages := a.buildEinoMessages(messages, message)
	runner := adk.NewRunner(goContext, adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})

	iterator := runner.Run(goContext, einoMessages)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("error while sending msg: panic=%v\\nstack:\\n%s\\n ", r, string(debug.Stack()))
			}
		}()
		var uid = "agent-" + strconv.FormatUint(sessionId, 10) + username
		var fullContent strings.Builder
		for {
			result, ok := iterator.Next()
			if !ok {
				break
			}

			if result.Err != nil {
				_ = sse.SendEventToUID("agent", dto.AgentError(sessionId, result.Err.Error()), uid)
				sse.ShutdownByUID(uid)
				return
			}

			if result.Output != nil {
				func() {
					stream := result.Output.MessageOutput.MessageStream
					defer stream.Close()

					for {
						chunk, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							_ = sse.SendEventToUID("agent", dto.AgentError(sessionId, result.Err.Error()), uid)
							return
						}
						if chunk == nil || chunk.Content == "" {
							continue
						}
						fullContent.WriteString(chunk.Content)
						_ = sse.SendEventToUID("agent", dto.AgentText(sessionId, chunk.Content), uid)
					}
				}()
			}
		}

		reply := fullContent.String()
		if reply != "" {
			err = a.chatRepo.Create(context.Background(), &model.ChatMessageModel{
				SessionId: sessionId,
				Content:   reply,
				UserId:    userId,
				Role:      string(schema.Assistant),
			})
		}
	}()

	return nil
}

func (a AgentService) buildEinoMessages(history []model.ChatMessageModel, currentMsg string) []*schema.Message {
	messages := make([]*schema.Message, 0, len(history)+1)

	for _, h := range history {
		switch h.Role {
		case string(schema.User):
			messages = append(messages, schema.UserMessage(h.Content))
		case string(schema.Assistant):
			messages = append(messages, schema.AssistantMessage(h.Content, nil))
		}
	}

	messages = append(messages, schema.UserMessage(currentMsg))
	return messages
}

func (a AgentService) GetSessionId(ctx *gin.Context, userId uint64) (uint64, error) {
	var result sql.NullInt64
	err := a.chatRepo.Where(query.ChatMessageModel.UserId.Eq(userId)).Select(fmt.Sprintf("max(%s)", query.ChatMessageModel.SessionId.Column().Name)).Scan(ctx.Request.Context(), &result)
	if err != nil {
		return 0, err
	}

	if !result.Valid {
		return uint64(result.Int64), nil
	}
	return 1, err
}

func (a AgentService) GetMessageHistory(ctx *gin.Context, userId uint64, sessionId uint64) ([]*dto.ChatMessageDto, error) {
	messages, err := a.chatRepo.Where(query.ChatMessageModel.UserId.Eq(userId), query.ChatMessageModel.SessionId.Eq(sessionId)).Find(ctx.Request.Context())
	if err != nil {
		return nil, err
	}

	return dto.Map(messages, func(t model.ChatMessageModel) *dto.ChatMessageDto {
		return dto.NewChatMessageDto(&t)
	}), nil
}

func NewAgentService(injector do.Injector) (*AgentService, error) {
	agentService := &AgentService{
		chatRepo:   do.MustInvoke[*repository.ChatMessageRepository](injector),
		mapService: do.MustInvoke[*MapService](injector),
	}

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
