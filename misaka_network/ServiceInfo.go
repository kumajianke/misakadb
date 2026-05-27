package misaka_network

type ServiceInfo struct {
	Port    int
	Address string
}

func NewServiceInfo(port *int, address string) *ServiceInfo {
	return &ServiceInfo{
		Port:    *port,
		Address: address,
	}
}
