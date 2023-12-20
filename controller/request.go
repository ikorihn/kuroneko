package controller

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type History struct {
	StatusCode    int
	Header        http.Header
	Body          []byte
	ContentLength int
	ExecutionTime time.Time

	Request *http.Request
}

func (c *Controller) Send(method, requestUrl, contentType string, headers map[string]string, body []byte) (*History, error) {
	var bbuf io.Reader
	if body != nil {
		bbuf = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, requestUrl, bbuf)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &History{
		StatusCode:    res.StatusCode,
		Header:        res.Header,
		Body:          b,
		ContentLength: int(res.ContentLength),
		ExecutionTime: time.Now(),
		Request:       req,
	}, err
}
