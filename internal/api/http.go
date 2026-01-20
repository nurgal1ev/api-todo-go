package api

import (
	"cli-todo/internal/auth"
	"cli-todo/internal/commands"
	"cli-todo/internal/errors"
	"cli-todo/internal/storage"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	var data commands.AddTaskData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}
	fmt.Println(data)
	err = commands.AddTask(r.Context(), &data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	atoi, err := strconv.Atoi(taskID)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}
	var data commands.AddTaskData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}
	task := storage.Task{Text: data.Text}
	err = commands.UpdateTask(r.Context(), atoi, &task)
	if err != nil {
		w.Write([]byte("fail to update task: " + err.Error()))
		return
	}
	w.Write([]byte("task updated"))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := gorm.G[storage.Task](storage.Db).Find(r.Context())
	dataTasks, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("fail to marshal tasks: " + err.Error()))
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

	err = commands.DoneTask(r.Context(), int64(atoi))
	if err != nil {
		return
	}
	write, err := w.Write([]byte("task done"))
	if err != nil {
		return
	}
	fmt.Println(write)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
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
			return
		}
	}
	err = commands.DeleteTask(r.Context(), int64(atoi))
	if err != nil {
		return
	}
	write, err := w.Write([]byte("task deleted"))
	if err != nil {
		return
	}
	fmt.Println(write)
}

func HTTPServer() {
	router := http.NewServeMux()
	router.Handle("/add", auth.AuthMiddleware(http.HandlerFunc(addHandler)))
	router.Handle("/list", auth.AuthMiddleware(http.HandlerFunc(listHandler)))
	router.Handle("/done", auth.AuthMiddleware(http.HandlerFunc(doneHandler)))
	router.Handle("/delete", auth.AuthMiddleware(http.HandlerFunc(deleteHandler)))
	router.Handle("/update", auth.AuthMiddleware(http.HandlerFunc(updateHandler)))
	router.HandleFunc("/auth/register", auth.Register)
	router.HandleFunc("/auth/login", auth.Login)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("http server started")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("HTTP server error", err)
		return
	}
}
