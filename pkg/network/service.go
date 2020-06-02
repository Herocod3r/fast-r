package network

type Service interface {
	GetServers(max int, client Client) ([]Server, error)
	//todo split the interface to isolate client information
	GetClientInfo() (Client, error)
}
