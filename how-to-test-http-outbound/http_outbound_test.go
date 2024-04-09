package how_to_test_http_outbound

import (
	"context"
	"net/http"
	"testing"

	"github.com/josestg/gotips/how-to-test-http-outbound/transporttest"
)

var (
	testDataPosts = []Post{
		{ID: 1, UserID: 1, Title: "title1", Body: "body1"},
		{ID: 2, UserID: 1, Title: "title2", Body: "body2"},
		{ID: 3, UserID: 2, Title: "title3", Body: "body3"},
	}
)

func TestJSONPlaceholderOutbound_GetPosts(t *testing.T) {
	client := transporttest.NewClient(
		transporttest.AssertHost(t, "example.com"),
		transporttest.AssertMethod(t, http.MethodGet),
		transporttest.AssertPath(t, "/posts"),
		transporttest.RespondJSON(testDataPosts, http.StatusOK),
	)

	jp := NewJSONPlaceholderOutbound(client, "https://example.com")
	posts, err := jp.GetPosts(context.Background())
	if err != nil {
		t.Fatalf("GetPosts failed: %v", err)
	}

	if len(posts) != len(testDataPosts) {
		t.Fatalf("GetPosts returned %d posts, want %d", len(posts), len(testDataPosts))
	}

	for i, p := range posts {
		if p != testDataPosts[i] {
			t.Errorf("GetPosts returned post %d: got %v, want %v", i, p, testDataPosts[i])
		}
	}
}

func TestJSONPlaceholderOutbound_GetPost(t *testing.T) {
	client := transporttest.NewClient(
		transporttest.AssertHost(t, "example.com"),
		transporttest.AssertMethod(t, http.MethodGet),
		transporttest.AssertPath(t, "/posts/1"),
		transporttest.RespondJSON(testDataPosts[0], http.StatusOK),
	)

	jp := NewJSONPlaceholderOutbound(client, "https://example.com")
	post, err := jp.GetPost(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPost failed: %v", err)
	}

	if *post != testDataPosts[0] {
		t.Errorf("GetPost returned: got %v, want %v", *post, testDataPosts[0])
	}
}
