package api

import (
	"net/http"
	"encoding/json"
	"github.com/go-redis/redis"
)

func AddTodoHandler(c *redis.Client) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		todo := UnmarshallBody(r)

		res := Response{"Ok", "Successfully posted a new todo for id: " + todo.Id}
		response, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		err = c.Set(todo.Id, todo.Content, 0).Err()
		if (err != nil) {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)

	}

}
