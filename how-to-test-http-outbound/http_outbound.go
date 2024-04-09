package how_to_test_http_outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Post struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type JSONPlaceholderOutbound struct {
	client  *http.Client
	baseURL string
}

func NewJSONPlaceholderOutbound(client *http.Client, baseURL string) *JSONPlaceholderOutbound {
	return &JSONPlaceholderOutbound{
		client:  client,
		baseURL: baseURL,
	}
}

func (jp *JSONPlaceholderOutbound) GetPosts(ctx context.Context) ([]Post, error) {
	fullURL := jp.baseURL + "/posts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could't create request: %w", err)
	}

	var posts []Post
	if err := jp.fetch(req, &posts); err != nil {
		return nil, fmt.Errorf("could't fetch posts: %w", err)
	}
	return posts, nil
}

func (jp *JSONPlaceholderOutbound) GetPost(ctx context.Context, id int64) (*Post, error) {
	fullURL := jp.baseURL + "/posts/" + strconv.FormatInt(id, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could't create request: %w", err)
	}

	var post Post
	if err := jp.fetch(req, &post); err != nil {
		return nil, fmt.Errorf("could't fetch post: %w", err)
	}
	return &post, nil
}

func (jp *JSONPlaceholderOutbound) fetch(req *http.Request, dst any) error {
	resp, err := jp.client.Do(req)
	if err != nil {
		return fmt.Errorf("could't do request to %q: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could't read response body: %w", err)
	}

	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("could't unmarshal response body to JSON, raw body=%q: %w", raw, err)
	}

	return nil
}
