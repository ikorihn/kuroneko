package controller

import (
	"bytes"
	"context"
	"io"
	"maps"
	"net/http"
	"time"

	"github.com/tcnksm/go-httpstat"
)

const hdrContentType = "Content-Type"

type Request struct {
	Method      string
	Url         string
	ContentType string
	Headers     headerMap
	Body        []byte
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}

func NewRequestWithValues(
	method string,
	url string,
	contentType string,
	headers map[string]string,
	body []byte,
) *Request {
	return &Request{
		Method:      method,
		Url:         url,
		ContentType: contentType,
		Headers:     headers,
		Body:        body,
	}
}

func (r Request) ToHttpReq() *http.Request {
	var req *http.Request
	if len(r.Body) > 0 {
		b := io.NopCloser(bytes.NewBuffer(r.Body))
		req, _ = http.NewRequest(r.Method, r.Url, b)
	} else {
		req, _ = http.NewRequest(r.Method, r.Url, nil)
	}
	if r.ContentType != "" {
		req.Header.Add(hdrContentType, r.ContentType)
	}
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	return req
}

type Favorite struct {
	Request []Request `toml:"request"`
}

type History struct {
	StatusCode    int
	Header        http.Header
	Body          []byte
	ContentLength int
	ExecutionTime time.Time

	Request  Request
	HttpStat httpstat.Result
}

func (c *Controller) Send(ctx context.Context, method, requestUrl, contentType string, headers headerMap, body []byte) (*History, error) {
	var bbuf io.Reader
	if body != nil {
		bbuf = bytes.NewBuffer(body)
	}

	// collect stats
	var result httpstat.Result
	ctx = httpstat.WithHTTPStat(ctx, &result)

	req, err := http.NewRequestWithContext(ctx, method, requestUrl, bbuf)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Add(hdrContentType, contentType)
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
		Request: Request{
			Method:      method,
			Url:         requestUrl,
			ContentType: contentType,
			Headers:     maps.Clone(headers),
			Body:        body,
		},
		HttpStat: result,
	}, err
}
