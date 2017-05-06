package api

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
	"bytes"
	"strings"
	"encoding/json"
)

func todoErrorHandler(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.WriteHeader(status)
	if (status == http.StatusNotFound) {
		fmt.Fprint(w, "Cannot find todo :" + content)
	} else if ( status == http.StatusInternalServerError) {
		fmt.Fprint(w, "Something went wrong while managing : " + content)
	} else if (status == http.StatusBadRequest ) {
		fmt.Fprint(w, "Not able to extract content : " + content)
	} else if (status == http.StatusProcessing) {
		fmt.Fprint(w, "Unable to marshal/unmarshal for todo : " + content)
	}
}

func GetTodoHandler(c *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		content := r.URL.Query().Get("content")

		todos, err := c.LRange("todos", 0, -1).Result()
		if (err == redis.Nil || len(todos) == 0) {
			todoErrorHandler(w, r, http.StatusNotFound, content)
			return
		} else if (err != nil) {
			todoErrorHandler(w, r, http.StatusProcessing, content)
			return
		}

		var parsedTodos []string
		if(len(content) > 0){
			for _, k := range todos {
				if (strings.Contains(k, content)) {
					parsedTodos = append(parsedTodos,k)
				}
			}
		} else {
			for _, k := range todos {
					parsedTodos = append(parsedTodos,k)
			}
		}
		

		var response []byte
		res := SliceResponse{"Ok", parsedTodos}
		response, _ = json.Marshal(res)

		if(len(parsedTodos) > 0){
			res := SliceResponse{"Ok", parsedTodos}
			response, _ = json.Marshal(res)
		} else {
			res := SliceResponse{"Err", parsedTodos}
			response, _ = json.Marshal(res)
		}

		response = bytes.Replace(response, []byte("\\"), []byte(""), -1)

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
