package network

import "io"

type ConnectionHandler interface {
	GetUploadConnection(server *Server) (error, io.Writer)
	GetDownloadConnection(server *Server) (error, io.Reader)
}
