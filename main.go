package main

import (
	"context"
	"flag"
	"fmt"
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
	var alertVersion string
	flag.IntVar(&port, "port", 3001, `Server port, default: 3001`)
	flag.StringVar(&alertVersion, "alertversion", "", `Grafana alert version, default: "", optional: "ngalert"`)
	flag.Parse()
	fmt.Printf("G2WW server running on port %v", port)
	fmt.Println()
	// 此处多此一举的原因是为了适应 grafana 的变量 可视情况改掉
	if alertVersion == "" {
		alertVersion = "old"
	}
	fmt.Printf("G2WW server is based on the %v alert version", alertVersion)
	fmt.Println()

	app := gin.Default()
	// Server Info
	app.GET("/", GetSendCount)
	if alertVersion == "ngalert" {
		app.POST("/send", SendMsgNgalert)
	} else if alertVersion == "old" {
		app.POST("/send", SendMsgOld)
	} else {
		fmt.Printf(`[Error] error param "alertversion": %v`, alertVersion)
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
