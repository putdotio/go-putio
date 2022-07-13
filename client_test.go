package putio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("got: %v, want: %v", r.Method, want)
	}
}

func testHeader(t *testing.T, r *http.Request, key, value string) { // nolint
	if r.Header.Get(key) != value {
		t.Errorf("missing header. want: %q: %q", key, value)
	}
}

func TestNewClient(t *testing.T) {
	cl := NewClient(nil)
	if cl.BaseURL.String() != defaultBaseURL {
		t.Errorf("got: %v, want: %v", cl.BaseURL.String(), defaultBaseURL)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	cl := NewClient(nil)
	_, err := cl.NewRequest(context.Background(), http.MethodGet, ":", nil)
	if err == nil {
		t.Errorf("bad URL accepted")
	}
}

func TestNewRequest_customUserAgent(t *testing.T) {
	userAgent := "test"
	cl := NewClient(nil)
	cl.UserAgent = userAgent

	req, _ := cl.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	if got := req.Header.Get("User-Agent"); got != userAgent {
		t.Errorf("got: %v, want: %v", got, userAgent)
	}
}
