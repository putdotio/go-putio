package putio

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors.
var (
	ErrEmptyFolderName          = errors.New("empty folder name")
	ErrNoFileIDIsGiven          = errors.New("no file id is given")
	ErrFilenameCanNotBeEmpty    = errors.New("filename cannot be empty")
	ErrNewFilenameCanNotBeEmpty = errors.New("new filename cannot be empty")
	ErrInvalidPageNumber        = errors.New("invalid page number")
	ErrNoQueryGiven             = errors.New("no query given")
	ErrEmptySubtileKey          = errors.New("empty subtitle key is given")
	ErrNegativeTimeValue        = errors.New("time cannot be negative")
	ErrNoFileIsGiven            = errors.New("no files given")
	ErrEmptyUserName            = errors.New("empty username")
	ErrEmptyURL                 = errors.New("empty URL")
	ErrUnexpected               = errors.New("unexpected error")
)

// ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// Original http.Response
	Response *http.Response `json:"-"`

	// Body read from Response
	Body []byte `json:"-"`

	// Error while parsing the response
	ParseError error

	// These fileds are parsed from response if JSON.
	Message string `json:"error_message"`
	Type    string `json:"error_type"`
}

// nolint:goerr113
func (e *ErrorResponse) Error() string {
	if e.ParseError != nil {
		return fmt.Errorf(
			"cannot parse response. code:%d error:%w body:%q",
			e.Response.StatusCode,
			e.ParseError,
			string(e.Body[:250]),
		).Error()
	}
	return fmt.Sprintf(
		"putio error. code:%d type:%q message:%q request:%v %v",
		e.Response.StatusCode,
		e.Type,
		e.Message,
		e.Response.Request.Method,
		e.Response.Request.URL,
	)
}
