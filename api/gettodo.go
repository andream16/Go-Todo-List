package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis"
	"strings"
	"strconv"
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
				val = k
			}
		}

		var m, parsedTodos []byte

		if (len(val) == 0) {
			for _, cont := range todos {
				currTodo, err := json.Marshal(Todo{cont})
				if(err != nil) {
					todoErrorHandler(w, r, http.StatusProcessing, content)
					return
				}
				parsedTodos = append(parsedTodos, currTodo...)
			}
			m, err = json.Marshal(string(parsedTodos))
			if (err != nil) {
				todoErrorHandler(w, r, http.StatusProcessing, content)
				return
			}
		} else {

			m, err = json.Marshal(Todo{val})
			if (err != nil) {
				todoErrorHandler(w, r, http.StatusProcessing, content)
				return
			}
		}

		res := Response{"Ok", string(m)}
		response, err := json.Marshal(res)
		if err != nil {
			todoErrorHandler(w, r, http.StatusProcessing, content)
			return
		}

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
