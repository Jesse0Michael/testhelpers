package testserver

import (
	"net/http"
	"net/http/httptest"
)

// Hander defines the response from your http test server
type Handler struct {
	Status   int
	Response []byte
}

// Serve will create and serve this test handler from a new http test server
func (h Handler) Serve() *Server {
	return New(h)
}

// Server wraps a httptest server and will loop over its handlers to return a response from the test server
type Server struct {
	*httptest.Server
	handlers []Handler
	index    int
}

// Creates a new server that will loop over the specified response handlers when serving a response
func New(h ...Handler) *Server {
	server := &Server{Server: httptest.NewServer(nil), handlers: h}
	server.Config.Handler = server.handle()
	return server
}

// handle serving the correct response handler by looping over the available handlers
func (s *Server) handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := s.handlers[s.index]
		if handler.Status == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(handler.Status)
		}

		_, _ = w.Write(handler.Response)

		if s.index == len(s.handlers)-1 {
			s.index = 0
		} else {
			s.index++
		}
	})
}
