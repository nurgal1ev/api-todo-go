package board

import (
	"api-todo-go/internal/storage"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CreateBoardData struct {
	Name        string `json:"board_name"`
	Description string `json:"board_description"`
}

var defaultStatuses = []string{"todo", "in progress", "done"}

func CreateBoard(ctx context.Context, b *CreateBoardData) error {
	if b.Name == "" {
		return errors.New("board name can not be empty")
	}

	currentBoard := &storage.Board{
		Name:        b.Name,
		Description: b.Description,
	}

	err := gorm.G[storage.Board](storage.Db).Create(ctx, currentBoard)
	if err != nil {
		return err
	}

	for _, status := range defaultStatuses {
		err := gorm.G[storage.Status](storage.Db).Create(ctx, &storage.Status{
			Name:    status,
			BoardID: currentBoard.ID,
		})
		if err != nil {
			return err
		}
	}
	fmt.Printf("Created board with ID: %d\n", currentBoard.ID)

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
