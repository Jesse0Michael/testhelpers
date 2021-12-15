package testserver

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Hander defines the endpoint and response from your http test server
type Handler struct {
	Path     string
	Status   int
	Response []byte
	Headers  map[string]string
}

// Serve will create and serve this test handler from a new http test server
func (h Handler) Serve() *Server {
	return New(h)
}

// RoundTrip mocks a HTTP round trip for the provided response values
func (h Handler) RoundTrip(r *http.Request) (*http.Response, error) {
	if h.Status == 0 {
		h.Status = http.StatusOK
	}
	body := ioutil.NopCloser(&bytes.Buffer{})
	if len(h.Response) > 0 {
		body = ioutil.NopCloser(bytes.NewBuffer(h.Response))
	}
	headers := http.Header{}
	for k, v := range h.Headers {
		headers[k] = []string{v}
	}
	return &http.Response{
		StatusCode: h.Status,
		Body:       body,
		Header:     headers,
	}, nil
}

// Request is an alias for RoundTrip
func (h Handler) Request(r *http.Request) (*http.Response, error) {
	return h.RoundTrip(r)
}
