package rest

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{}

func Client(method, url string, headers map[string]string, body []byte, timeout time.Duration) (*http.Response, []byte, error) {
	client.Timeout = timeout
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	rbody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return response, rbody, nil
}
