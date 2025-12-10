package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"michelprogram/lol-event/internal/utils"
	"net/http"
	"net/url"
	"time"
)

type Droplet struct {
	endpointPlayers *url.URL
	endpointEvents  *url.URL
	httpClient      *http.Client
	events          <-chan []byte
	players         <-chan []byte
}

func NewDroplet(endpoint string, events <-chan []byte, players <-chan []byte) (*Droplet, error) {
	endpointPlayers, err := url.Parse(fmt.Sprintf("%s/players", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse players endpoint: %w", err)
	}

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
		endpointPlayers: endpointPlayers,
		endpointEvents:  endpointEvents,
		events:          events,
		players:         players,
		httpClient:      httpClient,
	}, nil
}

func (d Droplet) SendEvents(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			log.Println("SendEvents shutting down")
			return ctx.Err()
		case msg, ok := <-d.events:
			if !ok {
				log.Println("Events channel closed")
				return nil
			}
			_, err := utils.HttpPostRequest(d.httpClient, d.endpointEvents, bytes.NewReader(msg))
			if err != nil {
				log.Printf("Failed to send event: %v", err)
				continue
			}
			log.Println("Event sent successfully")
		}
	}
}

func (d Droplet) SendPlayers(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			log.Println("SendPlayers shutting down")
			return ctx.Err()
		case msg, ok := <-d.players:
			if !ok {
				log.Println("Players channel closed")
				return nil
			}
			_, err := utils.HttpPostRequest(d.httpClient, d.endpointPlayers, bytes.NewReader(msg))
			if err != nil {
				log.Printf("Failed to send players: %v", err)
				continue
			}
			log.Println("Players sent successfully")
		}
	}
}
