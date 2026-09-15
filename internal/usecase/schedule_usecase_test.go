package usecase_test

import (
	"context"
	"testing"
	"time"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/usecase"
)

type mockMovieRepo struct {
	movies map[string]*entity.Movie
}

func (m *mockMovieRepo) FindByID(ctx context.Context, id string) (*entity.Movie, error) {
	return m.movies[id], nil
}
func (m *mockMovieRepo) FindAll(ctx context.Context) ([]entity.Movie, error) {
	return nil, nil
}

type mockStudioRepo struct {
	studios map[string]*entity.Studio
}

func (m *mockStudioRepo) FindByID(ctx context.Context, id string) (*entity.Studio, error) {
	return m.studios[id], nil
}
func (m *mockStudioRepo) FindAll(ctx context.Context, cinemaID string) ([]entity.Studio, error) {
	return nil, nil
}

type mockScheduleRepo struct {
	schedules   map[string]*entity.Schedule
	hasConflict bool
}

func (m *mockScheduleRepo) Create(ctx context.Context, s *entity.Schedule) error {
	s.ID = "new-sched-id"
	m.schedules[s.ID] = s
	return nil
}
func (m *mockScheduleRepo) FindByID(ctx context.Context, id string) (*entity.ScheduleDetailResponse, error) {
	s := m.schedules[id]
	if s == nil {
		return nil, nil
	}
	return &entity.ScheduleDetailResponse{
		ID:        s.ID,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Price:     s.Price,
		Status:    s.Status,
	}, nil
}
func (m *mockScheduleRepo) FindAll(ctx context.Context, filter entity.ScheduleFilter) ([]entity.ScheduleDetailResponse, int, error) {
	return nil, len(m.schedules), nil
}
func (m *mockScheduleRepo) Update(ctx context.Context, s *entity.Schedule) error {
	m.schedules[s.ID] = s
	return nil
}
func (m *mockScheduleRepo) Delete(ctx context.Context, id string) error {
	delete(m.schedules, id)
	return nil
}
func (m *mockScheduleRepo) HasTimeConflict(ctx context.Context, studioID string, startTime, endTime time.Time, excludeID string) (bool, error) {
	return m.hasConflict, nil
}

func TestScheduleUsecase_Create(t *testing.T) {
	movieRepo := &mockMovieRepo{
		movies: map[string]*entity.Movie{
			"m1": {ID: "m1", Title: "Dune 2", DurationMinutes: 120, IsActive: true},
		},
	}
	studioRepo := &mockStudioRepo{
		studios: map[string]*entity.Studio{
			"s1": {ID: "s1", Name: "Studio 1", IsActive: true},
		},
	}
	schedRepo := &mockScheduleRepo{
		schedules:   make(map[string]*entity.Schedule),
		hasConflict: false,
	}

	uc := usecase.NewScheduleUsecase(schedRepo, movieRepo, studioRepo)

	res, err := uc.Create(context.Background(), entity.CreateScheduleRequest{
		MovieID:   "m1",
		StudioID:  "s1",
		StartTime: "2025-09-15T13:00:00Z",
		Price:     50000,
	})
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if res.ID != "new-sched-id" {
		t.Errorf("expected schedule ID 'new-sched-id', got %s", res.ID)
	}

	schedRepo.hasConflict = true
	_, err = uc.Create(context.Background(), entity.CreateScheduleRequest{
		MovieID:   "m1",
		StudioID:  "s1",
		StartTime: "2025-09-15T13:30:00Z",
		Price:     50000,
	})
	if err != usecase.ErrScheduleConflict {
		t.Errorf("expected ErrScheduleConflict, got %v", err)
	}

	schedRepo.hasConflict = false
	_, err = uc.Create(context.Background(), entity.CreateScheduleRequest{
		MovieID:   "m-invalid",
		StudioID:  "s1",
		StartTime: "2025-09-15T13:00:00Z",
		Price:     50000,
	})
	if err != usecase.ErrMovieNotFound {
		t.Errorf("expected ErrMovieNotFound, got %v", err)
	}
}
