package main

import (
	"context"
	"flag"
	"log"
	"michelprogram/lol-event/internal"
	"michelprogram/lol-event/internal/database"
	"michelprogram/lol-event/internal/riot"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	var apiEndpoint string
	var dropletEndpoint string
	var wg sync.WaitGroup

	flag.StringVar(&apiEndpoint, "liveclient", "https://127.0.0.1:2999", "api endpoints to reach live client data")
	flag.StringVar(&dropletEndpoint, "droplet", "https://127.0.0.1:2999", "endpoint to send league of legend event")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventIds := make(chan []string, 100)

	conn, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to create Database: %v", err)
	}

	gameSessionRepository := database.NewGameSessionRepository(conn.Pool)
	riotEventRepository := database.NewRiotEventRepository(conn.Pool)

	liveClient, err := riot.NewLiveClient(apiEndpoint, eventIds, gameSessionRepository, riotEventRepository)
	if err != nil {
		log.Fatalf("Failed to create LiveClient: %v", err)
	}

	droplet, err := internal.NewDroplet(dropletEndpoint, eventIds)
	if err != nil {
		log.Fatalf("Failed to create Droplet: %v", err)
	}

	log.Println("Objects created for liveclient and droplet")

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := liveClient.Process(ctx); err != nil {
			log.Printf("LiveClient error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := droplet.SendEvents(ctx); err != nil {
			log.Printf("SendEvents error: %v", err)
		}
	}()

	log.Println("All services started, waiting for events...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, cleaning up...")

	cancel()

	close(eventIds)

	wg.Wait()
	log.Println("Shutdown complete")
}
