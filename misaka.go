package main

import (
	"flag"
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	"misakadb/network"
	"misakadb/network/RegisterCenter"
	"misakadb/network/core"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var (
	port    = flag.Int("port", 10032, "服务端口")
	address = flag.String("address", "0.0.0.0", "服务地址")
	configs = flag.String("configs", "", "配置文件路径[可选]")
	debug   = flag.Bool("debug", false, "调试模式")
)

func printTitle(cfg *config.MisakaConfigure) {
	title := fmt.Sprintf(" MisakaDB Service V%s ", cfg.Service.Version)
	border := strings.Repeat("-", len(title)+2)
	fmt.Printf("%s\n|%s|\n%s\n\n", border, title, border)
}

func main() {
	flag.Parse()

	// 加载参数信息到ServiceInfo 用于创建套接字
	var serviceInfo network.ServiceInfo
	var cfg *config.MisakaConfigure
	var err_load_cfg error
	if *configs == "" {
		// 从命令行加载
		serviceInfo = network.NewServiceInfo(*port, *address, *debug)
		*configs = "$misaka.yaml"
		if strings.HasPrefix(*configs, "$") {
			*configs = "./profiles/" + (*configs)[1:]
		}

		cfg, err_load_cfg = config.InitGlobalMisakaConfigure(*configs)

		if err_load_cfg != nil {
			clilog.Error("缺省配置文件失败, 请确认misakadb的profiles目录有misaka.yaml文件:", err_load_cfg)
			return
		}
		printTitle(cfg)

		cfg.Network.Address = *address
		cfg.Network.Port = *port
	} else {
		// 从配置文件中加载
		cfg, err_load_cfg = config.InitGlobalMisakaConfigure(*configs)
		if err_load_cfg != nil {
			clilog.Error("加载配置文件失败:", err_load_cfg)
			return
		}
		serviceInfo = config.ConvertServiceInfo(cfg)
	}

	if serviceInfo.Debug {
		go http.ListenAndServe(":6060", nil)
		clilog.Success("pprof running on 0.0.0.0:6060")
	}

	_ = RegisterCenter.NewRegisterCenter(cfg.Network.MaxConn)

	clilog.Info("\nmisakadb running on", serviceInfo.Address+":"+fmt.Sprint(serviceInfo.Port))
	serviceCore := core.NewServiceCore(serviceInfo) // 创建服务核心
	err := serviceCore.Run()                        // 启动服务核心

	if err != nil {
		clilog.Error("服务运行失败:", err)
	}

}
