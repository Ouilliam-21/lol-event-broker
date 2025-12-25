package riot

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	db "michelprogram/lol-event/internal/database"
	"michelprogram/lol-event/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type LiveClient struct {
	gameStatus            status
	endpointPlayers       *url.URL
	endpointEvents        *url.URL
	httpClient            *http.Client
	eventIds              chan<- []string
	gameSessionRepository db.IGameSessionRepository
	riotEventRepository   db.IRiotEventRepository
}

func NewLiveClient(endpoint string, eventIds chan<- []string, gameSessionRepository db.IGameSessionRepository, riotEventRepository db.IRiotEventRepository) (*LiveClient, error) {

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
		endpointEvents:        endpointEvents,
		gameStatus:            NotStarted,
		eventIds:              eventIds,
		httpClient:            httpClient,
		gameSessionRepository: gameSessionRepository,
		riotEventRepository:   riotEventRepository,
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

			eventsContainer, err := NewEventContainer(raw)
			if err != nil {
				log.Println("Can't unmarshal JSON:", err)
				continue
			}

			if len(eventsContainer.List.Items) <= 0 {
				continue
			}

			last, ok := eventsContainer.GetLast()

			if !ok || last.ID == eventID {
				continue
			}

			eventsContainer.FilterActiveEvents()

			ids := make([]string, 0, len(eventsContainer.List.Items))

			for _, event := range eventsContainer.List.Items {
				riotEvent := &db.RiotEvent{
					ID:            uuid.New().String(),
					GameSessionId: gameSessionID,
					RiotEventId:   event.ID,
					EventName:     string(event.Name),
					EventData:     event.Raw,
				}

				_, err = lc.riotEventRepository.CreateRiotEvent(ctx, riotEvent)
				if err != nil {
					log.Printf("Failed to create riot event: %v", err)
					continue
				}

				ids = append(ids, riotEvent.ID)
			}

			eventID = last.ID
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
			raw, err := utils.HttpGetRequest(lc.httpClient, lc.endpointPlayers)
			if err != nil && lc.gameStatus == NotStarted {
				log.Println("Game not started: couldn't reach live client endpoint")
				continue
			}

			now := time.Now()

			gameSessionItem := &db.GameSession{
				ID:         uuid.New().String(),
				RiotGameId: 0,
				Status:     db.GameStatusActive,
				StartedAt:  now,
				PlayerData: json.RawMessage(raw),
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			gameSession, err := lc.gameSessionRepository.CreateGameSession(ctx, gameSessionItem)

			log.Printf("Game session created: %d\n", gameSession.RiotGameId)
			if err != nil {
				return fmt.Errorf("failed to create game session: %w", err)
			}

			lc.poolGameEvent(ctx, gameSession.ID)

			gameSessionItem.Status = db.GameStatusPlayed
			gameSessionItem.EndedAt = time.Now()

			_, err = lc.gameSessionRepository.UpdateGameSession(ctx, gameSessionItem)
			if err != nil {
				return fmt.Errorf("failed to update game session: %w", err)
			}
		}

	}
}
