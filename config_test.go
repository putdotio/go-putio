package putio

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestConfig_WriteMethodsCloseResponseBody(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) error
	}{
		{
			name: "set all",
			call: func(client *Client) error {
				return client.Config.SetAll(context.Background(), map[string]string{"key": "value"})
			},
		},
		{
			name: "set",
			call: func(client *Client) error {
				return client.Config.Set(context.Background(), "key", "value")
			},
		},
		{
			name: "del",
			call: func(client *Client) error {
				return client.Config.Del(context.Background(), "key")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &trackedReadCloser{Reader: strings.NewReader(`{"status":"OK"}`)}
			client := NewClient(&http.Client{
				Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       body,
						Header:     make(http.Header),
					}, nil
				}),
			})

			if err := tt.call(client); err != nil {
				t.Fatal(err)
			}
			if !body.closed {
				t.Fatal("response body was not closed")
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type trackedReadCloser struct {
	*strings.Reader
	closed bool
}

func (r *trackedReadCloser) Close() error {
	r.closed = true
	return nil
}
