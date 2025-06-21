package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
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

var logger = log.New(os.Stderr, "", log.LstdFlags)

func main() {
	urls, err := readURLs(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "input error:", err)
		os.Exit(1)
	}

	logger.Printf("read %d urls", len(urls))

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
	logger.Printf("checking %s", url)
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Printf("request creation error for %s: %v", url, err)
		return result{url: url, err: err}
	}
	resp, err := client.Do(req)
	dur := time.Since(start)
	if err != nil {
		logger.Printf("request error for %s: %v", url, err)
		return result{url: url, duration: dur, err: err}
	}
	resp.Body.Close()
	logger.Printf("done %s status=%d duration=%s", url, resp.StatusCode, dur)
	return result{url: url, status: resp.StatusCode, duration: dur}
}
