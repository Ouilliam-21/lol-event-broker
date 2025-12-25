package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"michelprogram/lol-event/internal/utils"
	"net/http"
	"net/url"
	"time"
)

type Droplet struct {
	endpointEvents *url.URL
	httpClient     *http.Client
	eventIds       <-chan []string
}

type EventIds struct {
	EventIds []string `json:"eventIds"`
}

func NewDroplet(endpoint string, eventIds <-chan []string) (*Droplet, error) {

	endpointEvents, err := url.Parse(fmt.Sprintf("%s/events", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse events endpoint: %w", err)
	}

	endpointHealth, err := url.Parse(fmt.Sprintf("%s/health", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse health endpoint: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	_, err = utils.HttpGetRequest(httpClient, endpointHealth)
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	log.Println("Droplet health check passed")

	return &Droplet{
		endpointEvents: endpointEvents,
		eventIds:       eventIds,
		httpClient:     httpClient,
	}, nil
}

func (d Droplet) SendEvents(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			log.Println("SendEvents shutting down")
			return ctx.Err()
		case ids, ok := <-d.eventIds:
			if !ok {
				log.Println("Events channel closed")
				return nil
			}
			if len(ids) == 0 {
				continue
			}
			payload, err := json.Marshal(EventIds{EventIds: ids})
			if err != nil {
				log.Printf("Failed to marshal event ids: %v", err)
				continue
			}
			log.Printf("Sending event ids: %s", string(payload))
			_, err = utils.HttpPostRequest(d.httpClient, d.endpointEvents, bytes.NewReader(payload))
			if err != nil {
				log.Printf("Failed to send event: %v", err)
				continue
			}
			log.Println("Event sent successfully")
		}
	}
}
