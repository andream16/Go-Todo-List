package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis"
	"strings"
	"bytes"
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
		if (len(content) == 0) {
			todoErrorHandler(w, r, http.StatusBadRequest, content)
			return
		}

		todos, err := c.LRange("todos", 0, -1).Result()
		if (err == redis.Nil || len(todos) == 0) {
			todoErrorHandler(w, r, http.StatusNotFound, content)
			return
		} else if (err != nil) {
			todoErrorHandler(w, r, http.StatusProcessing, content)
			return
		}

		var val string
		for _, k := range todos {
			if (strings.Contains(k, content)) {
				tmp, _ := json.Marshal(Todo{k})
				val = string(tmp)
			}
		}

		var response []byte
		var parsedTodos []string

		if (len(val) == 0) {
			for _, cont := range todos {
				currTodo, err := json.Marshal(Todo{cont})
				if(err != nil) {
					todoErrorHandler(w, r, http.StatusProcessing, content)
					return
				}
				currParsedTodo := string(currTodo)
				parsedTodos  = append(parsedTodos, currParsedTodo)
			}
		}

		if(len(val) > 0){
			res := Response{"Ok", val}
			response, _ = json.Marshal(res)
		} else if(len(parsedTodos) > 0){
			res := SliceResponse{"Ok", parsedTodos}
			response, _ = json.Marshal(res)
		}

		response = bytes.Replace(response, []byte("\\"), []byte(""), -1)

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
