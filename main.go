package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sf-paas/config"
	"sf-paas/middleware"
	"sf-paas/router"
	"syscall"
	"time"
)

var log = middleware.Logger

func main() {
	r := router.InitRouter()

	log.Info("init router success")
	//启动服务
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 创建系统信号接收器
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// 创建5s的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("stop server failed:%v\n", err)
	}
	log.Info("stop server success")
}
