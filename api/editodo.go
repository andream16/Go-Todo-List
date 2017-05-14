package api

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/go-redis/redis"
	"log"
)

func EditTodoHandler (c *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		content := r.URL.Query().Get("content")

		if(len(id) == 0 && len(content) == 0){
			todoErrorHandler(w, r, http.StatusNotFound, id)
			return
		}

		todos, err := c.LRange("todos", 0, -1).Result()

		if (err == redis.Nil || len(todos) == 0) {
			todoErrorHandler(w, r, http.StatusNotFound, id)
			return
		} else if (err != nil) {
			todoErrorHandler(w, r, http.StatusProcessing, id)
			return
		}

		t := &Todo{Id:id, Content:content}

		var m = make(map[string]interface{})
		m["id"] = id
		m["content"] = content
		td, _ := json.Marshal(m)

		var response []byte
		res := &SliceResponse{Status: "Ok"}

		if (len(id) > 0) {
			for _, d := range todos {
				todo := &Todo{}
				json.Unmarshal([]byte(d), todo)
				if (strings.EqualFold(todo.Id, id)) {

					var mt = make(map[string]interface{})
					mt["id"] = todo.Id
					mt["content"] = todo.Content
					mrs, _ := json.Marshal(mt)

					err := c.LRem("todos", 0, string(mrs)).Err()
					if(err != nil){
						log.Fatalf("Something went wrong while deleting a todo on Redis : " + err.Error())
						return
					}
					err = c.LPush("todos", string(td)).Err()
					if(err != nil){
						log.Fatalf("Something went wrong while pushing a new todo on Redis : " + err.Error())
						return
					}
					res.Data = append(res.Data, t)
				}
			}
			response, _ = json.Marshal(res)
		} else {
			todoErrorHandler(w, r, http.StatusNotFound, id)
			return
		}

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}