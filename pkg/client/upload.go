package client

import (
	"context"
	"math"
	"time"

	"github.com/herocod3r/fast-r/pkg/network/http"
	"github.com/herocod3r/fast-r/pkg/packetio"

	"github.com/herocod3r/fast-r/pkg/network"
	"github.com/herocod3r/fast-r/pkg/watcher"
)

type Upload struct {
	server       *network.Server
	speedMonitor *watcher.Listener
	cancelFunc   context.CancelFunc
	ctx          context.Context
	Speed        float64
	started      bool
	curSpeedFunc func(curSpeed float64)
}

func NewUpload(server *network.Server, curSpeedFunc func(curSpeed float64)) *Upload {
	return &Upload{server: server, curSpeedFunc: curSpeedFunc}
}

func (u *Upload) Start() error {
	runner := http.NewHandler()
	u.ctx, u.cancelFunc = context.WithCancel(context.Background())
	u.speedMonitor = watcher.NewListiner(u.cancelFunc)
	uploadStream := packetio.NewUploadStream(func(i int64, duration time.Duration) {
		if i <= 500000 {
			return
		}
		speed := ((float64(i * 8)) / duration.Seconds()) / float64(1000) //KiloBit
		monitorSpeed := speed / 1000
		u.speedMonitor.Listen(i, float32(math.Floor(monitorSpeed*100)/100))
		go u.curSpeedFunc(speed)
	})
	u.started = true
	err := runner.ExecuteUpload(u.ctx, u.server, uploadStream)
	u.Speed = ((float64(uploadStream.TotalBytes * 8)) / uploadStream.TotalTime.Seconds()) / float64(1000)
	if err != nil {
		return err
	}

	return nil
}
