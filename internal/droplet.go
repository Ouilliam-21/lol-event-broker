package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Droplet struct {
	endpoint string
	queue    <-chan []byte
}

func NewDroplet(endpoint string, queue <-chan []byte) (*Droplet, error) {

	health := fmt.Sprintf("%s/health", endpoint)

	err := checkEndpoints(health)

	if err != nil {
		return nil, err
	}

	return &Droplet{
		endpoint: endpoint,
	}, nil
}

func checkEndpoints(url string) error {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call droplet api service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pricing service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (d Droplet) Send(data []byte) error {

	url := fmt.Sprintf("%s/generate", d.endpoint)

	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: time.Second * time.Duration(10)}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request at %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("url returned status %d: %s", resp.StatusCode, string(raw))
	}

	return nil
}
