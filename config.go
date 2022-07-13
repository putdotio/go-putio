package putio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ConfigService represents configuration related operations.
type ConfigService struct {
	client *Client
}

// GetAll all fills config.
func (f *ConfigService) GetAll(ctx context.Context, config interface{}) error {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config", nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	var r struct {
		Config json.RawMessage `json:"config"`
	}
	_, err = f.client.Do(req, &r) // nolint:bodyclose
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return json.Unmarshal(r.Config, &config) // nolint:wrapcheck
}

// Get fetches config item via given key.
func (f *ConfigService) Get(ctx context.Context, key string, value interface{}) (found bool, err error) {
	req, err := f.client.NewRequest(ctx, http.MethodGet, "/v2/config/"+key, nil)
	if err != nil {
		return false, err
	}
	var r struct {
		Value *json.RawMessage `json:"value"`
	}
	_, err = f.client.Do(req, &r) // nolint:bodyclose
	if err != nil {
		return false, err
	}
	if r.Value == nil {
		return false, nil
	}
	return true, json.Unmarshal(*r.Value, &value) // nolint:wrapcheck
}

// SetAll updates all config items.
func (f *ConfigService) SetAll(ctx context.Context, config interface{}) error {
	b, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	v := struct {
		Config json.RawMessage `json:"config"`
	}{
		Config: b,
	}
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req, err := f.client.NewRequest(ctx, http.MethodPut, "/v2/config", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Set("content-type", "application/json")
	_, err = f.client.Do(req, nil) // nolint:bodyclose
	return fmt.Errorf("%w", err)
}

// Set updates given config key's value.
func (f *ConfigService) Set(ctx context.Context, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	v := struct {
		Value json.RawMessage `json:"value"`
	}{
		Value: b,
	}
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req, err := f.client.NewRequest(ctx, http.MethodPut, "/v2/config/"+key, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Set("content-type", "application/json")
	_, err = f.client.Do(req, nil) // nolint:bodyclose
	return fmt.Errorf("%w", err)
}

// Del destroys config item via given key.
func (f *ConfigService) Del(ctx context.Context, key string) error {
	req, err := f.client.NewRequest(ctx, http.MethodDelete, "/v2/config/"+key, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	_, err = f.client.Do(req, nil) // nolint:bodyclose
	return fmt.Errorf("%w", err)
}
