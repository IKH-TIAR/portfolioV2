package uploader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrNotConfigured = errors.New("supabase storage is not configured")

type Client struct {
	baseURL   string
	apiKey    string
	bucket    string
	http      *http.Client
	publicURL string
}

func New(baseURL, apiKey, bucket string) (*Client, error) {
	if baseURL == "" || apiKey == "" || bucket == "" {
		return nil, ErrNotConfigured
	}
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		bucket:    bucket,
		publicURL: baseURL + "/storage/v1/object/public/" + bucket + "/",
		http:      &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) Upload(ctx context.Context, objectPath string, data []byte, contentType string) (string, error) {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.baseURL, c.bucket, objectPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("build upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload to supabase: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("upload to supabase: unexpected status %d", resp.StatusCode)
	}
	return c.publicURL + objectPath, nil
}
