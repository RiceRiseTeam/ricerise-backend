package ai

import (
	"context"
	"ricerise/internal/config"

	"github.com/Wood-Q/Eino-pgvector/indexer"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/samber/do/v2"
)

func NewIndexer(injector do.Injector) (*indexer.Indexer, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	embedder := do.MustInvoke[*dashscope.Embedder](injector)
	myIndexer, err := indexer.NewIndexer(context.Background(), &indexer.IndexerConfig{
		Host:      appConfig.SQLHost,
		Port:      appConfig.SQLPort,
		User:      appConfig.SQLUser,
		Password:  appConfig.SQLPassword,
		DBName:    "ricerise",
		TableName: "documents",
		Dimension: 1024,
		IndexType: "hnsw",
		Embedding: embedder,
	})

	if err != nil {
		panic("failed to create indexer: " + err.Error())
	}

	return myIndexer, nil
}
