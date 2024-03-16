package pokeapi

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
)

func (c *Client) GetSprite(spriteURL string) (image.Image, error) {
	var resourceResponse image.Image
	// check the cache
	if b, ok := c.cache.Get(spriteURL); ok {

		img, format, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return resourceResponse, err
		}
		if format != "png" {
			return resourceResponse, fmt.Errorf("unsupported image format: %v", format)
		}
		return img, nil
	}
	req, err := http.NewRequest("GET", spriteURL, nil)
	if err != nil {
		return resourceResponse, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resourceResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return resourceResponse, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return resourceResponse, err
	}
	c.cache.Add(spriteURL, b)

	img, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return resourceResponse, err
	}
	if format != "png" {
		return resourceResponse, fmt.Errorf("unsupported image format: %v", format)
	}
	return img, nil
}
