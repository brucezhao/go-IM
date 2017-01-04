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
const VER string = "v1.0 build 1612272219"

//全局日志变量
var g_loger *IMLog

func main() {
	//test code
	//	c := NewConnections(3)

	//	for i := 0; i < 3; i++ {
	//		m := ConnectionInfo{true, "127.0.0.1", "12171", uint32(i)}
	//		//		fmt.Println("m=", m)
	//		c.Append(m)
	//	}
	//	//	p := c.At(1)
	//	fmt.Println(c)
	//	p := c.At(1)
	//	fmt.Println(*p)

	//	for i := 3; i < 5; i++ {
	//		m := ConnectionInfo{true, "127.0.0.1", "12171", uint32(i)}
	//		//		fmt.Println("m=", m)
	//		c.Append(m)
	//	}
	//	fmt.Println(c)

	//	c1 := NewConnections(3)
	//	c1.Copy(c)
	//	fmt.Println(c1)

	//	p.Connected = false
	//	p.Count = 100
	//	fmt.Println(c)

	//	return

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
