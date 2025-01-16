package pocs

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

var ErrTokenRequired = errors.New("NEPLOX_TOKEN is required for this operation")

type Client struct {
	baseURL *url.URL
	token   string

	resty *resty.Client
}

func NewClient(baseURL *url.URL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		resty:   resty.New().SetBaseURL(baseURL.String()).SetHeader("X-Neplox", token),
	}
}

func (c *Client) Get(path string) (content []byte, mime string, err error) {
	resp, err := c.resty.R().Get(path)
	if err != nil {
		return nil, "", fmt.Errorf("request failure: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, "", fmt.Errorf("object not found")
	} else if resp.StatusCode() != http.StatusOK {
		return nil, "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), resp.Header().Get("Content-Type"), nil
}
