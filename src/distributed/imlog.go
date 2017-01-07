// 日志单元
// 赵亦平
// 2016.12.17
// 此外我是用了linux的syslog，loglevel作用不大，主要是为了将来使用其它方式时做保留

package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net"
)

const C_LOGLEVEL_RUN = 0
const C_LOGLEVEL_WARNING = 1
const C_LOGLEVEL_ERROR = 2
const C_LOGLEVEL_DEBUG = 3

var LOGLEVELSTR [4]string = [4]string{"Run", "Warning", "Error", "Debug"}

type IMLog struct {
	loger *log.Logger
}

func NewLog() (*IMLog, error) {
	var loger IMLog
	_loger, err := syslog.NewLogger(syslog.LOG_INFO, 0)
	if err != nil {
		return nil, err
	}
	loger.loger = _loger

	return &loger, nil
}

func (il *IMLog) log(conn *net.Conn, level int, args ...interface{}) {
	sIP := ""

	if conn != nil {
		sIP = (*conn).RemoteAddr().String()
	}

	sAtt := fmt.Sprint(args)
	il.loger.Println(LOGLEVELSTR[level], sIP, sAtt)
}
