package putio

import (
	"context"
	"net/http"
)

// EventsService is the service to gather information about user's events.
type EventsService struct {
	client *Client
}

// nolint:godox
// FIXME: events list returns inconsistent data structures.

// List gets list of dashboard events. It includes downloads and share events.
func (e *EventsService) List(ctx context.Context) ([]Event, error) {
	req, err := e.client.NewRequest(ctx, http.MethodGet, "/v2/events/list", nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Events []Event
	}
	_, err = e.client.Do(req, &r) // nolint:bodyclose
	if err != nil {
		return nil, err
	}
	return r.Events, nil
}

// Delete Clears all all dashboard events.
func (e *EventsService) Delete(ctx context.Context) error {
	req, err := e.client.NewRequest(ctx, http.MethodPost, "/v2/events/delete", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = e.client.Do(req, &struct{}{}) // nolint:bodyclose
	if err != nil {
		return err
	}

	return nil
}
