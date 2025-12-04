package commands

import (
	"cli-todo/storage"
	"context"
	"errors"
	"gorm.io/gorm"
)

type AddTaskData struct {
	Text string `json:"text"`
}
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func AddTask(ctx context.Context, a *AddTaskData) error {
	if a.Text == "" {
		return errors.New("empty text")
	}

	err := gorm.G[storage.Task](storage.Db).Create(ctx, &storage.Task{Text: a.Text, Done: false})
	if err != nil {
		return err
	}
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
	_, err := gorm.G[storage.Task](storage.Db).Where("id = ?", id).Updates(ctx, storage.Task{Done: true})
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
