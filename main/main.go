package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"../redismanager"
	"../cassandramanager"
	"../api"
	"os"
)

func main() {

	//Initialize Redis Client
	client, err := redismanager.InitRedisClient()
	if(err != nil){
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	cluster, err := cassandramanager.InitCassandraCluster()
	if(err != nil){
		log.Fatalf(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Cassandra Cluster Successfully Initialized . . .")
	}

	fmt.Println(cluster)

	//Initialize Router Handlers
	r := mux.NewRouter()
	r.HandleFunc("/", api.IndexTodoHandler).
			      Methods("GET")
	r.HandleFunc("/todo/", api.GetTodoHandler(&client)).
		              Methods("GET")
	r.HandleFunc("/todo", api.AddTodoHandler(&client)).
		              Methods("POST")
	r.HandleFunc("/todo/", api.EditTodoHandler(&client)).
		              Methods("PUT")
	r.HandleFunc("/todo/", api.DeleteTodoHandler(&client)).
		              Methods("DELETE")
	http.Handle("/", r)

	fmt.Println("Listening on port :8000 . . .")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
	
}
