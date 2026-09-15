package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
)

type scheduleRepoPG struct {
	db *sql.DB
}

func NewScheduleRepositoryPG(db *sql.DB) repository.ScheduleRepository {
	return &scheduleRepoPG{db: db}
}

func (r *scheduleRepoPG) Create(ctx context.Context, s *entity.Schedule) error {
	query := `
		INSERT INTO schedules (movie_id, studio_id, start_time, end_time, price, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	if s.Status == "" {
		s.Status = entity.ScheduleScheduled
	}

	err := r.db.QueryRowContext(
		ctx, query,
		s.MovieID, s.StudioID, s.StartTime, s.EndTime, s.Price, s.Status,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert schedule: %w", err)
	}

	seedSeatsQuery := `
		INSERT INTO schedule_seats (schedule_id, seat_id, status)
		SELECT $1, id, 'AVAILABLE'
		FROM seats
		WHERE studio_id = $2 AND is_active = true
		ON CONFLICT DO NOTHING
	`
	_, _ = r.db.ExecContext(ctx, seedSeatsQuery, s.ID, s.StudioID)

	return nil
}

func (r *scheduleRepoPG) FindByID(ctx context.Context, id string) (*entity.ScheduleDetailResponse, error) {
	query := `
		SELECT 
			s.id, s.start_time, s.end_time, s.price, s.status, s.created_at, s.updated_at,
			m.id, m.title, m.duration_minutes, m.genre, m.age_rating, COALESCE(m.poster_url, ''),
			st.id, st.name, st.studio_type, st.total_seats,
			c.id, c.name, ci.name, c.address
		FROM schedules s
		JOIN movies m ON s.movie_id = m.id
		JOIN studios st ON s.studio_id = st.id
		JOIN cinemas c ON st.cinema_id = c.id
		JOIN cities ci ON c.city_id = ci.id
		WHERE s.id = $1
	`

	var item entity.ScheduleDetailResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.StartTime, &item.EndTime, &item.Price, &item.Status, &item.CreatedAt, &item.UpdatedAt,
		&item.Movie.ID, &item.Movie.Title, &item.Movie.DurationMinutes, &item.Movie.Genre, &item.Movie.AgeRating, &item.Movie.PosterURL,
		&item.Studio.ID, &item.Studio.Name, &item.Studio.StudioType, &item.Studio.TotalSeats,
		&item.Cinema.ID, &item.Cinema.Name, &item.Cinema.CityName, &item.Cinema.Address,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query schedule by id: %w", err)
	}

	return &item, nil
}

func (r *scheduleRepoPG) FindAll(ctx context.Context, filter entity.ScheduleFilter) ([]entity.ScheduleDetailResponse, int, error) {
	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if filter.MovieID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.movie_id = $%d", argIdx))
		args = append(args, filter.MovieID)
		argIdx++
	}
	if filter.StudioID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.studio_id = $%d", argIdx))
		args = append(args, filter.StudioID)
		argIdx++
	}
	if filter.CinemaID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("st.cinema_id = $%d", argIdx))
		args = append(args, filter.CinemaID)
		argIdx++
	}
	if filter.CityID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.city_id = $%d", argIdx))
		args = append(args, filter.CityID)
		argIdx++
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Date != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("DATE(s.start_time) = $%d", argIdx))
		args = append(args, filter.Date)
		argIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM schedules s
		JOIN movies m ON s.movie_id = m.id
		JOIN studios st ON s.studio_id = st.id
		JOIN cinemas c ON st.cinema_id = c.id
		WHERE %s
	`, whereSQL)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT 
			s.id, s.start_time, s.end_time, s.price, s.status, s.created_at, s.updated_at,
			m.id, m.title, m.duration_minutes, m.genre, m.age_rating, COALESCE(m.poster_url, ''),
			st.id, st.name, st.studio_type, st.total_seats,
			c.id, c.name, ci.name, c.address
		FROM schedules s
		JOIN movies m ON s.movie_id = m.id
		JOIN studios st ON s.studio_id = st.id
		JOIN cinemas c ON st.cinema_id = c.id
		JOIN cities ci ON c.city_id = ci.id
		WHERE %s
		ORDER BY s.start_time ASC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query schedules list: %w", err)
	}
	defer rows.Close()

	var results []entity.ScheduleDetailResponse
	for rows.Next() {
		var item entity.ScheduleDetailResponse
		if err := rows.Scan(
			&item.ID, &item.StartTime, &item.EndTime, &item.Price, &item.Status, &item.CreatedAt, &item.UpdatedAt,
			&item.Movie.ID, &item.Movie.Title, &item.Movie.DurationMinutes, &item.Movie.Genre, &item.Movie.AgeRating, &item.Movie.PosterURL,
			&item.Studio.ID, &item.Studio.Name, &item.Studio.StudioType, &item.Studio.TotalSeats,
			&item.Cinema.ID, &item.Cinema.Name, &item.Cinema.CityName, &item.Cinema.Address,
		); err != nil {
			return nil, 0, err
		}
		results = append(results, item)
	}

	return results, total, nil
}

func (r *scheduleRepoPG) Update(ctx context.Context, s *entity.Schedule) error {
	query := `
		UPDATE schedules
		SET movie_id = $1, studio_id = $2, start_time = $3, end_time = $4, price = $5, status = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		s.MovieID, s.StudioID, s.StartTime, s.EndTime, s.Price, s.Status, s.ID,
	).Scan(&s.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("schedule not found")
		}
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}

func (r *scheduleRepoPG) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM schedules WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("schedule not found")
	}
	return nil
}

func (r *scheduleRepoPG) HasTimeConflict(ctx context.Context, studioID string, startTime, endTime time.Time, excludeID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM schedules
		WHERE studio_id = $1
		  AND status != 'CANCELLED'
		  AND (start_time < $3 AND end_time > $2)
		  AND ($4 = '' OR id != $4::UUID)
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, studioID, startTime, endTime, excludeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check schedule conflict: %w", err)
	}
	return count > 0, nil
}
