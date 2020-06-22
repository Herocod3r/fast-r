package watcher

import (
	"context"
)

//This package is responsible for monitoring latency

var ByteRate int64 = 102400 //100kb
var MinimumSpeedBlock int = 10

func NewListiner(cancelFunc context.CancelFunc) *Listener {
	return &Listener{cancelFunc: cancelFunc, packetsQueue: NewItemQueue()}
}

type Listener struct {
	lastSentByte  int64
	bufferedBytes int64
	packetsQueue  *PacketDurationQueue
	cancelFunc    context.CancelFunc
}

func (l *Listener) Listen(bytesTransfered int64, speed float32) {
	curBytes := bytesTransfered - l.lastSentByte
	l.bufferedBytes += curBytes

	if l.bufferedBytes < ByteRate {
		return
	}

	l.packetsQueue.Enqueue(speed)
	if l.packetsQueue.Size() >= MinimumSpeedBlock {
		isEven := true
		for i := 1; i < len(l.packetsQueue.items); i++ {
			if l.packetsQueue.items[i-1] != l.packetsQueue.items[i] {
				isEven = false
				break
			}
		}
		if isEven {
			l.cancelFunc()
			return
		}
		_ = l.packetsQueue.Dequeue()
	}
	l.bufferedBytes = 0
}
