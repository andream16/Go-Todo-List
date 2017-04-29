package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis"
)

func GetTodoHandler(c *redis.Client) func (w http.ResponseWriter, r *http.Request) {
     return func (w http.ResponseWriter, r *http.Request) {
	     vars := mux.Vars(r)
	     id, ok := vars["id"]
	     if (!ok) {
		     http.Error(w, "Cannot Extract Id from Request", http.StatusInternalServerError)
	     }

	     val, err := c.Get(id).Result()
	     if (err == redis.Nil) {
		     fmt.Println("Key " + id + " does not exist.")
	     } else if (err != nil) {
		     panic(err)
	     }

	     m, err := json.Marshal(Todo{id, val})
	     if (err != nil) {
		     panic(err)
	     }

	     res := Response{"Ok", string(m)}
	     response, err := json.Marshal(res)
	     if err != nil {
		     http.Error(w, err.Error(), http.StatusInternalServerError)
		     return
	     }

	     defer r.Body.Close()

	     w.WriteHeader(http.StatusOK)
	     w.Write(response)
     }
}
