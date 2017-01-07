// 外部端口监听
package main

import (
	"net"
)

type OuterListener struct {
	port          string          //监听端口
	blackList     map[string]bool //黑名单
	interListener *InterListener  //内部监听对象
}

func NewOuterListener(cfg *Config, il *InterListener) *OuterListener {
	var listener OuterListener

	listener.interListener = il
	listener.blackList = make(map[string]bool, len(cfg.BlackList))
	for _, v := range cfg.BlackList {
		listener.blackList[v] = true
	}
	listener.port = cfg.OuterPort

	return &listener
}

func (ol *OuterListener) Listen() {
	listener, err := net.Listen("tcp", ol.port)
	if err != nil {
		g_loger.log(nil, C_LOGLEVEL_ERROR, "监听外部端口失败：", err)
		return
	}

	defer func() {
		recover()
		listener.Close()
	}()

	g_loger.log(nil, C_LOGLEVEL_RUN, "开始监听外部端口：", ol.port)

	var sIP string
	for {
		conn, err := listener.Accept()
		if err != nil {
			g_loger.log(nil, C_LOGLEVEL_ERROR, "外部端口出现错误：", err)
			continue
		}

		sIP = TrimIP(conn.RemoteAddr().String())

		if _, ok := ol.blackList[sIP]; !ok {
			go ol.handler(conn)
		} else {
			conn.Close()
		}
	}
}

//链接处理函数
func (ol *OuterListener) handler(conn net.Conn) {
	defer conn.Close()

	s := ol.interListener.GetServerIP()

	if s == "" {
		s = ":"
	}

	conn.Write([]byte(s))
}
