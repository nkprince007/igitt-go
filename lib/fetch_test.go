package lib_test

import (
	"encoding/json"
	"testing"

	vcr "github.com/ad2games/vcr-go"
	"github.com/nkprince007/igitt-go/lib"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type response struct {
		UserId    int  `json:"userId"`
		Id        int  `json:"id"`
		Completed bool `json:"completed"`
	}

	vcr.Start("get", nil)
	defer vcr.Stop()

	var got response

	data, e := make(chan []byte), make(chan error)
	go lib.Get("bleh", "https://jsonplaceholder.typicode.com/todos/1", data, e)
	if err := <-e; err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(<-data, &got); err != nil {
		t.Error(err)
	}

	want := response{1, 1, false}
	if got != want {
		t.Errorf("Get failed, got: '%v', want: '%v'", got, want)
	}
}

func TestPost(t *testing.T) {
	type response struct {
		UserId string `json:"id"`
		Job    string `json:"job"`
		Name   string `json:"name"`
	}

	vcr.Start("post", nil)
	defer vcr.Stop()

	data, e := make(chan []byte), make(chan error)
	go lib.Post("", "https://reqres.in/api/users", data, e)
	data <- []byte(`{"name": "morpheus", "job": "leader"}`)

	if err := <-e; err != nil {
		t.Error(err)
	}

	var got response
	body := <-data
	if err := json.Unmarshal(body, &got); err != nil {
		t.Error(err)
	}

	want := response{"753", "leader", "morpheus"}
	if got != want {
		t.Errorf("Post failed, got: '%v', want: '%v'", got, want)
	}
}

func TestDelete(t *testing.T) {
	vcr.Start("delete", nil)
	defer vcr.Stop()

	data, e := make(chan []byte), make(chan error)
	go lib.Delete("", "https://reqres.in/api/users/2", data, e)
	data <- nil

	if err := <-e; err != nil {
		t.Error(err)
	}

	got, want := <-data, make([]byte, 0)
	assert.Equal(t, want, got)
}
