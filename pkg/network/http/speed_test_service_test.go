package http

import "testing"

func TestSpeedTestService_GetServers(t *testing.T) {
	inst := speedTestService{}

	client, er := inst.GetClientInfo()
	if er != nil {
		t.Fail()
		return
	}

	if len(client.Ip) < 1 || len(client.Isp) < 1 || client.Latitude == 0 || client.Longitude == 0 {
		t.Fail()
	}

	servers, err := inst.GetServers(5, client)
	if err != nil {
		t.Fail()
	}
	if len(servers) != 5 {
		t.Fail()
	}
}

func TestSpeedTestService_GetClientInfo(t *testing.T) {
	inst := speedTestService{}

	client, er := inst.GetClientInfo()
	if er != nil {
		t.Fail()
		return
	}

	if len(client.Ip) < 1 || len(client.Isp) < 1 || client.Latitude == 0 || client.Longitude == 0 {
		t.Fail()
	}
}
