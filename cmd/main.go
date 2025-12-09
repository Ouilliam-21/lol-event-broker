package main

import (
	"flag"
	"log"
	"michelprogram/lol-event/internal"
)

func main() {
	var apiEndpoint string
	var dropletEndpoint string

	flag.StringVar(&apiEndpoint, "Live client data", "https://127.0.0.1:2999/", "api endpoints to reach live client data")
	flag.StringVar(&dropletEndpoint, "Droplet URL", "https://127.0.0.1:2999/", "endpoint to send league of legend event")

	flag.Parse()

	queue := make(chan []byte)

	liveClient := internal.NewLiveClient(apiEndpoint, queue)
	droplet, err := internal.NewDroplet(dropletEndpoint, queue)

	log.Println("Object created for liveclient and droplet")

	if err != nil {
		log.Fatal(err)
	}

	go liveClient.Process()

	log.Println("Waiting to received message events...")

	for msg := range queue {
		droplet.Send(msg)
	}

}
