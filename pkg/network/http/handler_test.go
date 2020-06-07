package http

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/herocod3r/fast-r/pkg/packetio"
)

func TestHandler_ExecuteDownload(t *testing.T) {
	inst := speedTestService{}

	client, er := inst.GetClientInfo()
	if er != nil {
		t.Fail()
		return
	}

	if len(client.Ip) < 1 || len(client.Isp) < 1 || client.Latitude == 0 || client.Longitude == 0 {
		t.Fail()
	}

	servers, err := inst.GetServers(1, client)
	if err != nil {
		t.Fail()
	}

	handler := Handler{}
	stream, er := handler.ExecuteDownload(context.Background(), &servers[0])
	if er != nil {
		t.Fail()
		return
	}
	defer stream.Close()
	download := packetio.NewDownloadStream(func(i int64, duration time.Duration) {
		speed := (float64(i) / float64(125000)) / duration.Seconds()

		fmt.Println(fmt.Sprintf("Speed is %.1f Mbs", speed))
	})
	download.Process(stream)
	speed := (float64(download.TotalBytes) / float64(125000)) / download.TotalTime.Seconds()

	fmt.Println("Download test complete", download.TotalTime.Seconds(), fmt.Sprintf("Speed is %.1f Mbs", speed))

}

func TestHandler_ExecuteUpload(t *testing.T) {
	inst := speedTestService{}

	client, er := inst.GetClientInfo()
	if er != nil {
		t.Fail()
		return
	}

	if len(client.Ip) < 1 || len(client.Isp) < 1 || client.Latitude == 0 || client.Longitude == 0 {
		t.Fail()
	}

	servers, err := inst.GetServers(1, client)
	if err != nil {
		t.Fail()
	}

	handler := Handler{}
	uploadStream := packetio.NewUploadStream(func(i int64, duration time.Duration) {
		if i <= 500000 {
			return
		}
		speed := (float64(i) / float64(125000)) / duration.Seconds()

		fmt.Println(fmt.Sprintf("Speed is %.1f Mbs", speed))
	})
	er = handler.ExecuteUpload(context.Background(), &servers[0], uploadStream)
	if er != nil {
		t.Fail()
		return
	}

	speed := (float64(uploadStream.TotalBytes) / float64(125000)) / uploadStream.TotalTime.Seconds()

	fmt.Println("Download test complete", uploadStream.TotalTime.Seconds(), fmt.Sprintf("Speed is %.1f Mbs", speed))

}
