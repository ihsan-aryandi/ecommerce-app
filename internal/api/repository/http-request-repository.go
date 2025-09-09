package repository

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpHeaders map[string]string

func fetch(method, url string, body io.Reader, headers HttpHeaders) (response *http.Response, result []byte, err error) {
	client := &http.Client{}

	// Init request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	// Set Headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Fetch request
	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	result, err = io.ReadAll(response.Body)
	return
}

type QueryParams map[string]string

func (q QueryParams) ToString(withQuestionMark bool) (result string) {
	if len(q) == 0 {
		return ""
	}

	result = "&"
	if withQuestionMark {
		result = "?"
	}

	for key, value := range q {
		result += fmt.Sprintf("%s=%s&", key, value)
	}

	return strings.TrimRight(result, "&")
}
