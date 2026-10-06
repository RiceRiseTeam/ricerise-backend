package ai

import (
	"context"
	"ricerise/internal/config"

	"github.com/Wood-Q/Eino-pgvector/retriever"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/samber/do/v2"
)

func NewRetriever(injector do.Injector) (*retriever.Retriever, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	embedder := do.MustInvoke[*dashscope.Embedder](injector)
	newRetriever, err := retriever.NewRetriever(context.Background(), &retriever.RetrieverConfig{
		Host:      appConfig.SQLHost,
		Port:      appConfig.SQLPort,
		User:      appConfig.SQLUser,
		Password:  appConfig.SQLPassword,
		DBName:    "ricerise",
		TableName: "documents",
		Dimension: 1024,
		Embedding: embedder,
	})

	if err != nil {
		panic("failed to create retriever: " + err.Error())
	}

	return newRetriever, nil
}
