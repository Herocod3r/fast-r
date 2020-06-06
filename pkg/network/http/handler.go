package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/herocod3r/fast-r/pkg/network"
)

type Handler struct {
}

func (h Handler) ExecuteUpload(server *network.Server, uploadReader io.Reader) error {
	panic("implement me")
}

func (h Handler) ExecuteDownload(server *network.Server) (io.ReadCloser, error) {
	ul, _ := url.Parse(server.Address)
	ul.Path = ""
	downloadUrl := fmt.Sprintf("%s/speedtest/random3000x3000.jpg?r=2", ul.String())

	rsp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, &network.Error{InternalError: errors.New("Unable to connect to network")}
	}

	return rsp.Body, nil
}
