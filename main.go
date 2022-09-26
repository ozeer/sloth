package main

import (
	"flag"
	"github.com/facebookgo/grace/gracehttp"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"sloth/config"
	"sloth/model/storage"
	"sloth/routers"
	"sloth/service"
	"sloth/third"
	"sloth/tool"
)

var (
	configFile string
)

func main() {
	// 解析命令行参数
	ParseCommandArgs(configFile)

	// 加载配置文件
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Printf("Load ini config file error: '%v'", err)
		return
	}

	// 初始化日志配置
	if err := tool.InitLog(conf); err != nil {
		log.Fatalf("Can't load log module, error: %v", err)
	}

	// 初始化Redis配置
	if err := storage.InitRedis(conf); err != nil {
		log.Fatalf("Can't init redis, error: %v", err)
	}

	// 初始化consumer queue配置
	if err := service.InitConsumer(conf); err != nil {
		log.Fatalf("Can't init consumer redis, error: %v", err)
	}

	// 初始化钉钉机器人配置
	if err := third.InitDingTalk(conf); err != nil {
		log.Fatalf("Can't init ding talk, error: %v", err)
	}

	// 初始化队列
	service.Init(conf)

	// 初始化路由，启动http服务
	var g errgroup.Group
	g.Go(func() error {
		return gracehttp.Serve(&http.Server{
			Addr:    conf.Gin.Address + ":" + conf.Gin.Port,
			Handler: routers.Init(conf),
		})
	})
	tool.LogAccess.Infof("Start http server! %s:%s", conf.Gin.Address, conf.Gin.Port)

	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// ParseCommandArgs 解析命令行参数
func ParseCommandArgs(configFile string) {
	flag.StringVar(&configFile, "c", "", "./app-name -c /path/to/config.ini")
	flag.StringVar(&configFile, "config", "", "./app-name -c /path/to/config.ini")
	flag.Parse()
}
