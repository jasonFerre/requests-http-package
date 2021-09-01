package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type Requests interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client Requests
)

func init() {
	Client = &http.Client{}
}

// Get executes a http.Get request
func Fetch(url string, headers http.Header) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New(err.Error())
	}

	req.Header = headers
	resp, err := Client.Do(req)

	if err != nil {
		return "", errors.New(err.Error())
	}

	defer resp.Body.Close()
	return _readResponseBody(resp)
}

func Delete(url string, headers http.Header) (string, error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", errors.New(err.Error())
	}

	req.Header = headers
	resp, err := Client.Do(req)
	if err != nil {
		return "", errors.New(err.Error())
	}

	defer resp.Body.Close()
	return _readResponseBody(resp)
}

func Post(url string, body interface{}, headers http.Header) (string, error) {

	parsedBody, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	//resp, err := http.Post(url, "application/json", parsedBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(parsedBody))
	if err != nil {
		return "", err
	}

	req.Header = headers
	resp, err := Client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	return _readResponseBody(resp)
}

func _readResponseBody(response *http.Response) (string, error) {
	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
		return "", errors.New(err.Error())
	}

	responseMap := map[string]interface{}{
		"StatusCode": response.StatusCode,
		"Status":     response.Status,
		"Body":       string(bodyResp),
	}

	jsonResponse, err := json.Marshal(responseMap)

	if err != nil {
		return "", err
	}
	return string(jsonResponse), nil
}
