package network

type ServiceInfo struct {
	Port    int
	Address string
	Debug   bool
}

func NewServiceInfo(port *int, address string, debug bool) *ServiceInfo {
	return &ServiceInfo{
		Port:    *port,
		Address: address,
		Debug:   debug,
	}
}
