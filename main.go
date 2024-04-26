package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luo2pei4/base-server/configs"
	"github.com/luo2pei4/base-server/internal/logger"
	"github.com/luo2pei4/base-server/routers"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start base-server",
	Run:   start,
}

var flagServiceConfigFile string

func init() {
	startCmd.Flags().StringVarP(&flagServiceConfigFile, "config", "c", "./service-config.toml", "service config file path")
}

func main() {
	// 加载配置文件
	configs.LoadServiceConfig(flagServiceConfigFile)
	// 启动配置文件监控
	configs.StartServiceConfigWatch()
	// 执行命令行
	startCmd.Execute()
}

func start(cmd *cobra.Command, args []string) {

	// 初始化日志框架
	logger.InitLog(configs.GetLogLevel(), configs.GetLogFile())

	// 设置最大cpu使用数，默认50%
	configs.SetMaxCPUNum()

	// 初始化router
	router := routers.InitRouter()

	// 获取服务端口号
	port, err := configs.GetSerivePort()
	if err != nil {
		log.Fatal(err.Error())
	}

	srv := &http.Server{
		Addr:    port, // 设置端口
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panicf("server error, %s", err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Infof("Server Shutdown: %s", err.Error())
	}
	logger.Info("Server exiting")
}
