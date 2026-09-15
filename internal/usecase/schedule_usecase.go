package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
)

var (
	ErrScheduleNotFound  = errors.New("jadwal tayang tidak ditemukan")
	ErrScheduleConflict  = errors.New("jadwal tayang bentrok dengan penayangan lain di studio yang sama")
	ErrInvalidTimeRange  = errors.New("waktu selesai tayang harus lebih besar daripada waktu mulai tayang")
	ErrMovieNotFound     = errors.New("film tidak ditemukan atau tidak aktif")
	ErrStudioNotFound    = errors.New("studio bioskop tidak ditemukan atau tidak aktif")
	ErrInvalidTimeFormat = errors.New("format waktu tidak valid. Gunakan format ISO8601 (contoh: 2025-09-15T13:00:00Z atau 2025-09-15 13:00:00)")
)

type ScheduleUsecase interface {
	Create(ctx context.Context, req entity.CreateScheduleRequest) (*entity.ScheduleDetailResponse, error)
	GetByID(ctx context.Context, id string) (*entity.ScheduleDetailResponse, error)
	GetAll(ctx context.Context, filter entity.ScheduleFilter) ([]entity.ScheduleDetailResponse, int, error)
	Update(ctx context.Context, id string, req entity.UpdateScheduleRequest) (*entity.ScheduleDetailResponse, error)
	Delete(ctx context.Context, id string) error
}

type scheduleUsecase struct {
	scheduleRepo repository.ScheduleRepository
	movieRepo    repository.MovieRepository
	studioRepo   repository.StudioRepository
}

func NewScheduleUsecase(
	scheduleRepo repository.ScheduleRepository,
	movieRepo repository.MovieRepository,
	studioRepo repository.StudioRepository,
) ScheduleUsecase {
	return &scheduleUsecase{
		scheduleRepo: scheduleRepo,
		movieRepo:    movieRepo,
		studioRepo:   studioRepo,
	}
}

func (u *scheduleUsecase) Create(ctx context.Context, req entity.CreateScheduleRequest) (*entity.ScheduleDetailResponse, error) {
	movie, err := u.movieRepo.FindByID(ctx, req.MovieID)
	if err != nil {
		return nil, fmt.Errorf("error verifying movie: %w", err)
	}
	if movie == nil {
		return nil, ErrMovieNotFound
	}

	studio, err := u.studioRepo.FindByID(ctx, req.StudioID)
	if err != nil {
		return nil, fmt.Errorf("error verifying studio: %w", err)
	}
	if studio == nil {
		return nil, ErrStudioNotFound
	}

	startTime, err := parseDateTime(req.StartTime)
	if err != nil {
		return nil, ErrInvalidTimeFormat
	}

	var endTime time.Time
	if req.EndTime != "" {
		endTime, err = parseDateTime(req.EndTime)
		if err != nil {
			return nil, ErrInvalidTimeFormat
		}
	} else {
		cleaningBuffer := 15 * time.Minute
		endTime = startTime.Add(time.Duration(movie.DurationMinutes)*time.Minute + cleaningBuffer)
	}

	if !endTime.After(startTime) {
		return nil, ErrInvalidTimeRange
	}

	conflict, err := u.scheduleRepo.HasTimeConflict(ctx, req.StudioID, startTime, endTime, "")
	if err != nil {
		return nil, fmt.Errorf("error checking schedule conflict: %w", err)
	}
	if conflict {
		return nil, ErrScheduleConflict
	}

	schedule := &entity.Schedule{
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     req.Price,
		Status:    entity.ScheduleScheduled,
	}

	if err := u.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return u.scheduleRepo.FindByID(ctx, schedule.ID)
}

func (u *scheduleUsecase) GetByID(ctx context.Context, id string) (*entity.ScheduleDetailResponse, error) {
	item, err := u.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding schedule: %w", err)
	}
	if item == nil {
		return nil, ErrScheduleNotFound
	}
	return item, nil
}

func (u *scheduleUsecase) GetAll(ctx context.Context, filter entity.ScheduleFilter) ([]entity.ScheduleDetailResponse, int, error) {
	return u.scheduleRepo.FindAll(ctx, filter)
}

func (u *scheduleUsecase) Update(ctx context.Context, id string, req entity.UpdateScheduleRequest) (*entity.ScheduleDetailResponse, error) {
	existing, err := u.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrScheduleNotFound
	}

	updatedSchedule := &entity.Schedule{
		ID:        existing.ID,
		MovieID:   existing.Movie.ID,
		StudioID:  existing.Studio.ID,
		StartTime: existing.StartTime,
		EndTime:   existing.EndTime,
		Price:     existing.Price,
		Status:    existing.Status,
	}

	if req.MovieID != nil {
		movie, err := u.movieRepo.FindByID(ctx, *req.MovieID)
		if err != nil || movie == nil {
			return nil, ErrMovieNotFound
		}
		updatedSchedule.MovieID = *req.MovieID
	}

	if req.StudioID != nil {
		studio, err := u.studioRepo.FindByID(ctx, *req.StudioID)
		if err != nil || studio == nil {
			return nil, ErrStudioNotFound
		}
		updatedSchedule.StudioID = *req.StudioID
	}

	if req.StartTime != nil {
		st, err := parseDateTime(*req.StartTime)
		if err != nil {
			return nil, ErrInvalidTimeFormat
		}
		updatedSchedule.StartTime = st
	}

	if req.EndTime != nil {
		et, err := parseDateTime(*req.EndTime)
		if err != nil {
			return nil, ErrInvalidTimeFormat
		}
		updatedSchedule.EndTime = et
	}

	if !updatedSchedule.EndTime.After(updatedSchedule.StartTime) {
		return nil, ErrInvalidTimeRange
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, errors.New("harga tidak boleh negatif")
		}
		updatedSchedule.Price = *req.Price
	}

	if req.Status != nil {
		updatedSchedule.Status = *req.Status
	}

	conflict, err := u.scheduleRepo.HasTimeConflict(
		ctx, updatedSchedule.StudioID, updatedSchedule.StartTime, updatedSchedule.EndTime, id,
	)
	if err != nil {
		return nil, fmt.Errorf("error checking conflict on update: %w", err)
	}
	if conflict {
		return nil, ErrScheduleConflict
	}

	if err := u.scheduleRepo.Update(ctx, updatedSchedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return u.scheduleRepo.FindByID(ctx, id)
}

func (u *scheduleUsecase) Delete(ctx context.Context, id string) error {
	existing, err := u.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrScheduleNotFound
	}
	return u.scheduleRepo.Delete(ctx, id)
}

func parseDateTime(val string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("cannot parse time")
}
