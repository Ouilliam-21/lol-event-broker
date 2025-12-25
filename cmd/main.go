package main

import (
	"context"
	"flag"
	"log"
	"michelprogram/lol-event/internal"
	"michelprogram/lol-event/internal/config"
	"michelprogram/lol-event/internal/database"
	"michelprogram/lol-event/internal/riot"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var configFile string
	var wg sync.WaitGroup

	flag.StringVar(&configFile, "config", "config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventIds := make(chan []string, 100)

	conn, err := database.NewDatabase(cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}

	gameSessionRepository := database.NewGameSessionRepository(conn.Pool)
	riotEventRepository := database.NewRiotEventRepository(conn.Pool)

	liveClient, err := riot.NewLiveClient(cfg.Endpoints.LiveClient, eventIds, gameSessionRepository, riotEventRepository, cfg.GetWatchedPlayers(), cfg.GetWatchedEvents())
	if err != nil {
		log.Fatalf("Failed to create LiveClient: %v", err)
	}

	droplet, err := internal.NewDroplet(cfg.Endpoints.Droplet, eventIds)
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
