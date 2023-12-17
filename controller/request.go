package controller

import (
	"io"
	"net/http"
)

var httpClient *http.Client = http.DefaultClient

type Response struct {
	StatusCode    int
	Header        http.Header
	Body          []byte
	ContentLength int

	Request *http.Request
}

func Send(method, requestUrl string) (*Response, error) {
	req, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode:    res.StatusCode,
		Header:        res.Header,
		Body:          b,
		ContentLength: int(res.ContentLength),
		Request:       req,
	}, err
}
