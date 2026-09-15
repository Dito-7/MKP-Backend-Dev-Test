package entity

import "time"

type StudioType string

const (
	StudioRegular  StudioType = "REGULAR"
	StudioIMAX     StudioType = "IMAX"
	StudioPremiere StudioType = "PREMIERE"
	Studio4DX      StudioType = "4DX"
)

type Studio struct {
	ID         string     `json:"id"`
	CinemaID   string     `json:"cinema_id"`
	CinemaName string     `json:"cinema_name,omitempty"`
	CityName   string     `json:"city_name,omitempty"`
	Name       string     `json:"name"`
	StudioType StudioType `json:"studio_type"`
	TotalRows  int        `json:"total_rows"`
	TotalCols  int        `json:"total_cols"`
	TotalSeats int        `json:"total_seats"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
