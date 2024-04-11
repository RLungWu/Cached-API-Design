package api

import "time"

type AdRequest struct {
	Title   string `json:"title"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Condition []Condition `json:"condition"`
}

type Condition struct {
	AgeStart int `json:"age_start"`
	AgeEnd   int `json:"age_end"`
	Country []string `json:"country"`
	Platform []string `json:"platform"`
}
