package main

import (
	"autodg/server"
	"autodg/service"
	"flag"
	"github.com/pkg/errors"
	"github.com/wentaojin/transferdb/pkg/signal"
	"go.uber.org/zap"
	"log"
	_ "net/http/pprof"
	"os"
)

var (
	conf    = flag.String("config", "config.toml", "specify the configuration file, default is config.toml")
	mode    = flag.String("mode", "", "specify the program running mode: [check prepare]")
	version = flag.Bool("version", false, "view AutoDG version info")
)

func main() {
	// 命令行参数解析，初始化
	flag.Parse()

	// 获取程序版本
	service.GetAppVersion(*version)

	// 读取配置文件
	cfg, err := service.ReadConfigFile(*conf)
	if err != nil {
		log.Fatalf("read config file [%s] failed: %v", *conf, err)
	}

	//go func() {
	//	if err = http.ListenAndServe(cfg.AppConfig.PprofPort, nil); err != nil {
	//		service.Logger.Fatal("listen and serve pprof failed", zap.Error(errors.Cause(err)))
	//	}
	//	os.Exit(0)
	//}()

	// 初始化日志 logger
	if err = service.NewZapLogger(cfg); err != nil {
		log.Fatalf("create global zap logger failed: %v", err)
	}
	service.RecordAppVersion("autoDG", service.Logger, cfg)

	// 信号量监听处理
	signal.SetupSignalHandler(func() {
		os.Exit(1)
	})

	// 程序运行
	if err = server.Run(cfg, *mode); err != nil {
		service.Logger.Fatal("server run failed", zap.Error(errors.Cause(err)))
	}
}
