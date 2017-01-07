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
const VER string = "v1.0 build 170107171136"

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
	cfg := NewConfig()

	//创建内部监听
	interListener := NewInterListener(cfg)
	go interListener.Listen()

	//创建外部监听
	outerListener := NewOuterListener(cfg, interListener)
	go outerListener.Listen()

	//监听退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	s := <-signalChan
	g_loger.log(nil, C_LOGLEVEL_RUN, "服务器收到退出信息:", s, "，正常退出。")
}
