package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	username, password string
	httpClient         *http.Client
}

func New(username, password string, timeout time.Duration) *Client {
	return &Client{
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) MakeGETRequest(ctx context.Context, url string, queryString map[string]string, data interface{}) error {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("can't create request %q: %w", url, err)
	}

	if queryString != nil {
		q := req.URL.Query()
		for k, v := range queryString {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed %q: %s", url, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("can't read json %q: %w", url, err)
	}
	return nil
}

type KeyGetter interface {
	GetByKey(ctx context.Context, key string) (interface{}, error)
}

func (c *Client) BulkKeyRequests(ctx context.Context, getter KeyGetter, keys []string) ([]interface{}, error) {
	var wg sync.WaitGroup
	var errs []error
	var results []interface{}

	wg.Add(len(keys))
	var mu sync.Mutex
	setResult := func(data interface{}, err error) {
		mu.Lock()
		if err != nil {
			errs = append(errs, err)
		} else {
			results = append(results, data)
		}
		mu.Unlock()
	}
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					setResult(nil, fmt.Errorf("panic while getting %s: %v", key, err))
				}
			}()
			setResult(getter.GetByKey(ctx, key))
		}(key)
	}
	wg.Wait()
	var err error
	if len(errs) != 0 {
		err = fmt.Errorf("can't get data: %v", errs)
	}
	return results, err
}
