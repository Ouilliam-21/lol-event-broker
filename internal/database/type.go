package database

import (
	"context"
	"encoding/json"
	"time"
)

type GameStatus string

const (
	GameStatusActive GameStatus = "ACTIVE"
	GameStatusPlayed GameStatus = "PLAYED"
)

type GameSession struct {
	ID         string          `json:"id"`
	Status     GameStatus      `json:"status"`
	StartedAt  time.Time       `json:"startedAt"`
	EndedAt    time.Time       `json:"endedAt"`
	PlayerData json.RawMessage `json:"playerData"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type RiotEvent struct {
	ID            string          `json:"id"`
	GameSessionId string          `json:"gameSessionId"`
	RiotEventId   int64           `json:"riotEventId"`
	EventName     string          `json:"eventName"`
	EventData     json.RawMessage `json:"eventData"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type IGameSessionRepository interface {
	CreateGameSession(ctx context.Context, gameSession *GameSession) (*GameSession, error)
	UpdateGameSession(ctx context.Context, gameSession *GameSession) (*GameSession, error)
}

type IRiotEventRepository interface {
	CreateRiotEvent(ctx context.Context, riotEvent *RiotEvent) (*RiotEvent, error)
	UpdateRiotEvent(ctx context.Context, riotEvent *RiotEvent) (*RiotEvent, error)
}
