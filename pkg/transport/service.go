package transport

type Service interface {
	GetServers(max int) []Server
}
