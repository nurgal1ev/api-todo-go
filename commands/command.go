package commands

import (
	"cli-todo/storage"
	"errors"
)

type AddTaskData struct {
	Text string `json:"text"`
}
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func AddTask(a *AddTaskData) error {
	if a.Text == "" {
		return errors.New("empty text")
	}

	task := Task{
		Text: a.Text,
		Done: false,
	}
	statement, err := storage.Db.Prepare("INSERT INTO tasks (task, done) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(task.Text, task.Done)
	if err != nil {
		return err
	}
	return nil
}

func DoneTask(id int64) error {
	statement, err := storage.Db.Prepare("UPDATE tasks SET done = true WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(id int64) error {
	statement, err := storage.Db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
