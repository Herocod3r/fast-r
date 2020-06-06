package packetio

import (
	"io"
	"math"
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxUploadSize = 10485760 //10mb
)

type UploadStream struct {
	TotalBytes    int64
	TotalTime     time.Duration
	TimeStarted   time.Time
	TimeEnded     time.Time
	Callback      func(int64, time.Duration)
	payloadBuffer []byte
}

func (u *UploadStream) Read(p []byte) (n int, err error) {
	if u.payloadBuffer == nil {
		u.payloadBuffer = u.getDefaultBuffer(512)
	}

	lenP := len(p)

	if lenP <= cap(u.payloadBuffer) {
		n = copy(p, u.payloadBuffer[:len(p)])
	} else {
		//extra := lenP - cap(u.payloadBuffer)
		runs := math.Ceil(float64(lenP) / float64(cap(u.payloadBuffer)))
		for i := 0; i < int(runs); i++ {
			if i == 0 {
				tempn := copy(p[i:], u.payloadBuffer)
				n += tempn
				continue
			}
			extra := lenP - n
			var tempn int
			if extra > cap(u.payloadBuffer) {
				tempn = copy(p[n:], u.payloadBuffer)
			} else {
				tempn = copy(p[n:], u.payloadBuffer[:extra])
			}
			n += tempn
		}
	}

	u.TotalBytes += int64(n)
	u.TotalTime = time.Now().Sub(u.TimeStarted)
	go u.Callback(u.TotalBytes, u.TotalTime)
	if n >= maxUploadSize {
		err = io.EOF
		u.TimeEnded = time.Now()
	}
	return n, err
}

func (UploadStream) getDefaultBuffer(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
