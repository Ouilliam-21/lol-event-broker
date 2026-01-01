package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func HttpPostRequest(client *http.Client, targetURL *url.URL, body io.Reader, token string) ([]byte, error) {
	return httpRequest(client, "POST", targetURL, body, token)
}

func HttpGetRequest(client *http.Client, targetURL *url.URL, token string) ([]byte, error) {
	return httpRequest(client, "GET", targetURL, nil, token)
}

func httpRequest(client *http.Client, method string, targetURL *url.URL, body io.Reader, token string) ([]byte, error) {
	request, err := http.NewRequest(method, targetURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to request at %s: %w", targetURL.String(), err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("url returned status %d: %s", resp.StatusCode, string(raw))
	}

	return raw, nil
}
