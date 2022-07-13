package putio

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UploadService uses TUS (resumable upload protocol) for sending files to put.io.
type UploadService struct {
	// Log is a user supplied function to collect log messages from upload methods.
	Log func(message string)

	client *Client
}

func (u *UploadService) log(message string) {
	if u.Log != nil {
		u.Log(message)
	}
}

// CreateUpload is used for beginning new upload. Use returned location in SendFile function.
func (u *UploadService) CreateUpload(
	ctx context.Context,
	filename string,
	parentID, length int64,
	overwrite bool,
) (location string, err error) {
	u.log(fmt.Sprintf("Creating upload %q at parent=%d", filename, parentID))
	req, err := u.client.NewRequest(ctx, http.MethodPost, "$upload-tus$", nil)
	if err != nil {
		return
	}
	metadata := map[string]string{
		"name":       filename,
		"parent_id":  strconv.FormatInt(parentID, 10),
		"no-torrent": "true",
		"overwrite":  strconv.FormatBool(overwrite),
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(length, 10))
	req.Header.Set("Upload-Metadata", encodeMetadata(metadata))

	resp, err := u.client.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("%w status: %d", ErrUnexpected, resp.StatusCode)
		return
	}
	location = resp.Header.Get("Location")
	return
}

// SendFile sends the contents of the file to put.io.
// In case of an transmission error, you can resume upload but you have to get the correct offset from server by
// calling GetOffset and must seek to the new offset on io.Reader.
func (u *UploadService) SendFile(
	ctx context.Context,
	r io.Reader,
	location string,
	offset int64,
) (fileID int64, crc32 string, err error) {
	u.log(fmt.Sprintf("Sending file %q offset=%d", location, offset))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stop upload if speed is too slow.
	// Wrap reader so each read call resets the timer that cancels the request on certain duration.
	if u.client.Timeout > 0 {
		r = &timerResetReader{r: r, timer: time.AfterFunc(u.client.Timeout, cancel), timeout: u.client.Timeout}
	}

	req, err := u.client.NewRequest(ctx, http.MethodPatch, location, r)
	if err != nil {
		return
	}

	// putio.Client.NewRequests add another context for handling Client.Timeout. Replace it with original context.
	// Request must not be cancelled on timeout because sending upload body takes a long time.
	// We will be using timerResetReader for tracking uploaded bytes and doint cancellation there.
	if u.client.Timeout > 0 {
		req = req.WithContext(ctx)
	}

	req.Header.Set("content-type", "application/offset+octet-stream")
	req.Header.Set("upload-offset", strconv.FormatInt(offset, 10))
	resp, err := u.client.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("%w status: %d", ErrUnexpected, resp.StatusCode)
		return
	}
	fileID, err = strconv.ParseInt(resp.Header.Get("putio-file-id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse putio-file-id header: %w", err)
		return
	}
	crc32 = resp.Header.Get("putio-file-crc32")
	return
}

// GetOffset returns the offset at the server.
func (u *UploadService) GetOffset(ctx context.Context, location string) (n int64, err error) {
	u.log(fmt.Sprintf("Getting upload offset %q", location))
	req, err := u.client.NewRequest(ctx, http.MethodHead, location, nil)
	if err != nil {
		return
	}

	resp, err := u.client.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%w status: %d", ErrUnexpected, resp.StatusCode)
		return
	}
	n, err = strconv.ParseInt(resp.Header.Get("upload-offset"), 10, 64)
	u.log(fmt.Sprintln("uploadJob offset:", n))
	return n, fmt.Errorf("%w", err)
}

// TerminateUpload removes incomplete file from the server.
func (u *UploadService) TerminateUpload(ctx context.Context, location string) (err error) {
	u.log(fmt.Sprintf("Terminating upload %q", location))
	req, err := u.client.NewRequest(ctx, http.MethodDelete, location, nil)
	if err != nil {
		return
	}

	resp, err := u.client.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	u.log(fmt.Sprintln("Status code:", resp.StatusCode))
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("%w status: %d", ErrUnexpected, resp.StatusCode)
		return
	}
	return nil
}

func encodeMetadata(metadata map[string]string) string {
	encoded := make([]string, 0, len(metadata))
	for k, v := range metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, base64.StdEncoding.EncodeToString([]byte(v))))
	}
	return strings.Join(encoded, ",")
}

type timerResetReader struct {
	r       io.Reader
	timer   *time.Timer
	timeout time.Duration
}

func (r *timerResetReader) Read(p []byte) (int, error) {
	r.timer.Reset(r.timeout)
	return r.r.Read(p) // nolint:wrapcheck
}
