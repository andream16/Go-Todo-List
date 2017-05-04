package api

import (
	"net/http"
	"encoding/json"
	"github.com/go-redis/redis"
	"strconv"
	"log"
)

func AddTodoHandler(c *redis.Client) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		todo := UnmarshallBody(r)

		var uid = "0"
		err := c.Get("i").Err()
		if(err != nil) {
		   	err := c.Set("i", "0", 0).Err()
			if(err != nil){
				log.Fatalf("Something went wrong while setting a new counter i on Redis : " + err.Error())
				return
			}
		} else {
		   	index, err := c.Incr("i").Result()
			uid = strconv.FormatInt(index, 10)
		   	if(err != nil){
				log.Fatalf("Something went wrong while incrementing counter i: " + err.Error())
			   	return
		   	}

			err = c.Set("i", index, 0).Err()
			if( err!= nil ){
				log.Fatalf("Something went wrong while setting counter i to a new value: " + err.Error())
				return
			}
		}

		var m = make(map[string]interface{})
		m["id"] = uid
		m["content"] = todo.Content

		t, err := json.Marshal(m)
		if(err != nil){
			log.Fatalf("Something went wrong while marshalling m: " + err.Error())
			return
		}
		s, _ := strconv.Unquote(string(t))
		
		err = c.LPush("todos", s).Err()
		if(err != nil){
			log.Fatalf("Something went while pushing a new todo " + s +" to todos:" + err.Error())
			return
		}

		res := Response{"Ok", "Successfully posted a new todo for id: " + s}
		response, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Something went wrong while unmarshaling JSON : " + err.Error())
			return
		}

		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write(response)
		
	}

}
