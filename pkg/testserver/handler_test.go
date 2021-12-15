package testserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jesse0michael/testhelpers/pkg/testhelpers"
)

func ExampleHandler_serve() {
	testServer := Handler{Response: testhelpers.LoadFile(&testing.T{}, "testdata/identity.json")}.Serve()
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	if err != nil {
		fmt.Println(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//     "id": "1a",
	//     "type": "hero"
	// }
}

func TestHandler_Request(t *testing.T) {
	tests := []struct {
		name        string
		handler     Handler
		wantStatus  int
		wantBody    string
		wantHeaders http.Header
	}{
		{
			name:        "round trip - empty",
			handler:     Handler{},
			wantStatus:  200,
			wantHeaders: http.Header{},
		},
		{
			name:        "round trip - status",
			handler:     Handler{Status: http.StatusAccepted},
			wantStatus:  202,
			wantHeaders: http.Header{},
		},
		{
			name:        "round trip - body",
			handler:     Handler{Response: []byte("test-response")},
			wantStatus:  200,
			wantBody:    "test-response",
			wantHeaders: http.Header{},
		},
		{
			name:        "round trip - headers",
			handler:     Handler{Headers: map[string]string{"Content-Type": "application/json"}},
			wantStatus:  200,
			wantHeaders: http.Header{"Content-Type": []string{"application/json"}},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp, _ := tt.handler.Request(req)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Handler.Request().Status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
			if tt.wantBody != "" {
				b, _ := io.ReadAll(resp.Body)
				if string(b) != tt.wantBody {
					t.Errorf("Handler.Request().Body = %v, want %v", string(b), tt.wantBody)
				}
				resp.Body.Close()
			}
			if !reflect.DeepEqual(resp.Header, tt.wantHeaders) {
				t.Errorf("Handlers.Request().Header = %v, want %v", resp.Header, tt.wantHeaders)
			}
		})
	}
}
