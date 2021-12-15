package testserver

import (
	"net/http"
	"net/http/httptest"
)

// Server wraps a httptest server and will loop over its handlers to return a response from the test server
type Server struct {
	*httptest.Server
	*Handlers
}

// Creates a new server that will loop over the specified response handlers when serving a response
func New(handle ...Handler) *Server {
	server := &Server{Server: httptest.NewServer(nil), Handlers: NewHandlers(handle...)}
	server.Config.Handler = server.handle()
	return server
}

// handle serving the correct response handler by looping over the available handlers
func (s *Server) handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, err := s.Get(r)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		for k, v := range handler.Headers {
			w.Header().Add(k, v)
		}
		if handler.Status == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(handler.Status)
		}

		_, _ = w.Write(handler.Response)
	})
}
