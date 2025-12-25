package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RiotEventRepository struct {
	db *pgxpool.Pool
}

func NewRiotEventRepository(db *pgxpool.Pool) *RiotEventRepository {
	return &RiotEventRepository{db: db}
}

func (r *RiotEventRepository) CreateRiotEvent(ctx context.Context, riotEvent *RiotEvent) (*RiotEvent, error) {

	var insertedID string
	err := r.db.QueryRow(ctx, "INSERT INTO riot_events (id, game_session_id, riot_event_id, event_name, event_data, received_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", riotEvent.ID, riotEvent.GameSessionId, riotEvent.RiotEventId, riotEvent.EventName, riotEvent.EventData, riotEvent.ReceivedAt).Scan(&insertedID)
	if err != nil {
		return nil, fmt.Errorf("failed to create riot event: %w", err)
	}

	return riotEvent, nil
}

func (r *RiotEventRepository) UpdateRiotEvent(ctx context.Context, riotEvent *RiotEvent) (*RiotEvent, error) {

	var updatedID string
	err := r.db.QueryRow(ctx, "UPDATE riot_events SET game_session_id = $1, riot_event_id = $2, event_name = $3, event_data = $4, received_at = $5 WHERE id = $6 RETURNING id", riotEvent.GameSessionId, riotEvent.RiotEventId, riotEvent.EventName, riotEvent.EventData, riotEvent.ReceivedAt, riotEvent.ID).Scan(&updatedID)
	if err != nil {
		return nil, fmt.Errorf("failed to update riot event: %w", err)
	}

	return riotEvent, nil
}
