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

func (ds *DownloadStream) Process(stream io.Reader) {
	ds.TimeStarted = time.Now()
	buffer := new(bytes.Buffer)
	for {
		//todo find a way to do reads without in memory buffers, maybe a file instead?
		n, err := buffer.ReadFrom(stream)
		if err != nil {
			if errors.Is(err, io.EOF) {
				//maybe some custom close action ?
			}
			break
		}

		ds.TotalBytes += n
		ds.TotalTime = time.Now().Sub(ds.TimeStarted)
		go ds.Callback(ds.TotalBytes, ds.TotalTime)
	}
}
