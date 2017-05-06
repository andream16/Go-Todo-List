package api

import (
	"net/http"
	"fmt"
)

type Todo struct {
	Content string `json:"Content"`
}

type SliceResponse struct {
	Status string
	Data []string
}

type Response struct {
	Status string
	Data string
}

func IndexTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "BE is alive!")
}
