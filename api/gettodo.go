package api

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
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

		var response []byte

		if(len(todos) > 0){
			res := &SliceResponse{Status: "Ok"}
			for _, d := range todos {
				todo := &Todo{}
				json.Unmarshal([]byte(d), todo)
				res.Data = append(res.Data, todo)
			}
			response, _ = json.Marshal(res)
			
		} else {
			todoErrorHandler(w, r, http.StatusNotFound, content)
			return
		}

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
