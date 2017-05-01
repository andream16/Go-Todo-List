package api

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
)

func todosErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if (status == http.StatusNotFound) {
		fmt.Fprint(w, "No todos Found")
	} else if(status == http.StatusInternalServerError){
		fmt.Fprint(w, "Something went wrong!")
	}
}

func GetTodosHandler(c *redis.Client) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {


		todos, err := c.LRange("todos", 0, -1).Result()
		if(len(todos) < 1) {
			todosErrorHandler(w, r, http.StatusNotFound)
			return
		}

		finalTodos := make([]string, len(todos))

		for v := range todos {
			val, err := json.Marshal(v)
			if(err != nil){
				todosErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
			finalTodos = append(finalTodos, string(val))
		}

		if(err != nil){
			todosErrorHandler(w, r, http.StatusInternalServerError)
			return
		}


		response, err := json.Marshal(todos)
		if err != nil {
			panic(err)
		}

		fmt.Println(response)
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)

	}
}
