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

func TestNewRequest_resolveAbsoluteUrl(t *testing.T) {
	cl := NewClient(nil)
	fullUrl := "https://api.put.io/v2/files/list/continue"

	req, _ := cl.NewRequest(context.Background(), http.MethodGet, "/v2/files/list/continue", nil)
	if got := req.URL.String(); got != fullUrl {
		t.Errorf("got: %v, want: %v", got, fullUrl)
	}
}

func TestNewRequest_overrideBaseUrl(t *testing.T) {
	cl := NewClient(nil)
	fullUrl := "https://upload.put.io/v2/files/upload"

	req, _ := cl.NewRequest(context.Background(), http.MethodGet, "/v2/files/upload", nil)
	if got := req.URL.String(); got != fullUrl {
		t.Errorf("got: %v, want: %v", got, fullUrl)
	}
}

func TestNewRequest_overrideWithAbsoluteUrl(t *testing.T) {
	cl := NewClient(nil)
	fullUrl := "https://my-custom-url.com/endpoint?query=param"

	req, _ := cl.NewRequest(context.Background(), http.MethodGet, "https://my-custom-url.com/endpoint?query=param", nil)
	if got := req.URL.String(); got != fullUrl {
		t.Errorf("got: %v, want: %v", got, fullUrl)
	}
}
