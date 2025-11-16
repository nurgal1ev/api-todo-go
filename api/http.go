package api

import (
	"cli-todo/commands"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	var data commands.AddTaskData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := "fail to write HTTP response: " + err.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}
	fmt.Println(data)
	err = commands.AddTask(&data)
	if err != nil {
		msg := "fail to write HTTP response: " + err.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}
}

// todo: convert to sql
func listHandler(w http.ResponseWriter, r *http.Request) {
	dataTasks, err := json.Marshal(commands.Tasks)
	if err != nil {
		msg := "fail to write HTTP response: " + err.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(dataTasks)
	if err != nil {
		fmt.Println("fail to write HTTP response: " + err.Error())
		return
	}
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")
	if taskId == "" {
		msg := "fail to write HTTP response: task id is required"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
		}
		return
	}
	atoi, err := strconv.Atoi(taskId)
	if err != nil {
		msg := "taskId is invalid"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
		}
		return
	}

	err = commands.DoneTask(int64(atoi))
	if err != nil {
		return
	}
	write, err := w.Write([]byte("task done"))
	if err != nil {
		return
	}
	fmt.Println(write)

}

// todo: convert to sql
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	commands.LoadTasks()
	taskId := r.URL.Query().Get("id")
	if taskId == "" {
		msg := "fail to write HTTP response: task id is required"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
		}
		return
	}
	commands.DeleteTask([]string{taskId})
	write, err := w.Write([]byte("task deleted"))
	if err != nil {
		return
	}
	fmt.Println(write)
}

func HTTPServer() {
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/done", doneHandler)
	http.HandleFunc("/delete", deleteHandler)
	fmt.Println("http server started")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server error", err)
		return
	}
}
