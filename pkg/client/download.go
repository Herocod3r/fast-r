package client

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	"github.com/herocod3r/fast-r/pkg/packetio"

	"github.com/herocod3r/fast-r/pkg/network"
	"github.com/herocod3r/fast-r/pkg/network/http"
	"github.com/herocod3r/fast-r/pkg/watcher"
)

type Download struct {
	server       *network.Server
	speedMonitor *watcher.Listener
	cancelFunc   context.CancelFunc
	ctx          context.Context
	Speed        float64
	started      bool
	curSpeedFunc func(curSpeed float64)
}

func NewDownload(server *network.Server, curSpeedFunc func(curSpeed float64)) *Download {
	return &Download{server: server, curSpeedFunc: curSpeedFunc}
}

func (d *Download) Start() error {
	runner := http.NewHandler()
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	d.speedMonitor = watcher.NewListiner(d.cancelFunc)
	stream, err := runner.ExecuteDownload(d.ctx, d.server)
	if err != nil {
		return err
	}
	defer stream.Close()
	d.started = true
	d.processDownloadTest(stream)
	return nil
}

func (d *Download) Stop() error {
	if !d.started {
		return errors.New("Process Has Not Started")
	}
	d.cancelFunc()
	d.started = false
	return nil
}

func (d *Download) processDownloadTest(stream io.ReadCloser) {
	downloadProcessor := packetio.NewDownloadStream(func(i int64, duration time.Duration) {
		speed := ((float64(i * 8)) / duration.Seconds()) / float64(1000) //KiloBit
		monitorSpeed := speed / 1000
		d.speedMonitor.Listen(i, float32(math.Floor(monitorSpeed*100)/100))
		go d.curSpeedFunc(speed)
	})
	downloadProcessor.Process(stream)
	d.Speed = ((float64(downloadProcessor.TotalBytes * 8)) / downloadProcessor.TotalTime.Seconds()) / float64(1000)
}
