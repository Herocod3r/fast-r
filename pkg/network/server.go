package network

import (
	"errors"
	"net/http"
	"time"
)

type Server struct {
	Name        string
	Address     string
	PingAddress string
	Latency     float32
}

func (s Server) PingForLatency() (time.Duration, error) {
	currentTime := time.Now()

	rsp, err := http.Get(s.PingAddress)
	if err != nil {
		return 0, err
	}

	if rsp.StatusCode != http.StatusOK {
		return 0, errors.New("Unable To Contact Server")
	}

	return time.Now().Sub(currentTime), nil
}
