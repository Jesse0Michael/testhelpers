package testserver

import (
	"net/http"
	"net/http/httptest"
)

// Hander defines the endpoint and response from your http test server
type Handler struct {
	Path     string
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
	handlers map[string][]Handler
	index    map[string]int
}

// Creates a new server that will loop over the specified response handlers when serving a response
func New(handle ...Handler) *Server {
	handlers := map[string][]Handler{}
	index := map[string]int{}
	for _, h := range handle {
		handlers[h.Path] = append(handlers[h.Path], h)
		index[h.Path] = 0
	}
	server := &Server{Server: httptest.NewServer(nil), handlers: handlers, index: index}
	server.Config.Handler = server.handle()
	return server
}

// handle serving the correct response handler by looping over the available handlers
func (s *Server) handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		handlers := s.handlers[path]
		if len(handlers) == 0 {
			handlers = s.handlers[""]
			path = ""
			if len(handlers) == 0 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}

		handler := handlers[s.index[path]]
		if handler.Status == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(handler.Status)
		}

		_, _ = w.Write(handler.Response)

		if s.index[path] == len(handlers)-1 {
			s.index[path] = 0
		} else {
			s.index[path]++
		}
	})
}
