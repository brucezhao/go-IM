// 分发服务器v1.0
// 赵亦平
// 2016.12.17

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//版本号
const VER string = "v1.0 build 1612180030"

//全局日志变量
var g_loger *IMLog

func main() {
	var err error

	//创建日志
	g_loger, err = NewLog()
	if err != nil {
		fmt.Println(err)
		return
	}

	//显示版本信息
	g_loger.log(nil, C_LOGLEVEL_RUN, "服务器启动，版本：", VER)

	//读配置文件
	cfg, err := NewConfig()
	if err != nil {
		g_loger.log(nil, C_LOGLEVEL_ERROR, err)
		return
	}

	//创建内部监听
	NewInterListener(cfg)

	//创建外部监听
	NewOuterListener(cfg)

	//监听退出信息
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	s := <-signalChan
	g_loger.log(nil, C_LOGLEVEL_RUN, "服务器收到退出信息:", s, "，正常退出。")
}
