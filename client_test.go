package putio

import (
	"context"
	"errors"
	"fmt"
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

func writeResponse(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := fmt.Fprintln(w, body); err != nil {
		t.Errorf("write response: %v", err)
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

func TestValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		setup()
		defer teardown()

		mux.HandleFunc("/v2/oauth2/validate", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			_, _ = w.Write([]byte(`{"result":true,"user_id":123}`))
		})

		got, err := client.ValidateToken(context.Background())
		if err != nil {
			t.Fatalf("ValidateToken returned error: %v", err)
		}
		if got == nil || *got != 123 {
			t.Fatalf("ValidateToken userID = %v, want 123", got)
		}
	})

	t.Run("invalid token result", func(t *testing.T) {
		setup()
		defer teardown()

		mux.HandleFunc("/v2/oauth2/validate", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			_, _ = w.Write([]byte(`{"result":false,"token_id":null,"token_scope":null,"user_id":null}`))
		})

		got, err := client.ValidateToken(context.Background())
		if got != nil {
			t.Fatalf("ValidateToken userID = %v, want nil", got)
		}
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("ValidateToken error = %v, want %v", err, ErrInvalidToken)
		}
	})

	t.Run("missing user id", func(t *testing.T) {
		setup()
		defer teardown()

		mux.HandleFunc("/v2/oauth2/validate", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			_, _ = w.Write([]byte(`{}`))
		})

		got, err := client.ValidateToken(context.Background())
		if got != nil {
			t.Fatalf("ValidateToken userID = %v, want nil", got)
		}
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("ValidateToken error = %v, want %v", err, ErrInvalidToken)
		}
	})
}
