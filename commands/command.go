package commands

import (
	"cli-todo/storage"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type AddTaskData struct {
	Text string `json:"text"`
}
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var Tasks []Task

func AddTask(a *AddTaskData) error {
	if a.Text == "" {
		return errors.New("empty text")
	}

	task := Task{
		ID:   len(Tasks) + 1,
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

	fmt.Printf("Добавлено: %d. [ ] %s\n", len(Tasks), a.Text)
	return nil
}

func LoadTasks() {
	_, err := os.Stat("tasks.json")
	if err != nil {
		fmt.Println("Файл не найден")
		return
	} else {
		data, _ := os.ReadFile("tasks.json")
		_ = json.Unmarshal(data, &Tasks)
	}
}

func SaveTask() {
	jsonData, err := json.MarshalIndent(Tasks, "", " ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_ = os.WriteFile("tasks.json", jsonData, 0644)
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

func DeleteTask(args []string) {
	LoadTasks()
	if len(args) == 0 {
		fmt.Println("Нужно ввести номер задачи")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Номер задачи должен быть числом")
		return
	}

	for i, task := range Tasks {
		if task.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			fmt.Printf("Задача %d удалена: %s\n", id, task.Text)
			break
		}
	}
	SaveTask()
}
