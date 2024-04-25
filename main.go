package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/luo2pei4/base-server/logger"
	"github.com/luo2pei4/base-server/routers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start base-server",
	Run:   start,
}

var (
	flagServerPort string
	flagLogLevel   string
	flagLogFile    string
)

func init() {
	startCmd.Flags().StringVarP(&flagServerPort, "port", "p", ":8080", "server port, like '8080' or ':8080'")
	startCmd.Flags().StringVarP(&flagLogLevel, "log-level", "l", "info", "the level of log")
	startCmd.Flags().StringVarP(&flagLogFile, "log-file", "f", "./base-server.log", "the name of log file with full path, the file name must with suffix '.log'")
}

func main() {
	startCmd.Execute()
}

func checkArgs() error {
	// 检查端口
	matched, err := regexp.MatchString("^:?[0-9]{4,5}", flagServerPort)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("port is invalid")
	}
	// 检查日志等级
	switch flagLogLevel {
	case "info":
	case "warn":
	case "debug":
	default:
		return errors.New("log level must be info/warn/debug")
	}
	// 检查日志文件名称
	if !strings.HasSuffix(flagLogFile, ".log") {
		return errors.New("invalid log file name")
	}
	return nil
}

func start(cmd *cobra.Command, args []string) {
	// 检查参数有效性
	if err := checkArgs(); err != nil {
		log.Fatalln(err.Error())
	}

	// 解析日志等级
	logLevel, _ := logrus.ParseLevel(flagLogLevel)

	// 初始化日志框架
	logger.InitLog(logLevel, flagLogFile)

	// 初始化router
	router := routers.InitRouter()

	srv := &http.Server{
		Addr:    flagServerPort, // 设置端口
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
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
