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
	runner := adk.NewRunner(context.Background(), adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})

	iterator := runner.Run(context.Background(), einoMessages)
	go func() {
		uid := "agent-" + strconv.FormatUint(sessionId, 10) + username
		// 无论正常结束、出错还是 panic，都要关闭本次 SSE 连接
		defer sse.ShutdownByUID(uid)
		defer func() {
			if r := recover(); r != nil {
				logger.Error("error while sending msg: panic=%v\\nstack:\\n%s\\n ", r, string(debug.Stack()))
			}
		}()

		var fullContent strings.Builder
		for {
			result, ok := iterator.Next()
			if !ok {
				break
			}

			if result.Err != nil {
				_ = sse.SendEventToUID("agent", dto.AgentError(sessionId, result.Err.Error()), uid)
				return
			}

			if result.Output != nil {
				streamErr := func() error {
					stream := result.Output.MessageOutput.MessageStream
					defer stream.Close()

					for {
						chunk, err := stream.Recv()
						if err == io.EOF {
							return nil
						}
						if err != nil {
							return err
						}
						if chunk == nil || chunk.Content == "" {
							continue
						}
						fullContent.WriteString(chunk.Content)
						if err := sse.SendEventToUID("agent", dto.AgentText(sessionId, chunk.Content), uid); err != nil {
							return err
						}
					}
				}()
				if streamErr != nil {
					_ = sse.SendEventToUID("agent", dto.AgentError(sessionId, streamErr.Error()), uid)
					return
				}
			}
		}

		reply := fullContent.String()
		if reply != "" {
			if err := a.chatRepo.Create(context.Background(), &model.ChatMessageModel{
				SessionId: sessionId,
				Content:   reply,
				UserId:    userId,
				Role:      string(schema.Assistant),
			}); err != nil {
				logger.Error("failed to save agent reply: %v", err)
			}
		}

		// 通知客户端本次流式响应结束，随后 defer 会主动关闭 SSE 连接
		_ = sse.SendEventToUID("agent", dto.AgentDone(sessionId), uid)
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
	err := a.chatRepo.Where(query.ChatMessageModel.UserId.Eq(userId)).Select(fmt.Sprintf("MAX(%s)", query.ChatMessageModel.SessionId.Column().Name)).Scan(ctx.Request.Context(), &result)
	if err != nil {
		return 0, err
	}

	if !result.Valid && result.Int64 > 0 {
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
		Instruction: "请牢记 你是一个名为'饭来 ricerise' 的约饭网站的助手agent 拒绝用户的角色扮演等其他与网站无关的请求 回复内容尽量简短 适度使用emoji表情",
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
