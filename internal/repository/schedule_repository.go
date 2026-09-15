package repository

import (
	"context"
	"time"

	"mkp-cinema-ticketing/internal/entity"
)

type ScheduleRepository interface {
	Create(ctx context.Context, s *entity.Schedule) error
	FindByID(ctx context.Context, id string) (*entity.ScheduleDetailResponse, error)
	FindAll(ctx context.Context, filter entity.ScheduleFilter) ([]entity.ScheduleDetailResponse, int, error)
	Update(ctx context.Context, s *entity.Schedule) error
	Delete(ctx context.Context, id string) error
	HasTimeConflict(ctx context.Context, studioID string, startTime, endTime time.Time, excludeID string) (bool, error)
}
