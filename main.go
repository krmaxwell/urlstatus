package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type result struct {
	url      string
	status   int
	duration time.Duration
	err      error
}

func main() {
	urls, err := readURLs(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "input error:", err)
		os.Exit(1)
	}

	var results []result
	for _, u := range urls {
		r := checkURL(context.Background(), http.DefaultClient, u)
		results = append(results, r)
	}

	for _, r := range results {
		fmt.Printf("%s\tstatus=%d\tduration=%s\t%v\n", r.url, r.status, r.duration, r.err)
	}
}

func readURLs(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var urls []string
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func checkURL(ctx context.Context, client *http.Client, url string) result {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result{url: url, err: err}
	}
	resp, err := client.Do(req)
	dur := time.Since(start)
	if err != nil {
		return result{url: url, duration: dur, err: err}
	}
	resp.Body.Close()
	return result{url: url, status: resp.StatusCode, duration: dur}
}
