package lib_test

import (
	"encoding/json"
	"testing"

	vcr "github.com/ad2games/vcr-go"
	"github.com/nkprince007/igitt-go/lib"
)

type response struct {
	UserId    int  `json:"userId"`
	Id        int  `json:"id"`
	Completed bool `json:"completed"`
}

func TestGet(t *testing.T) {
	vcr.Start("get", nil)
	defer vcr.Stop()

	var got response

	body, err := lib.Get("bleh", "https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Error(err)
	}

	want := response{1, 1, false}
	if got != want {
		t.Errorf("Get failed, got: '%v', want: '%v'", got, want)
	}
}
