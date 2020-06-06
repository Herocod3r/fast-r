package packetio

import (
	"testing"
	"time"
)

func TestUploadStream_Read(t *testing.T) {
	size := (512 * 3) + 4
	buff := make([]byte, size)
	us := UploadStream{Callback: func(i int64, duration time.Duration) {

	}}
	n, er := us.Read(buff)
	if n != size {
		t.Fail()
	}

	if er != nil {
		t.Fail()
	}
}

func TestUploadStream_Read_EOF(t *testing.T) {

	buff := make([]byte, maxUploadSize)
	us := UploadStream{Callback: func(i int64, duration time.Duration) {

	}}
	n, er := us.Read(buff)
	if n != maxUploadSize {
		t.Fail()
	}

	if er == nil {
		t.Fail()
	}
}
