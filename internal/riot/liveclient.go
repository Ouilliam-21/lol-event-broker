package riot

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	db "michelprogram/lol-event/internal/database"
	"michelprogram/lol-event/internal/riot/events"
	"michelprogram/lol-event/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type status string

const (
	NotStarted status = "NOT STARTED"
	Running    status = "RUNNING"
)

type LiveClient struct {
	gameStatus            status
	endpointsPlayers      *url.URL
	endpointEvents        *url.URL
	httpClient            *http.Client
	eventIds              chan<- []string
	gameSessionRepository db.IGameSessionRepository
	riotEventRepository   db.IRiotEventRepository
	eventManager          *events.EventManager
}

func NewLiveClient(endpoint string, eventIds chan<- []string, gameSessionRepository db.IGameSessionRepository, riotEventRepository db.IRiotEventRepository) (*LiveClient, error) {

	endpointEvents, err := url.Parse(fmt.Sprintf("%s/liveclientdata/eventdata", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse events endpoint: %w", err)
	}

	endpointsPlayers, err := url.Parse(fmt.Sprintf("%s/liveclientdata/playerlist", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse players endpoint: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &LiveClient{
		endpointEvents:        endpointEvents,
		endpointsPlayers:      endpointsPlayers,
		gameStatus:            NotStarted,
		eventIds:              eventIds,
		httpClient:            httpClient,
		gameSessionRepository: gameSessionRepository,
		riotEventRepository:   riotEventRepository,
		eventManager:          events.NewEventManager(),
	}, nil
}

func (lc *LiveClient) poolGameEvent(ctx context.Context, gameSessionID string) {
	var eventID int64

	log.Println("Game started")
	lc.gameStatus = Running

	for lc.gameStatus == Running {
		select {
		case <-ctx.Done():
			log.Println("poolGameEvent shutting down")
			lc.gameStatus = NotStarted
			return
		case <-time.After(5 * time.Second):
			q := lc.endpointEvents.Query()
			q.Set("eventID", strconv.FormatInt(eventID, 10))
			lc.endpointEvents.RawQuery = q.Encode()

			raw, err := utils.HttpGetRequest(lc.httpClient, lc.endpointEvents)

			if err != nil {
				log.Println("Game likely ended: ", err)
				lc.gameStatus = NotStarted
				return
			}

			err = lc.eventManager.ProcessEvent(raw)
			if err != nil {
				log.Println("Can't process event: ", err)
				continue
			}

			if lc.eventManager.IsEmpty() {
				continue
			}

			lastEvent := lc.eventManager.GetLast()

			if lastEvent == nil || lastEvent.GetEventID() == eventID {
				continue
			}

			events := lc.eventManager.FilterEvents()

			ids := make([]string, 0, len(events))

			for _, event := range events {

				rawEvent, err := event.ToJson()
				if err != nil {
					log.Printf("Failed to marshal event: %v", err)
					continue
				}

				riotEvent := &db.RiotEvent{
					ID:            uuid.New().String(),
					GameSessionId: gameSessionID,
					RiotEventId:   event.GetEventID(),
					EventName:     string(event.GetEventName()),
					EventData:     rawEvent,
				}

				_, err = lc.riotEventRepository.CreateRiotEvent(ctx, riotEvent)
				if err != nil {
					log.Printf("Failed to create riot event: %v", err)
					continue
				}

				ids = append(ids, riotEvent.ID)
			}

			eventID = lastEvent.GetEventID()
			lc.eventIds <- ids
			log.Printf("Next event id %d", eventID)
		}
	}
}

func (lc *LiveClient) Process(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			log.Println("Process shutting down")
			return ctx.Err()
		case <-time.After(5 * time.Second):
			raw, err := utils.HttpGetRequest(lc.httpClient, lc.endpointsPlayers)
			if err != nil && lc.gameStatus == NotStarted {
				log.Println("Game not started: couldn't reach live client endpoint")
				continue
			}

			now := time.Now()

			gameSessionItem := &db.GameSession{
				ID:         uuid.New().String(),
				Status:     db.GameStatusActive,
				StartedAt:  now,
				PlayerData: json.RawMessage(raw),
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			_, err = lc.gameSessionRepository.CreateGameSession(ctx, gameSessionItem)
			if err != nil {
				return fmt.Errorf("failed to create game session: %w", err)
			}
			log.Printf("Game session created: %s\n", gameSessionItem.ID)

			lc.poolGameEvent(ctx, gameSessionItem.ID)

			gameSessionItem.Status = db.GameStatusPlayed
			gameSessionItem.EndedAt = time.Now()

			_, err = lc.gameSessionRepository.UpdateGameSession(ctx, gameSessionItem)
			if err != nil {
				return fmt.Errorf("failed to update game session: %w", err)
			}
		}

	}
}
