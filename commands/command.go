package commands

import (
	"cli-todo/storage"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AddTaskData struct {
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
}
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func AddTask(ctx context.Context, a *AddTaskData) error {
	var board storage.Board

	if a.Text == "" {
		return errors.New("empty text")
	}

	err := storage.Db.WithContext(ctx).
		Where("user_id = ? AND name = ?", a.UserID, "Todo").
		First(&board).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		board = storage.Board{
			UserID: a.UserID,
			Name:   "Todo",
		}
		err = storage.Db.WithContext(ctx).Create(&board).Error
		if err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to find board: %w", err)
	}

	err = gorm.G[storage.Task](storage.Db).Create(ctx, &storage.Task{BoardID: board.ID, Text: a.Text})
	if err != nil {
		return err
	}
	fmt.Println(a.UserID)
	return nil
}

func UpdateTask(ctx context.Context, id int, task *storage.Task) error {
	_, err := gorm.G[storage.Task](storage.Db).Where("id = ?", id).Updates(ctx, *task)
	if err != nil {
		return err
	}
	return nil
}

func DoneTask(ctx context.Context, id int64) error {
	task, err := gorm.G[storage.Task](storage.Db).Where("id = ?", id).First(ctx)
	if err != nil {
		return err
	}

	var board storage.Board
	err = storage.Db.WithContext(ctx).Where("id = ?", task.BoardID).First(&board).Error
	if err != nil {
		return fmt.Errorf("failed to find current board: %w", err)
	}

	var doneBoard storage.Board
	err = storage.Db.WithContext(ctx).Where("user_id = ? AND name = ?", board.UserID, "Done").First(&doneBoard).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		doneBoard = storage.Board{
			UserID: board.UserID,
			Name:   "Done",
		}
		err = storage.Db.WithContext(ctx).Create(&doneBoard).Error
		if err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to find board: %w", err)
	}

	task.BoardID = doneBoard.ID
	err = storage.Db.WithContext(ctx).Save(&task).Error
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func DeleteTask(ctx context.Context, id int64) error {
	_, err := gorm.G[storage.Task](storage.Db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
