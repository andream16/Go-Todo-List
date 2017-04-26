package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
)

type Todo struct {
	 Id      string
	 Content string
}

type Response struct {
	Status string
	Description string
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexTodoHandler).
			      Methods("GET")
	r.HandleFunc("/todo", GetTodoHandler).
		              Methods("GET")
	r.HandleFunc("/todo", AddTodoHandler).
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

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got a GET!")
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//http.Error(w, http.StatusText(http.NoBody), http.NoBody)
		panic(err)
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		//w.WriteHeader(http.ErrBodyNotAllowed)
		panic(err)
	}

	res := Response{"Ok", "Successfully posted a new todo for id: " + todo.Id}
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	
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
