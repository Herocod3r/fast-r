package packetio

import (
	"bytes"
	"errors"
	"io"
	"time"
)

type DownloadStream struct {
	TotalBytes  int64
	TotalTime   time.Duration
	TimeStarted time.Time
	TimeEnded   time.Time
	Callback    func(int64, time.Duration)
}

func NewDownloadStream(callback func(int64, time.Duration)) *DownloadStream {
	return &DownloadStream{Callback: callback}
}

func (ds *DownloadStream) Process(stream io.ReadCloser) {
	ds.TimeStarted = time.Now()
	buffer := make([]byte, bytes.MinRead)
	for {
		n, err := stream.Read(buffer)
		ds.TotalBytes += int64(n)
		ds.TotalTime = time.Now().Sub(ds.TimeStarted)
		go ds.Callback(ds.TotalBytes, ds.TotalTime)

		if err != nil {
			if errors.Is(err, io.EOF) {
				//maybe some custom close action ?
			}
			break
		}

	}
	ds.TimeEnded = time.Now()
}
