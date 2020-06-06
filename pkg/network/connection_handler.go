package network

import "io"

type ConnectionHandler interface {
	ExecuteUpload(server *Server, uploadReader io.Reader) error
	ExecuteDownload(server *Server) (io.ReadCloser, error)
}
