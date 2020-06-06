package packetio

import (
	"bytes"
	"errors"
	"io"
	"time"
)

type UploadStream struct {
	TotalBytes    int64
	TotalTime     time.Duration
	TimeStarted   time.Time
	TimeEnded     time.Time
	Callback      func(int64, time.Duration)
	payloadBuffer bytes.Buffer
	cur           int
}

func (u *UploadStream) Read(p []byte) (n int, err error) {
	n, err = u.payloadBuffer.Read(p)
	if err != nil {
		if errors.Is(err, io.EOF) {
			//maybe some custom close action ?
		}
		return n, err
	}
	u.TotalBytes += int64(n)
	u.TotalTime = time.Now().Sub(u.TimeStarted)
	go u.Callback(u.TotalBytes, u.TotalTime)
	return n, err
}
