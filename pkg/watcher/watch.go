package watcher

import (
	"context"
	"math"
	"time"
)

//This package is responsible for monitoring and pulling the plugs when the speed has been measured

var ByteRate = int64(1048576)

func NewListiner(cancelFunc context.CancelFunc) *Listener {
	return &Listener{cancelFunc: cancelFunc, packetsQueue: NewItemQueue(), currentPacketTime: time.Now()}
}

type Listener struct {
	lastSentByte      int64
	bufferedBytes     int64
	currentPacketTime time.Time
	packetsQueue      *PacketDurationQueue
	cancelFunc        context.CancelFunc
}

func (l *Listener) Listen(bytesTransfered int64) {
	curBytes := bytesTransfered - l.lastSentByte
	l.bufferedBytes += curBytes

	if l.bufferedBytes < ByteRate {
		return
	}

	l.packetsQueue.Enqueue(time.Now().Sub(l.currentPacketTime))
	l.currentPacketTime = time.Now()
	if l.packetsQueue.Size() >= 3 {
		_ = l.packetsQueue.Dequeue()
		pointA := time.Duration(math.Abs(float64(l.packetsQueue.items[0] - l.packetsQueue.items[1])))
		pointB := time.Duration(math.Abs(float64(l.packetsQueue.items[1] - l.packetsQueue.items[2])))

		isUniform := pointA.Milliseconds() < 10 && pointB.Milliseconds() < 10
		if isUniform {
			l.cancelFunc()
		}
	}
	l.bufferedBytes = 0
}
