package admin

import "time"

type AdRequest struct {
	Title      string      `json:"title"`
	StartAt    time.Time   `json:"start_at"`
	EndAt      time.Time   `json:"end_at"`
	Conditions Contition `json:"conditions"`
}

type Contition struct {
	AgeStart *int      `json:"age_start"`
	AgeEnd   *int      `json:"age_end"`
	Gender   []string `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}
