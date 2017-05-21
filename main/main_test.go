package main

import (
	"testing"
	"../redismanager"
	"../api"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func TestRedisInitialization(t *testing.T){
	_, err := redismanager.InitRedisClient()
	if(err != nil){
		t.Fatalf("Unable to connect to redis")
	}
}

func TestIndexGet(t *testing.T){
	var client http.Client
	resp, err := client.Get("http://localhost:8000/")
	if(err != nil){
		t.Fatalf("Something went wrong while: GET /")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		r, err := ioutil.ReadAll(resp.Body)
		if(err != nil){
			t.Error("Unable to extract body from response: ", err)
		} else {
			b := string(r)
			if( b != "BE is alive!"){
				t.Fatalf("Received an invalid response: " + b)
			}
		}
	}

}

func TestPostTodo(t *testing.T){

	var jsonStr = []byte(`{"content" : "testing_testing"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if(err != nil){
		t.Fatalf("Unable to post a new todo")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Unable to post a new todo, got status bad status")
	}
	r, err := ioutil.ReadAll(resp.Body)
	if(err != nil){
		t.Error("Unable to extract body from response: ", err)
	} else {
		d := &api.Response{}
		json.Unmarshal([]byte(r), d)
		if(d.Status != "Ok"){
			t.Fatalf("Got a bad Status: " + d.Status)
		}
	}
}

func TestGetTodo(t *testing.T){

	var client http.Client
	resp, err := client.Get("http://localhost:8000/todo/")
	if(err != nil){
		t.Fatalf("Something went wrong while: GET /todo/")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		r, err := ioutil.ReadAll(resp.Body)
		if(err != nil){
			t.Error("Unable to extract body from response: ", err)
		} else {
			d := &api.SliceResponse{}
			json.Unmarshal([]byte(r), d)
			if(d.Status != "Ok"){
				t.Fatalf("Got a bad Status: " + d.Status)
			}

			l := len(d.Data)
			var i int
			for _, k := range d.Data {
				if(i == l-1){
					var jsonStr = []byte(`{"id" : "`+ k.Id +`"}`)
					req, err := http.NewRequest("GET", "http://localhost:8000/todo/", bytes.NewBuffer(jsonStr))
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if(err != nil){
						t.Fatalf("Unable to get a todo")
					}
					defer resp.Body.Close()
					if resp.StatusCode != 200 {
						t.Fatalf("Unable to get a todo, got status bad status")
					}
					r, err := ioutil.ReadAll(resp.Body)
					if(err != nil){
						t.Error("Unable to extract body from response: ", err)
					} else {
						d := &api.Response{}
						json.Unmarshal([]byte(r), d)
						if(d.Status != "Ok"){
							t.Fatalf("Got a bad Status: " + d.Status)
						}
					}

				}
				i++
			}
		}
	}

}

func TestUpdateTodo(t *testing.T){
	var client http.Client
	resp, err := client.Get("http://localhost:8000/todo/")
	if(err != nil){
		t.Fatalf("Something went wrong while: GET /todo/")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		r, err := ioutil.ReadAll(resp.Body)
		if(err != nil){
			t.Error("Unable to extract body from response: ", err)
		} else {
			d := &api.SliceResponse{}
			json.Unmarshal([]byte(r), d)
			if(d.Status != "Ok"){
				t.Fatalf("Got a bad Status: " + d.Status)
			}

			l := len(d.Data)
			var i int
			for _, k := range d.Data {
				if(i == l-1){
					var jsonStr = []byte(`{"id" : "`+ k.Id +`","content" : "testing_put_todo"}`)
					req, err := http.NewRequest("PUT", "http://localhost:8000/todo/", bytes.NewBuffer(jsonStr))
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if(err != nil){
						t.Fatalf("Unable to put a new todo")
					}
					defer resp.Body.Close()
					if resp.StatusCode != 200 {
						t.Fatalf("Unable to put a todo, got status bad status")
					}
					r, err := ioutil.ReadAll(resp.Body)
					if(err != nil){
						t.Error("Unable to extract body from response: ", err)
					} else {
						d := &api.Response{}
						json.Unmarshal([]byte(r), d)
						if(d.Status != "Ok"){
							t.Fatalf("Got a bad Status: " + d.Status)
						}
					}

				}
				i++
			}
		}
	}
}