package client

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/herocod3r/fast-r/pkg/cache"
	"github.com/herocod3r/fast-r/pkg/network"
	"github.com/herocod3r/fast-r/pkg/network/http"
)

type serverCacheEntry struct {
	Servers    []network.Server
	TimeCached time.Time
}

func GetServersList(considerCache bool, cacheClient cache.Client) ([]network.Server, error) {
	if considerCache && cacheClient != nil {
		servers, err := getServerFromCache(cacheClient)
		if len(servers) > 0 && err == nil {
			return sortByLatency(servers)
		}
	}
	service := http.NewSpeedTestService()
	client, er := service.GetClientInfo()
	if er != nil {
		return nil, er
	}
	servers, er := service.GetServers(10, client)
	if er != nil {
		return nil, er
	}

	if cacheClient != nil {
		data, _ := json.Marshal(&serverCacheEntry{Servers: servers, TimeCached: time.Now()})
		_ = cacheClient.Set("serversData", string(data))
	}

	return sortByLatency(servers)
}

func getServerFromCache(cacheCLient cache.Client) ([]network.Server, error) {
	clientData, er := cacheCLient.Get("serversData")
	if er != nil {
		if errors.Is(er, cache.StoreNotActiveErr) {
			//log caching not supported
		}
		return nil, er
	}
	serversCacheEntry := serverCacheEntry{}
	er = json.Unmarshal([]byte(clientData), &serversCacheEntry)
	if er != nil {
		return nil, er
	}

	if len(serversCacheEntry.Servers) < 1 {
		return nil, errors.New("Cache Is Invalid")
	}

	periodOfCache := time.Now().Sub(serversCacheEntry.TimeCached)
	if periodOfCache > (time.Minute * 10) { //10Minutes Limit
		return nil, errors.New("Cache Is Invalid")
	}

	return serversCacheEntry.Servers, nil
}

func sortByLatency(servers []network.Server) ([]network.Server, error) {
	wg := new(sync.WaitGroup)
	latencyServers := make([]network.Server, 0)
	for _, server := range servers {
		wg.Add(1)
		go func(serv network.Server) {
			latency, err := serv.PingForLatency()
			serv.Latency = float32(latency)
			if err == nil {
				latencyServers = append(latencyServers, serv)
			}
			wg.Done()
		}(server)
	}
	wg.Wait()

	sort.Slice(latencyServers, func(i, j int) bool {
		return latencyServers[i].Latency <= latencyServers[j].Latency
	})
	return latencyServers, nil
}
