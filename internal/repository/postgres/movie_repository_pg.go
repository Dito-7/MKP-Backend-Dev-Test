package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
)

type movieRepoPG struct {
	db *sql.DB
}

func NewMovieRepositoryPG(db *sql.DB) repository.MovieRepository {
	return &movieRepoPG{db: db}
}

func (r *movieRepoPG) FindByID(ctx context.Context, id string) (*entity.Movie, error) {
	query := `
		SELECT id, title, description, duration_minutes, genre, age_rating, COALESCE(poster_url, ''), release_date, is_active, created_at, updated_at
		FROM movies
		WHERE id = $1 AND is_active = true
	`
	var m entity.Movie
	var releaseDate timeFormatDate
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Title, &m.Description, &m.DurationMinutes, &m.Genre, &m.AgeRating, &m.PosterURL, &releaseDate, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query movie by id: %w", err)
	}
	m.ReleaseDate = releaseDate.String()
	return &m, nil
}

func (r *movieRepoPG) FindAll(ctx context.Context) ([]entity.Movie, error) {
	query := `
		SELECT id, title, description, duration_minutes, genre, age_rating, COALESCE(poster_url, ''), release_date, is_active, created_at, updated_at
		FROM movies
		WHERE is_active = true
		ORDER BY release_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list movies: %w", err)
	}
	defer rows.Close()

	var movies []entity.Movie
	for rows.Next() {
		var m entity.Movie
		var releaseDate timeFormatDate
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.DurationMinutes, &m.Genre, &m.AgeRating, &m.PosterURL, &releaseDate, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.ReleaseDate = releaseDate.String()
		movies = append(movies, m)
	}
	return movies, nil
}

type timeFormatDate struct {
	sql.NullTime
}

func (t timeFormatDate) String() string {
	if t.Valid {
		return t.Time.Format("2006-01-02")
	}
	return ""
}
