package api

import (
	"net/http"
	"encoding/json"
	"github.com/go-redis/redis"
	"fmt"
	"strconv"
)

func AddTodoHandler(c *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		todo := UnmarshallBody(r)

		res := Response{"Ok", "Successfully posted a new todo for id: " + todo.Id}
		response, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		var uid = "0"
		err = c.Get("i").Err()
		if(err != nil) {
		   	err := c.Set("i", "0", 0).Err()
			if(err != nil){
				panic(err)
			}
		} else {
		   	index, err := c.Incr("i").Result()
			uid = strconv.FormatInt(index, 10)
		   	if(err != nil){
			   	panic(err)
		   	}

			err = c.Set("i", index, 0).Err()
			if( err!= nil ){
				panic(err)
			}
		}

		////////JUST A TEST BUT WORKS
		var m = make(map[string]interface{})
		m["id"] = uid
		m["content"] = todo.Content

		fmt.Println(uid)

		hash, err := c.HMSet("k", m).Result()
		if(err != nil){
			panic(err)
		}

		mmap, err := json.Marshal(m)
		if(err != nil){
			panic(err)
		}
		s, _ := strconv.Unquote(string(mmap))
		
		err = c.LPush("mah", s).Err()
		if(err != nil){
			panic(err)
		}
		fmt.Println(hash)

		w.WriteHeader(http.StatusOK)
		w.Write(response)

	}

}
