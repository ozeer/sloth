package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"syscall"

	"github.com/ozeer/sloth/config"
	"github.com/ozeer/sloth/global"
	"github.com/ozeer/sloth/model/storage"
	"github.com/ozeer/sloth/routers"
	"github.com/ozeer/sloth/service"

	graceHttp "github.com/facebookgo/grace/gracehttp"
	errGroup "golang.org/x/sync/errgroup"
)

var (
	configFile string
)

func main() {
	// 解析命令行参数
	flag.StringVar(&configFile, "c", "", "./app-name -c /path/to/config.ini")
	flag.StringVar(&configFile, "config", "", "./app-name -c /path/to/config.ini")
	flag.Parse()

	// 加载配置文件
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Printf("Load ini config file error: '%v'", err)
		return
	}

	// 初始化日志配置
	// if err := tool.InitLog(conf); err != nil {
	// 	log.Fatalf("Can't load log module, error: %v", err)
	// }

	// 初始化日志配置
	global.Logger = config.InitLogger(conf)
	defer func() {
		err := global.Logger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			// 处理错误的逻辑
			global.Error("日志错误：", err.Error())
		}
	}()

	// 初始化Redis配置
	if err := storage.InitRedis(conf); err != nil {
		log.Fatalf("Can't init redis, error: %v", err)
	}

	// 初始化consumer queue配置
	if err := service.InitConsumer(conf); err != nil {
		log.Fatalf("Can't init consumer redis, error: %v", err)
	}

	// 初始化钉钉机器人配置
	// if err := third.InitDingTalk(conf); err != nil {
	// 	log.Fatalf("Can't init ding talk, error: %v", err)
	// }

	// 初始化队列
	service.InitTimerTask(conf)

	// 初始化路由，启动http服务
	var g errGroup.Group
	g.Go(func() error {
		return graceHttp.Serve(&http.Server{
			Addr:    conf.Gin.Address + ":" + conf.Gin.Port,
			Handler: routers.Init(conf),
		})
	})
	global.InfoF("Start http server! %s:%s", conf.Gin.Address, conf.Gin.Port)

	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}
}
