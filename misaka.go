package main

import (
	"flag"
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	"misakadb/misaka_network"
	"misakadb/misaka_network/active"
	"os"
)

func main() {
	fmt.Println("\r\nMisakaDB Service V0.0.1. \r\n")

	// 解析命令行参数
	port := flag.Int("port", 10032, "服务端口")
	address := flag.String("address", "0.0.0.0", "服务地址")
	configer := flag.String("config", "", "配置文件路径")
	flag.Parse()

	// 加载参数信息到ServiceInfo 用于创建套接字
	var serviceInfo *misaka_network.ServiceInfo
	if *configer == "" {

		// 从命令行加载
		serviceInfo = misaka_network.NewServiceInfo(port, *address)
		clilog.Info("misakadb running on", *address+":"+fmt.Sprint(*port))
	} else {

		// 从配置文件中加载
		cfg, err := config.InitGlobalMisakaConfigure(*configer)
		if err != nil {
			clilog.Error("加载配置文件失败:", err)
			os.Exit(1)
		}
		serviceInfo = config.ConvertServiceInfo(cfg)
		clilog.Info("misakadb running on", serviceInfo.Address+":"+fmt.Sprint(serviceInfo.Port))
	}

	serviceCore := active.NewServiceCore(serviceInfo)
	err := serviceCore.Run()

	if err != nil {
		clilog.Error("服务运行失败:", err)
	}
}
