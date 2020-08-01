package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"

	"github.com/herocod3r/fast-r/pkg/network"
)

const (
	configUrl = "https://www.speedtest.net/speedtest-config.php"
	serverUrl = "https://www.speedtest.net/speedtest-servers-static.php"
)

type speedTestService struct {
	ignoreIds []string
}

func NewSpeedTestService() network.Service {
	return &speedTestService{}
}

func (s *speedTestService) GetServers(max int, client network.Client) (servers []network.Server, er error) {
	if len(client.MetaData) < 1 {
		return nil, errors.New("client metadata missing")
	}

	s.ignoreIds = strings.Split(client.MetaData, ",")

	ul, _ := url.Parse(serverUrl)
	queries := ul.Query()
	queries.Add("X", time.Now().UTC().String())
	ul.RawQuery = queries.Encode()
	configEndpointUrl := ul.String()
	rsp, er := http.Get(configEndpointUrl)
	if er != nil {
		er = network.NetworkAccessErr
		return
	}
	defer rsp.Body.Close()
	buf := new(bytes.Buffer)
	_, er = buf.ReadFrom(rsp.Body)
	if er != nil {
		return
	}

	return s.parseServers(max, buf)

}

func (s *speedTestService) parseServers(max int, buffer *bytes.Buffer) (servers []network.Server, er error) {
	servers = make([]network.Server, 0)
	ignoreIdsMap := make(map[string]int)
	for index, value := range s.ignoreIds {
		ignoreIdsMap[value] = index
	}

	doc := etree.NewDocument()

	count, er := doc.ReadFrom(buffer)
	print(count)
	if er != nil {
		return
	}

	serversElm := doc.FindElement("//settings/servers")

	if serversElm == nil {
		er = errors.New("Unable to parse xml")
		return
	}

	for _, serverElm := range serversElm.ChildElements() {
		if len(servers) >= max {
			break
		}
		id := serverElm.SelectAttrValue("id", "")
		if _, ok := ignoreIdsMap[id]; ok {
			continue
		}
		ul, _ := url.Parse(serverElm.SelectAttrValue("url", ""))
		ul.Path = ""
		servers = append(servers, network.Server{
			Name:        serverElm.SelectAttrValue("sponsor", ""),
			Address:     ul.String(),
			PingAddress: ul.String() + "/latency.txt",
		})
	}

	return
}

func (s *speedTestService) GetClientInfo() (client network.Client, er error) {
	ul, _ := url.Parse(configUrl)
	queries := ul.Query()
	queries.Add("X", time.Now().UTC().String())
	ul.RawQuery = queries.Encode()
	configEndpointUrl := ul.String()
	rsp, er := http.Get(configEndpointUrl)
	if er != nil {
		er = network.NetworkAccessErr
		return
	}
	defer rsp.Body.Close()
	buf := new(bytes.Buffer)
	_, er = buf.ReadFrom(rsp.Body)
	if er != nil {
		return
	}

	return s.parseClientXml(buf)
}

func (s *speedTestService) parseClientXml(buffer *bytes.Buffer) (client network.Client, er error) {
	doc := etree.NewDocument()

	count, er := doc.ReadFrom(buffer)
	print(count)
	if er != nil {
		return
	}

	settingsElm := doc.FindElement("//settings/client")
	serverConfigElm := doc.FindElement("//settings/server-config")

	if settingsElm == nil {
		er = errors.New("Unable to parse xml")
		return
	}

	latitude, _ := strconv.ParseFloat(settingsElm.SelectAttrValue("lat", "0"), 32)
	longitude, _ := strconv.ParseFloat(settingsElm.SelectAttrValue("lon", "0"), 32)

	client = network.Client{
		Ip:        settingsElm.SelectAttrValue("ip", "::"),
		Latitude:  float32(latitude),
		Longitude: float32(longitude),
		Isp:       settingsElm.SelectAttrValue("isp", ""),
		Location:  settingsElm.SelectAttrValue("country", ""),
	}
	if serverConfigElm != nil {
		client.MetaData = serverConfigElm.SelectAttrValue("ignoreids", "")
	}
	return
}
