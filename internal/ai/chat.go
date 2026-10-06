package ai

import (
	"context"
	"ricerise/internal/config"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/samber/do/v2"
)

func NewChatModel(injector do.Injector) (*openai.ChatModel, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey:  appConfig.AIToken,
		Model:   appConfig.AIModel,
		BaseURL: appConfig.AIBaseURL, // 关键：指向 DeepSeek
	})

	if err != nil {
		panic("failed to initialize chat model: " + err.Error())
	}

	return chatModel, nil
}
