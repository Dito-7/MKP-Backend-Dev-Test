package repository

import (
	"context"

	"mkp-cinema-ticketing/internal/entity"
)

type StudioRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Studio, error)
	FindAll(ctx context.Context, cinemaID string) ([]entity.Studio, error)
}
