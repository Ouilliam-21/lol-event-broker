package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GameSessionRepository struct {
	db *pgxpool.Pool
}

func NewGameSessionRepository(db *pgxpool.Pool) *GameSessionRepository {
	return &GameSessionRepository{db: db}
}

func (r *GameSessionRepository) CreateGameSession(ctx context.Context, gameSession *GameSession) (*GameSession, error) {

	var insertedID string
	err := r.db.QueryRow(ctx, "INSERT INTO game_sessions (id, status, started_at, ended_at, player_data) VALUES ($1, $2, $3, $4, $5) RETURNING id", gameSession.ID, gameSession.Status, gameSession.StartedAt, gameSession.EndedAt, gameSession.PlayerData).Scan(&insertedID)
	if err != nil {
		return nil, fmt.Errorf("failed to create game session: %w", err)
	}

	return gameSession, nil
}

func (r *GameSessionRepository) UpdateGameSession(ctx context.Context, gameSession *GameSession) (*GameSession, error) {

	var updatedID string
	err := r.db.QueryRow(ctx, "UPDATE game_sessions SET status = $1, started_at = $2, ended_at = $3, player_data = $4 WHERE id = $5 RETURNING id", gameSession.Status, gameSession.StartedAt, gameSession.EndedAt, gameSession.PlayerData, gameSession.ID).Scan(&updatedID)
	if err != nil {
		return nil, fmt.Errorf("failed to update game session: %w", err)
	}

	return gameSession, nil
}
