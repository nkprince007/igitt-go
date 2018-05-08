package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ResponseError indicates an error while fetching a HTTP response
type ResponseError struct {
	code int
	s    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("ResponseError: Response('%d', '%s')", e.code, e.s)
}

var client = &http.Client{}

// Get retrieves data from URL
func Get(token string, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "token: "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 && resp.StatusCode != 304 {
		return nil, &ResponseError{resp.StatusCode, string(body)}
	}
	return body, err
}
