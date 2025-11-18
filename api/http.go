package api

import (
	"cli-todo/commands"
	"cli-todo/storage"
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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	atoi, err := strconv.Atoi(taskID)
	if err != nil {
		msg := "fail to write HTTP response: "
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}
	var data commands.AddTaskData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := "fail to write HTTP response: "
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}
	task := commands.Task{Text: data.Text}
	err = commands.UpdateTask(atoi, &task)
	if err != nil {
		w.Write([]byte("fail to update task: " + err.Error()))
		return
	}
	w.Write([]byte("task updated"))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := storage.Db.Query("SELECT id, task, done FROM tasks")
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
	defer rows.Close()

	var tasks []commands.Task
	for rows.Next() {
		var t commands.Task
		err := rows.Scan(&t.ID, &t.Text, &t.Done)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("fail to scan task: " + err.Error()))
			return
		}
		tasks = append(tasks, t)
	}

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
	err = commands.DeleteTask(int64(atoi))
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
	router.HandleFunc("/add", addHandler)
	router.HandleFunc("/list", listHandler)
	router.HandleFunc("/done", doneHandler)
	router.HandleFunc("/delete", deleteHandler)
	router.HandleFunc("/update", updateHandler)

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
