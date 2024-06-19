package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IClient interface {
	GetData(ctx context.Context, endpoint string) (*ClientResponse, error)
}

type client struct {
	httpClient *http.Client
}

func closeChannels(channels ...interface{}) {
	for _, ch := range channels {
		switch ch := ch.(type) {
		case chan error:
			close(ch)
		case chan map[string]interface{}:
			close(ch)
		}
	}
}

func NewClient() IClient {
	return &client{
		httpClient: &http.Client{},
	}
}

func (c *client) GetData(ctx context.Context, url string) (*ClientResponse, error) {
	result := make(chan map[string]interface{})
	errCh := make(chan error, 1)

	go func() {
		defer closeChannels(result, errCh)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errCh <- fmt.Errorf("error request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var jsonData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
			errCh <- fmt.Errorf("error decoding JSON: %w", err)
		}
		time.Sleep(2 * time.Second)
		result <- jsonData
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	case hasResult := <-result:
		response := &ClientResponse{hasResult}
		return response, nil
	}
}
