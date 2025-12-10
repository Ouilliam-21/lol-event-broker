package riot

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"michelprogram/lol-event/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LiveClient struct {
	gameStatus      status
	endpointPlayers *url.URL
	endpointEvents  *url.URL
	httpClient      *http.Client
	events          chan<- []byte
	players         chan<- []byte
}

func NewLiveClient(endpoint string, events chan<- []byte, players chan<- []byte) (*LiveClient, error) {

	endpointPlayers, err := url.Parse(fmt.Sprintf("%s/liveclientdata/playerlist", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse players endpoint: %w", err)
	}

	endpointEvents, err := url.Parse(fmt.Sprintf("%s/liveclientdata/eventdata", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse events endpoint: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &LiveClient{
		endpointPlayers: endpointPlayers,
		endpointEvents:  endpointEvents,
		gameStatus:      NotStarted,
		events:          events,
		players:         players,
		httpClient:      httpClient,
	}, nil
}

func (lc *LiveClient) poolGameEvent() {
	var events Events
	var eventID int64

	log.Println("Game started")
	lc.gameStatus = Running

	for lc.gameStatus == Running {
		time.Sleep(5 * time.Second)

		eventURL := *lc.endpointEvents
		q := eventURL.Query()
		q.Set("eventID", strconv.FormatInt(eventID, 10))
		raw, err := utils.HttpGetRequest(lc.httpClient, &eventURL)

		if err != nil {
			log.Println("Game likely ended: ", err)
			lc.gameStatus = NotStarted
			return
		}

		if err := json.Unmarshal(raw, &events); err != nil {
			log.Println("Can't unmarshal JSON:", err)
			continue
		}

		if len(events.Events) <= 0 {
			continue
		}

		lastId := events.GetLast().ID

		if lastId == eventID {
			continue
		}

		eventsJSON, err := json.Marshal(events.FilterActiveEvents())
		if err != nil {
			log.Printf("Failed to marshal players: %v", err)
			continue
		}

		eventID = lastId
		lc.events <- eventsJSON
		log.Printf("Next event id %d", eventID)
	}
}

func (lc *LiveClient) Process() error {
	players := make([]Player, 0, 10)

	for {
		time.Sleep(5 * time.Second)

		raw, err := utils.HttpGetRequest(lc.httpClient, lc.endpointPlayers)
		if err != nil && lc.gameStatus == NotStarted {
			log.Println("Game not started: couldn't reach live client endpoint")
			continue
		}

		players = players[:0]

		if err := json.Unmarshal(raw, &players); err != nil {
			log.Println("Can't unmarshal JSON:", err)
			continue
		}

		playersJSON, err := json.Marshal(players)
		if err != nil {
			log.Printf("Failed to marshal players: %v", err)
			continue
		}

		lc.players <- playersJSON
		lc.poolGameEvent()
	}
}
