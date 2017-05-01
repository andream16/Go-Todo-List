package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis"
)

func todoErrorHandler(w http.ResponseWriter, r *http.Request, status int, id string) {
	w.WriteHeader(status)
	if (status == http.StatusNotFound) {
		fmt.Fprint(w, "Cannot find " + id)
	} else if( status == http.StatusInternalServerError){
		fmt.Fprint(w, "Something went wrong while managing " + id)
	} else if(status == http.StatusBadRequest ){
		fmt.Fprint(w, "Not able to extract id " + id)
	} else if(status == http.StatusProcessing){
		fmt.Fprint(w, "Unable to marshal/unmarshal for todo " + id)
	}
}

func GetTodoHandler(c *redis.Client) func (w http.ResponseWriter, r *http.Request) {
     return func (w http.ResponseWriter, r *http.Request) {
	     vars := mux.Vars(r)
	     id, ok := vars["id"]
	     if (!ok) {
		     todoErrorHandler(w, r, http.StatusBadRequest, id)
		     return
	     }

	     val, err := c.Get(id).Result()
	     if (err == redis.Nil) {
		     todoErrorHandler(w, r, http.StatusNotFound, id)
		     return
	     } else if (err != nil) {
		     todoErrorHandler(w, r, http.StatusProcessing, id)
		     return
	     }

	     m, err := json.Marshal(Todo{id, val})
	     if (err != nil) {
		     todoErrorHandler(w, r, http.StatusProcessing, id)
		     return
	     }

	     res := Response{"Ok", string(m)}
	     response, err := json.Marshal(res)
	     if err != nil {
		     todoErrorHandler(w, r, http.StatusProcessing, id)
		     return
	     }

	     defer r.Body.Close()

	     w.WriteHeader(http.StatusOK)
	     w.Write(response)
     }
}
