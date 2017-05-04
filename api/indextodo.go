package api

import (
	"net/http"
	"fmt"
)

type Todo struct {
	Content string `json:"Content"`
}

type Response struct {
	Status string
	Description string
}

func IndexTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "BE is alive!")
}
