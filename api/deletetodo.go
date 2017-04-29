package api

import (
	"net/http"
	"fmt"
)

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got a DELETE!")
}
