package testserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func ExampleServer() {
	testServer := New(Handler{Status: http.StatusCreated, Response: []byte("success")},
		Handler{Status: http.StatusRequestTimeout, Response: []byte("timeout")})
	defer testServer.Close()

	resp, _ := http.Get(testServer.URL)
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s\n", resp.StatusCode, string(b))
	resp.Body.Close()

	resp, _ = http.Get(testServer.URL)
	b, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s\n", resp.StatusCode, string(b))
	resp.Body.Close()

	resp, _ = http.Get(testServer.URL)
	b, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s\n", resp.StatusCode, string(b))
	resp.Body.Close()
	// Output:
	// 201: success
	// 408: timeout
	// 201: success
}

func TestServer_handle(t *testing.T) {
	testServer := New(Handler{Status: http.StatusCreated, Response: []byte("success")},
		Handler{Response: []byte("OK")},
		Handler{Status: http.StatusConflict})
	defer testServer.Close()

	tests := []struct {
		status int
		body   string
	}{
		{status: 201, body: "success"},
		{status: 200, body: "OK"},
		{status: 409, body: ""},
		{status: 201, body: "success"},
	}
	for _, tt := range tests {
		resp, err := http.Get(testServer.URL)
		if err != nil {
			t.Errorf("Unexpected server response failure = %v", err)
		}
		if resp.StatusCode != tt.status {
			t.Errorf("Handler.Serve().Status = %v, want %v", resp.StatusCode, tt.status)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Unexpected response body failure = %v", err)
		}
		if string(b) != tt.body {
			t.Errorf("Handler.Serve().Body = %v, want %v", string(b), tt.body)
		}
		resp.Body.Close()
	}
}
