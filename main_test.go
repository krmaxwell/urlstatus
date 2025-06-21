package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadURLs(t *testing.T) {
	r := strings.NewReader("http://example.com\nhttp://foo.com\n")
	urls, err := readURLs(r)
	if err != nil {
		t.Fatalf("readURLs returned error: %v", err)
	}
	if len(urls) != 2 {
		t.Fatalf("expected 2 urls, got %d", len(urls))
	}
	if urls[0] != "http://example.com" || urls[1] != "http://foo.com" {
		t.Fatalf("unexpected urls: %v", urls)
	}
}

func TestCheckURLSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	r := checkURL(context.Background(), ts.Client(), ts.URL)
	if r.err != nil {
		t.Fatalf("checkURL returned error: %v", r.err)
	}
	if r.status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", r.status)
	}
}

func TestCheckURLRequestError(t *testing.T) {
	r := checkURL(context.Background(), http.DefaultClient, ":bad")
	if r.err == nil {
		t.Fatal("expected error for bad url")
	}
}
