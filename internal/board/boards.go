package board

import (
	"api-todo-go/internal/storage"
	"context"
	"errors"
	"gorm.io/gorm"
)

type CreateBoardData struct {
	Name        string `json:"board_name"`
	Description string `json:"board_description"`
	Status      string `json:"board_status"`
}

func CreateBoard(ctx context.Context, b *CreateBoardData) error {
	if b.Name == "" {
		return errors.New("board name can not be empty")
	}

	err := gorm.G[storage.Board](storage.Db).Create(ctx, &storage.Board{Name: b.Name, Description: b.Description, Status: b.Status})
	if err != nil {
		return err
	}

	return nil
}

func DeleteBoard(ctx context.Context, boardID int) error {
	_, err := gorm.G[storage.Board](storage.Db).Where("id = ?", boardID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBoard(ctx context.Context, id int, board *storage.Board) error {
	_, err := gorm.G[storage.Board](storage.Db).Where("id = ?", id).Updates(ctx, *board)
	if err != nil {
		return err
	}
	return nil
}
