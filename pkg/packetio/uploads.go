package packetio

import (
	"io"
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxUploadSize = 5242880 //5mb
)

type UploadStream struct {
	TotalBytes    int64
	TotalTime     time.Duration
	TimeStarted   time.Time
	TimeEnded     time.Time
	Callback      func(int64, time.Duration)
	payloadBuffer []byte
}

func NewUploadStream(callback func(int64, time.Duration)) *UploadStream {
	return &UploadStream{Callback: callback}
}

func (u *UploadStream) Read(p []byte) (n int, err error) {
	if u.payloadBuffer == nil {
		u.payloadBuffer = u.getDefaultBuffer(1024)
		u.TimeStarted = time.Now()
	}
	if u.TotalBytes >= maxUploadSize {
		return 0, io.EOF
	}

	n = copy(p, u.payloadBuffer)
	//
	//lenP := len(p)
	//
	//if lenP <= cap(u.payloadBuffer) {
	//	n = copy(p, u.payloadBuffer[:len(p)])
	//} else {
	//	//extra := lenP - cap(u.payloadBuffer)
	//	runs := math.Ceil(float64(lenP) / float64(cap(u.payloadBuffer)))
	//	for i := 0; i < int(runs); i++ {
	//		if i == 0 {
	//			tempn := copy(p[i:], u.payloadBuffer)
	//			n += tempn
	//			continue
	//		}
	//		extra := lenP - n
	//		var tempn int
	//		if extra > cap(u.payloadBuffer) {
	//			tempn = copy(p[n:], u.payloadBuffer)
	//		} else {
	//			tempn = copy(p[n:], u.payloadBuffer[:extra])
	//		}
	//		n += tempn
	//	}
	//}

	u.TotalBytes += int64(n)
	u.TotalTime = time.Now().Sub(u.TimeStarted)
	u.Callback(u.TotalBytes, u.TotalTime)
	return n, err
}

func (UploadStream) getDefaultBuffer(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
