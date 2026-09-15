package entity

import "time"

type Movie struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	Genre           string    `json:"genre"`
	AgeRating       string    `json:"age_rating"`
	PosterURL       string    `json:"poster_url"`
	ReleaseDate     string    `json:"release_date"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
