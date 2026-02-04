package commands

import (
	"api-todo-go/internal/storage"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AddTaskData struct {
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
}

type MoveTaskData struct {
	StatusID uint `json:"status_id"`
}

func AddTask(ctx context.Context, a *AddTaskData) error {
	if a.Text == "" {
		return errors.New("empty text")
	}

	err := gorm.G[storage.Task](storage.Db).Create(ctx, &storage.Task{UserID: a.UserID, Text: a.Text, StatusID: 1})
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

func MoveTask(ctx context.Context, taskID uint, statusID uint) error {
	_, err := gorm.G[storage.Task](storage.Db).Where("id = ?", taskID).Update(ctx, "status_id", statusID)
	if err != nil {
		return err
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
