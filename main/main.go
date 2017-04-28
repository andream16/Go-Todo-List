package main

import (
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
)

func InitRedisClient() redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr    : "localhost:6379",
		Password: "",
		DB      : 0, //default
	})

	pong, err := client.Ping().Result()
	if( err != nil ){
		fmt.Println("Cannot Initialize Redis Client ", err)
	}
	fmt.Println("Redis Client Successfully Initialized . . .", pong)

	return *client

}

func UnmarshallBody(r *http.Request) Todo {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		panic(err)
	}

	return todo
}

type Todo struct {
	 Id      string `json:"Id"`
	 Content string `json:"Content"`
}

type Response struct {
	Status string
	Description string
}

type RedisInstance struct {
	RInstance *redis.Client
}


func main() {

	//Initialize Redis Client
	client := InitRedisClient()
	//Get current redis instance to get passed to different Gorilla-Mux Handlers
	redisHandler := &RedisInstance{RInstance:&client}

	//Initialize Router Handlers
	r := mux.NewRouter()
	r.HandleFunc("/", IndexTodoHandler).
			      Methods("GET")
	r.HandleFunc("/todo/{id}", redisHandler.GetTodoHandler).
		              Methods("GET")
	r.HandleFunc("/todo", redisHandler.AddTodoHandler).
		              Methods("POST")
	r.HandleFunc("/todo", EditTodoHandler).
			      Methods("PUT")
	r.HandleFunc("/todo", DeleteTodoHandler).
		              Methods("DELETE")
	http.Handle("/", r)

	fmt.Println("Listening on port :8000 . . .")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))


}

func IndexTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "BE is alive!")
}

func (c *RedisInstance) GetTodoHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if(!ok){
		http.Error(w, "Cannot Extract Id from Request", http.StatusInternalServerError)
	}

	val, err := c.RInstance.Get(id).Result()
	if(err == redis.Nil){
		fmt.Println("Key " + id + " does not exist.")
	} else if(err != nil){
		panic(err)
	}

	m, err := json.Marshal(Todo{id, val})
	if(err != nil){
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

func (c *RedisInstance) AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todo := UnmarshallBody(r)

	//http.Error(w, http.StatusText(http.NoBody), http.NoBody)
	//w.WriteHeader(http.ErrBodyNotAllowed)

	res := Response{"Ok", "Successfully posted a new todo for id: " + todo.Id}
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err = c.RInstance.Set(todo.Id, todo.Content, 0).Err()
	if(err != nil){
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	                                         
}

func EditTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got a PUT!")
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got a DELETE!")
}
