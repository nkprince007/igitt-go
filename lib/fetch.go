package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ResponseError indicates an error while fetching a HTTP response.
type ResponseError struct {
	code int
	s    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("ResponseError: Response('%d', '%s')", e.code, e.s)
}

var client = &http.Client{
	Timeout: time.Second * 5,
}

func prepareRequest(token string, url string, method string, data []byte) (
	*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func sendRequest(req *http.Request, data chan<- []byte, e chan<- error) {
	defer close(e)
	defer close(data)

	resp, err := client.Do(req)
	if err != nil {
		e <- err
		data <- nil
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e <- err
		data <- nil
		return
	}

	if resp.StatusCode >= 300 && resp.StatusCode != 304 {
		e <- &ResponseError{resp.StatusCode, string(body)}
		data <- nil
		return
	}

	e <- nil
	data <- body
}

// Get retrieves data from URL.
func Get(token string, url string, data chan<- []byte, e chan<- error) {
	req, err := prepareRequest(token, url, "GET", nil)
	if err != nil {
		e <- err
		data <- nil
	}

	go sendRequest(req, data, e)
}

// Post sends data to URL.
func Post(token string, url string, data chan []byte, e chan<- error) {
	req, err := prepareRequest(token, url, "POST", <-data)
	if err != nil {
		e <- err
		data <- nil
	}

	go sendRequest(req, data, e)
}

// Delete removes data from URL.
func Delete(token string, url string, data chan []byte, e chan<- error) {
	req, err := prepareRequest(token, url, "DELETE", <-data)
	if err != nil {
		e <- err
		data <- nil
	}
	go sendRequest(req, data, e)
}
