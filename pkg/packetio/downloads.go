package packetio

import (
	"errors"
	"io"
	"time"
)

const MinBuffer = 10485760 //10Mb buffer

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
	buffer := make([]byte, MinBuffer)
	for {
		n, err := stream.Read(buffer)
		ds.TotalBytes += int64(n)
		ds.TotalTime = time.Now().Sub(ds.TimeStarted)
		ds.Callback(ds.TotalBytes, ds.TotalTime)

		if err != nil {
			if errors.Is(err, io.EOF) {
				//maybe some custom close action ?
			}
			break
		}

	}
	ds.TimeEnded = time.Now()
}
