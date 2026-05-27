package main

import (
	"flag"
	"fmt"
	"misakadb/clilog"
	"misakadb/misaka_network"
	"misakadb/misaka_network/active"
)

func main() {

	fmt.Println("\r\nMisakaDB Service V0.0.1. \r\n")

	port := flag.Int("port", 10032, "服务端口")
	address := flag.String("address", "0.0.0.0", "服务地址")

	flag.Parse()

	var serviceInfo *misaka_network.ServiceInfo = misaka_network.NewServiceInfo(port, *address)
	_ = serviceInfo
	clilog.Info("misakadb running on", *address+":"+fmt.Sprint(*port))

	serviceCore := active.NewServiceCore(serviceInfo)
	err := serviceCore.Run()

	if err != nil {
		clilog.Error("服务运行失败:", err)
	}
}
