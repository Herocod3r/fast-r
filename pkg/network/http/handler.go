package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/herocod3r/fast-r/pkg/network"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) ExecuteUpload(ctx context.Context, server *network.Server, uploadReader io.Reader) error {
	req, _ := http.NewRequestWithContext(ctx, "POST", server.Address, uploadReader)
	req.ContentLength = 5242880 //5MB (WorstCase)
	req.Header.Set("Content-Type", "text/plain")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK {
		return network.NetworkAccessErr
	}

	return nil
}

func (h Handler) ExecuteDownload(ctx context.Context, server *network.Server) (io.ReadCloser, error) {
	ul, _ := url.Parse(server.Address)
	ul.Path = ""
	downloadUrl := fmt.Sprintf("%s/speedtest/random6000x6000.jpg?r=2", ul.String())
	req, _ := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, network.NetworkAccessErr
	}

	return rsp.Body, nil
}
