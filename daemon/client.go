package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/VinylStage/animux-like/pet"
)

func getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", getSocketPath())
			},
		},
	}
}

func SendRequest(path string, query map[string]string) (*pet.State, error) {
	client := getClient()
	
	reqURL := "http://unix" + path
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			q.Add(k, v)
		}
		reqURL += "?" + q.Encode()
	}

	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("is the daemon running? error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, pet.ErrNoPet
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("daemon error: status %d", resp.StatusCode)
	}

	var state pet.State
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}
	return &state, nil
}
