package gcd

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
)

const (
	DefaultPrefix    = "https://www.comics.org/api"
	DefaultUserAgent = "GCD Client/Go 2024.11.13"
)

var defaultHTTPClient = &http.Client{Transport: &http2.Transport{}}

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type API struct {
	Prefix string
	Client HTTPDoer // override the client. Note that the gcd api only accept requests with HTTP/2, so http.DefaultClient is not compatible

	SessionID string // optional cookie value for gcdsessionid
}

func (a API) client() HTTPDoer {
	if a.Client != nil {
		return a.Client
	}

	return defaultHTTPClient
}

func (a API) req(ctx context.Context, url string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	if a.SessionID != "" {
		httpReq.AddCookie(&http.Cookie{
			Name:   "gcdsessionid",
			Value:  a.SessionID,
			Quoted: false,
			Path:   "/",
		})
	}

	httpReq.Header.Set("User-Agent", DefaultUserAgent)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Accept-Charset", "utf-8")

	resp, err := a.client().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	return resp, nil
}
