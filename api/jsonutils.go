package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

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
