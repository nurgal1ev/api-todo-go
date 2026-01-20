package board

import (
	"cli-todo/internal/storage"
	"context"
	"errors"
	"gorm.io/gorm"
)

type CreateBoardData struct {
	Name string `json:"name"`
}

func CreateBoard(ctx context.Context, b *CreateBoardData) error {
	if b.Name == "" {
		return errors.New("board name can not be empty")
	}

	err := gorm.G[storage.Board](storage.Db).Create(ctx, &storage.Board{Name: b.Name})
	if err != nil {
		return err
	}

	return nil
}

func DeleteBoard(ctx context.Context, boardID string) error {
	_, err := gorm.G[storage.Board](storage.Db).Where("id = ?", boardID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
