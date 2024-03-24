package main

import (
	"Mointor/config"
	"Mointor/pkg/iptables"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configFilename string
	flag.StringVar(&configFilename, "config", "./config/config.yml", "yaml config filename")
	flag.Parse()
	_, err := configs.LoadConfig(configFilename)
	if err != nil {
		return
	}
	// 创建一个通道来接收通知
	c := make(chan os.Signal, 1)
	// 通知通道接收特定信号
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		iptables.StartRecord()
		// 阻塞等待信号
		select {
		case <-c:
			// 收到信号，准备退出
			fmt.Println("收到退出信号，准备停止程序...")
			return
		default:
			// 没有收到信号，继续执行
		}
	}
}
