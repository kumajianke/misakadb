package main

import (
	"flag"
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	"misakadb/network"
	"misakadb/network/core"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	fmt.Println()
	fmt.Println("MisakaDB Service V0.0.1.")
	fmt.Println()

	// 解析命令行参数
	port := flag.Int("port", 10032, "服务端口")
	address := flag.String("address", "0.0.0.0", "服务地址")
	configs := flag.String("configs", "", "配置文件路径[可选]")
	debug := flag.Bool("debug", false, "调试模式")
	flag.Parse()

	// 加载参数信息到ServiceInfo 用于创建套接字
	var serviceInfo *network.ServiceInfo

	if *configs == "" {
		// 从命令行加载
		serviceInfo = network.NewServiceInfo(port, *address, *debug)
		*configs = "misaka.yaml"
		cfg, err := config.InitGlobalMisakaConfigure(*configs)
		if err != nil {
			clilog.Error("缺省配置文件失败, 请确认misakadb的根目录有misaka.yaml文件:", err)
			os.Exit(1)
		}
		cfg.Network.Address = *address
		cfg.Network.Port = *port
	} else {
		// 从配置文件中加载
		cfg, err := config.InitGlobalMisakaConfigure(*configs)
		if err != nil {
			clilog.Error("加载配置文件失败:", err)
			os.Exit(1)
		}
		serviceInfo = config.ConvertServiceInfo(cfg)
	}

	if serviceInfo.Debug {
		go http.ListenAndServe(":6060", nil)
		clilog.Success("pprof running on 0.0.0.0:6060")
	}

	clilog.Info("misakadb running on", serviceInfo.Address+":"+fmt.Sprint(serviceInfo.Port))
	serviceCore := core.NewServiceCore(serviceInfo) // 创建服务核心
	err := serviceCore.Run()                        // 启动服务核心

	if err != nil {
		clilog.Error("服务运行失败:", err)
	}

}
