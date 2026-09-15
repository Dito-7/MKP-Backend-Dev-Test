package repository

import (
	"context"

	"mkp-cinema-ticketing/internal/entity"
)

type MovieRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Movie, error)
	FindAll(ctx context.Context) ([]entity.Movie, error)
}
