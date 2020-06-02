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
)

type speedTestService struct {
	ignoreIds []string
}

func (s *speedTestService) GetServers(max int, client network.Client) ([]network.Server, error) {
	if len(client.MetaData) < 1 {
		return nil, errors.New("client metadata missing")
	}

	s.ignoreIds = strings.Split(client.MetaData, ",")

}

func (s *speedTestService) GetClientInfo() (client network.Client, er error) {
	ul, _ := url.Parse(configUrl)
	queries := ul.Query()
	queries.Add("X", time.Now().UTC().String())
	ul.RawQuery = queries.Encode()
	configEndpointUrl := ul.String()
	rsp, er := http.Get(configEndpointUrl)
	if er != nil {
		er = &network.Error{InternalError: er}
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
	}
	if serverConfigElm != nil {
		client.MetaData = serverConfigElm.SelectAttrValue("ignoreids", "")
	}
	return
}
