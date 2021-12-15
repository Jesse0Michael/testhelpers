package testserver

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Handlers struct {
	handlers map[string][]Handler
	index    map[string]int
	sync.RWMutex
}

// NewHandlers initializes a new set of handlers
func NewHandlers(handle ...Handler) *Handlers {
	handlers := map[string][]Handler{}
	index := map[string]int{}
	for _, h := range handle {
		path := strings.TrimPrefix(h.Path, "/")
		handlers[path] = append(handlers[path], h)
		index[path] = 0
	}
	return &Handlers{
		handlers: handlers,
		index:    index,
	}
}

func (h *Handlers) Get(r *http.Request) (Handler, error) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	handlers := h.handlers[path]
	if len(handlers) == 0 {
		handlers = h.handlers[""]
		path = ""
		if len(handlers) == 0 {
			return Handler{}, fmt.Errorf("no path available")
		}
	}

	h.Lock()
	handler := handlers[h.index[path]]
	if h.index[path] == len(handlers)-1 {
		h.index[path] = 0
	} else {
		h.index[path]++
	}
	h.Unlock()

	return handler, nil
}

// RoundTrip mocks a HTTP round trip for the provided handlers
func (h *Handlers) RoundTrip(r *http.Request) (*http.Response, error) {
	handler, err := h.Get(r)
	if err != nil {
		return nil, err
	}
	return handler.RoundTrip(r)
}

// Request is an alias for RoundTrip
func (h *Handlers) Request(r *http.Request) (*http.Response, error) {
	return h.RoundTrip(r)
}
