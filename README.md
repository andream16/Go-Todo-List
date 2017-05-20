# Go-Todo-List

A simple RESTful application built with [go](https://github.com/golang), [gorilla/mux](https://github.com/gorilla/mux) and [go-redis](https://github.com/go-redis/redis).

# What you need to run it

 - [go1.8](https://golang.org/doc/devel/release.html#go1.8)
 - [gorilla/mux](https://github.com/gorilla/mux). You can get it using `go get -u github.com/gorilla/mux` inside your `$GOPATH`
 - [go-redis](https://github.com/go-redis/redis). You can get it using `go get -u github.com/go-redis/redis`

# How to run it

 - Start redis daemon `sudo service start redis` or `sudo systemctl start redis`
 - run the main using `go run main/main.go` or build the project with `go build` and run it's bin

# What you can Actually do

 - Get all todos : `curl -H "Content-Type: application/json" -X GET http://localhost:8080`
 - Get one todo  : `curl -H "Content-Type: application/json" -X GET -d '{"id":"1"}' http://localhost:8080`
 - Add a new todo: `curl -H "Content-Type: application/json" -X POST -d '{"content":"Have some Nugs"}' http://localhost:8080`
 - Edit a todo   : `curl -H "Content-Type: application/json" -X PUT -d '{"id":"1", "content":"New Content"}' http://localhost:8080`
 - Delete a todo : `curl -H "Content-Type: application/json" -X DELETE -d '{"id":"1"}' http://localhost:8080`
 
 Swagger reference available [here](https://app.swaggerhub.com/apis/AndreaM16/Go-Todo-List/1.0.0). 
