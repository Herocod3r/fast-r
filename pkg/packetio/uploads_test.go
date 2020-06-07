package packetio

import (
	"errors"
	"io"
	"testing"
	"time"
)

func TestUploadStream_Read(t *testing.T) {
	size := (512 * 3) + 4
	buff := make([]byte, size)
	us := UploadStream{Callback: func(i int64, duration time.Duration) {

	}}
	_, er := us.Read(buff)

	if er != nil {
		t.Fail()
	}
}

func TestUploadStream_Read_EOF(t *testing.T) {

	buff := make([]byte, maxUploadSize)
	us := UploadStream{Callback: func(i int64, duration time.Duration) {

	}}

	for {
		n, er := us.Read(buff)
		if n == 0 {
			break
		}
		if er != nil {
			if !errors.Is(er, io.EOF) {
				t.Fail()
			}
		}

	}

}
