package errors

import (
	"fmt"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error, msg string) {
	msg = msg + ": " + err.Error()
	_, err = w.Write([]byte(msg))
	if err != nil {
		fmt.Println("fail to write HTTP response: " + err.Error())
		return
	}
	fmt.Println(msg)
	return
}
