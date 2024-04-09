package how_to_test_http_outbound

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var (
	testServer    *httptest.Server
	testDataPosts = []Post{
		{ID: 1, UserID: 1, Title: "title1", Body: "body1"},
		{ID: 2, UserID: 1, Title: "title2", Body: "body2"},
		{ID: 3, UserID: 2, Title: "title3", Body: "body3"},
	}
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(testDataPosts)
	})

	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		for _, p := range testDataPosts {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(p)
				return
			}
		}

		http.Error(w, "post not found", http.StatusNotFound)
	})

	testServer = httptest.NewServer(mux)

	exitCode := m.Run()
	testServer.Close()
	os.Exit(exitCode)
}

func TestJSONPlaceholderOutbound_GetPosts(t *testing.T) {
	client := testServer.Client()
	jp := NewJSONPlaceholderOutbound(client, testServer.URL)

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
	client := testServer.Client()
	jp := NewJSONPlaceholderOutbound(client, testServer.URL)

	for _, want := range testDataPosts {
		got, err := jp.GetPost(context.Background(), want.ID)
		if err != nil {
			t.Fatalf("GetPost failed: %v", err)
		}

		if *got != want {
			t.Errorf("GetPost returned post %d: got %v, want %v", want.ID, got, want)
		}
	}
}
