package ai

import (
	"context"
	"ricerise/internal/config"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/samber/do/v2"
)

func NewEmbedder(injector do.Injector) (*dashscope.Embedder, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	embedder, err := dashscope.NewEmbedder(context.Background(), &dashscope.EmbeddingConfig{
		APIKey:  appConfig.EmbeddingToken,
		Model:   appConfig.EmbeddingModel,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		panic("failed to initialize embedder: " + err.Error())
	}

	return embedder, nil
}
