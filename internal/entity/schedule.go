package entity

import (
	"time"
)

type ScheduleStatus string

const (
	ScheduleScheduled ScheduleStatus = "SCHEDULED"
	ScheduleActive    ScheduleStatus = "ACTIVE"
	ScheduleCompleted ScheduleStatus = "COMPLETED"
	ScheduleCancelled ScheduleStatus = "CANCELLED"
)

type Schedule struct {
	ID        string         `json:"id"`
	MovieID   string         `json:"movie_id"`
	StudioID  string         `json:"studio_id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Price     float64        `json:"price"`
	Status    ScheduleStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ScheduleDetailResponse struct {
	ID        string         `json:"id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Price     float64        `json:"price"`
	Status    ScheduleStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Movie     MovieSummary   `json:"movie"`
	Studio    StudioSummary  `json:"studio"`
	Cinema    CinemaSummary  `json:"cinema"`
}

type MovieSummary struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	DurationMinutes int    `json:"duration_minutes"`
	Genre           string `json:"genre"`
	AgeRating       string `json:"age_rating"`
	PosterURL       string `json:"poster_url"`
}

type StudioSummary struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	StudioType StudioType `json:"studio_type"`
	TotalSeats int        `json:"total_seats"`
}

type CinemaSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CityName string `json:"city_name"`
	Address  string `json:"address"`
}

type CreateScheduleRequest struct {
	MovieID   string  `json:"movie_id" binding:"required,uuid"`
	StudioID  string  `json:"studio_id" binding:"required,uuid"`
	StartTime string  `json:"start_time" binding:"required"`
	EndTime   string  `json:"end_time"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type UpdateScheduleRequest struct {
	MovieID   *string         `json:"movie_id,omitempty"`
	StudioID  *string         `json:"studio_id,omitempty"`
	StartTime *string         `json:"start_time,omitempty"`
	EndTime   *string         `json:"end_time,omitempty"`
	Price     *float64        `json:"price,omitempty"`
	Status    *ScheduleStatus `json:"status,omitempty"`
}

type ScheduleFilter struct {
	MovieID  string
	StudioID string
	CinemaID string
	CityID   string
	Date     string
	Status   string
	Page     int
	Limit    int
}
