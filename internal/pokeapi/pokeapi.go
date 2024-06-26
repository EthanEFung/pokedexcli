package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://pokeapi.co/api/v2"

type Cache interface {
	Add(key string, value []byte)
	Get(key string) ([]byte, bool)
}

type Client struct {
	cache      Cache
	httpClient http.Client
}

func NewClient(cacheInterval time.Duration, cache Cache) Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		cache: cache,
	}
}

// GetJSON is a generic method to make a GET request to the pokeapi
// and unmarshal the response into the provided interface
// it also will return cached responses if a http request has already been made
func (c *Client) GetJSON(url string, v any) error {
	// check the cache
	if b, ok := c.cache.Get(url); ok {
		// fmt.Println("cached results")
		err := json.Unmarshal(b, v)
		if err != nil {
			return err
		}
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad status code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.cache.Add(url, b)

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
