package fast_r

import "time"

type Summary struct {
	TotalDownloaded float32
	TotalUploaded   float32
	StartedAt       time.Time
	FinishedAt      time.Time
}
