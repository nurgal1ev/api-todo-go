package board

import (
	"api-todo-go/internal/errors"
	"api-todo-go/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	var data CreateBoardData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	err = CreateBoard(r.Context(), &data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	if boardID == "" {
		msg := "fail to write HTTP response: task id is required"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
		}
		return
	}

	atoi, err := strconv.Atoi(boardID)
	if err != nil {
		msg := "BoardId is invalid"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		return
	}

	err = DeleteBoard(r.Context(), atoi)
	if err != nil {
		return
	}
	write, err := w.Write([]byte("board deleted"))
	if err != nil {
		return
	}
	fmt.Println(write)

	w.WriteHeader(http.StatusOK)
}

func GetBoardHandler(w http.ResponseWriter, r *http.Request) {
	boards, err := gorm.G[storage.Board](storage.Db).Find(r.Context())
	dataBoards, err := json.Marshal(boards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("fail to marshal boards: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(dataBoards)
	if err != nil {
		fmt.Println("fail to write HTTP response: " + err.Error())
		return
	}
}

func UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	atoi, err := strconv.Atoi(boardID)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	var data CreateBoardData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	board := storage.Board{Name: data.Name, Description: data.Description}
	err = UpdateBoard(r.Context(), atoi, &board)
	if err != nil {
		w.Write([]byte("fail to update board: " + err.Error()))
		return
	}

	w.Write([]byte("board updated"))
}

func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "boardID")

	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	type InviteUserData struct {
		UserID uint   `json:"userId"`
		Role   string `json:"role"`
	}

	var data InviteUserData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to write HTTP response: ")
		return
	}

	err = InviteUserToBoard(r.Context(), data.UserID, uint(boardID), data.Role)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("board invited"))
}
