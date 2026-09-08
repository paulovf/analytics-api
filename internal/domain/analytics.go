package domain

import (
	"time"
)

type DailyStatusStat struct {
	DayBucket   time.Time `json:"day"`
	Status      string    `json:"status"`
	TotalEvents int64     `json:"total_events"`
}
