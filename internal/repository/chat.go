package repository

import (
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ChatMessageRepository struct {
	gorm.Interface[model.ChatMessageModel]
	db *gorm.DB
}

func (u ChatMessageRepository) WithTx(tx *gorm.DB) *ChatMessageRepository {
	return &ChatMessageRepository{
		Interface: gorm.G[model.ChatMessageModel](tx),
		db:        tx,
	}
}

func NewChatMessageRepository(injector do.Injector) (*ChatMessageRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &ChatMessageRepository{
		Interface: gorm.G[model.ChatMessageModel](db),
		db:        db,
	}, nil
}
