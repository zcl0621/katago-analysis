package main

import (
	"ktago-server-go/logger"
	"ktago-server-go/process"
	"ktago-server-go/server"
	"os"
)

func main() {
	process.Init()
	if err := server.StartGinServer(); err != nil {
		logger.Logger("main 启动http服务失败", "error", err, "")
		os.Exit(2)
		return
	}
}
