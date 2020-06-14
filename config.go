package fast_r

import "time"

type Config struct {
	MaxTime            time.Duration
	MaximumConnections int
	EnableCaching      bool
}
