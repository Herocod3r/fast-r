package network

import (
	"context"
	"io"
)

type ConnectionHandler interface {
	ExecuteUpload(ctx context.Context, server *Server, uploadReader io.Reader) error
	ExecuteDownload(ctx context.Context, server *Server) (io.ReadCloser, error)
}
