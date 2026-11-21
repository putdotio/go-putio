package putio

import (
	"context"
	"net/http"
)

// AccountService is the service to gather information about user account.
type AccountService struct {
	client *Client
}

// Info retrieves user account information.
func (a *AccountService) Info(ctx context.Context) (AccountInfo, error) {
	return a.info(ctx, "/v2/account/info")
}

// DownloadToken retrieves the account download token. Media requests, such as
// HLS playlists, authenticate with this token instead of the OAuth token.
func (a *AccountService) DownloadToken(ctx context.Context) (string, error) {
	info, err := a.info(ctx, "/v2/account/info?download_token=1")
	if err != nil {
		return "", err
	}
	if info.DownloadToken == "" {
		return "", ErrMissingDownloadToken
	}
	return info.DownloadToken, nil
}

func (a *AccountService) info(ctx context.Context, relURL string) (AccountInfo, error) {
	req, err := a.client.NewRequest(ctx, http.MethodGet, relURL, nil)
	if err != nil {
		return AccountInfo{}, err
	}

	var r struct {
		Info AccountInfo
	}
	_, err = a.client.Do(req, &r) // nolint:bodyclose
	if err != nil {
		return AccountInfo{}, err
	}
	return r.Info, nil
}

// Settings retrieves user preferences.
func (a *AccountService) Settings(ctx context.Context) (Settings, error) {
	req, err := a.client.NewRequest(ctx, http.MethodGet, "/v2/account/settings", nil)
	if err != nil {
		return Settings{}, err
	}
	var r struct {
		Settings Settings
	}
	_, err = a.client.Do(req, &r) // nolint:bodyclose
	if err != nil {
		return Settings{}, err
	}

	return r.Settings, nil
}
