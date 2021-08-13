package main

import (
	"context"
	"flag"
	"fmt"
	"g2ww/grafana/ngalert"
	"g2ww/grafana/old"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// http服务
var srv *http.Server

func main() {
	// 定义变量 用于接收命令行参数
	var port int
	var version string
	flag.IntVar(&port, "port", 3001, `Server port, default: 3001`)
	flag.StringVar(&version, "version", "old", `Grafana alert version, default: "old", optional: "ngalert"`)
	flag.Parse()
	fmt.Println("G2WW server running on port", port)

	app := gin.Default()
	// Server Info
	if version == "ngalert" {
		app.GET("/", ngalert.GetSendCount)
		app.POST("/send", ngalert.SendMsg)
	} else if version == "old" {
		app.GET("/", old.GetSendCount)
		app.POST("/send", old.SendMsg)
	} else {
		fmt.Printf(`[Error] error param "version": %v`, version)
		fmt.Println()
		shutdown()
	}
	srv = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: app,
	}

	// 启动http请求
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	shutdown()
}

// 优雅的关闭
func shutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println()
		fmt.Println(fmt.Errorf("server shutdown:[%v]", err))
		return
	}
}
