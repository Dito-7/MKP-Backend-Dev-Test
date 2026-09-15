package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
)

type studioRepoPG struct {
	db *sql.DB
}

func NewStudioRepositoryPG(db *sql.DB) repository.StudioRepository {
	return &studioRepoPG{db: db}
}

func (r *studioRepoPG) FindByID(ctx context.Context, id string) (*entity.Studio, error) {
	query := `
		SELECT s.id, s.cinema_id, c.name as cinema_name, ci.name as city_name, s.name, s.studio_type, s.total_rows, s.total_cols, s.total_seats, s.is_active, s.created_at, s.updated_at
		FROM studios s
		JOIN cinemas c ON s.cinema_id = c.id
		JOIN cities ci ON c.city_id = ci.id
		WHERE s.id = $1 AND s.is_active = true
	`
	var s entity.Studio
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.CinemaID, &s.CinemaName, &s.CityName, &s.Name, &s.StudioType, &s.TotalRows, &s.TotalCols, &s.TotalSeats, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query studio by id: %w", err)
	}
	return &s, nil
}

func (r *studioRepoPG) FindAll(ctx context.Context, cinemaID string) ([]entity.Studio, error) {
	query := `
		SELECT s.id, s.cinema_id, c.name as cinema_name, ci.name as city_name, s.name, s.studio_type, s.total_rows, s.total_cols, s.total_seats, s.is_active, s.created_at, s.updated_at
		FROM studios s
		JOIN cinemas c ON s.cinema_id = c.id
		JOIN cities ci ON c.city_id = ci.id
		WHERE ($1 = '' OR s.cinema_id::text = $1) AND s.is_active = true
		ORDER BY s.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, cinemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to query studios: %w", err)
	}
	defer rows.Close()

	var studios []entity.Studio
	for rows.Next() {
		var s entity.Studio
		if err := rows.Scan(
			&s.ID, &s.CinemaID, &s.CinemaName, &s.CityName, &s.Name, &s.StudioType, &s.TotalRows, &s.TotalCols, &s.TotalSeats, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		studios = append(studios, s)
	}
	return studios, nil
}
