package http

import (
	"io"
	"net/http"

	"github.com/herocod3r/fast-r/pkg/network"
)

type Handler struct {
}

func (h Handler) GetUploadConnection(server *network.Server) (error, io.Writer) {
	panic("implement me")
}

func (h Handler) GetDownloadConnection(server *network.Server) (error, io.Reader) {
	http.Post()
	panic("implement me")
}
