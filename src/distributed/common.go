// 通用函数单元
package main

import (
	"errors"
	"net"
	"strings"
	"time"
)

//设置一个conn的deadline，单位为秒
func SetDeadLine(conn net.Conn, sec int) {
	d := time.Duration(int64(time.Second) * int64(sec))
	conn.SetDeadline(time.Now().Add(d))
}

//从conn读一个32位无符号整数
func ReadUint32(conn net.Conn, buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, errors.New("The read buffer is not long enought.")
	}
	n, err := conn.Read(buf)
	if err != nil || n != 4 {
		return 0, err
	} else {
		result := (uint32)(buf[0])<<24 + (uint32)(buf[1])<<16 + (uint32)(buf[2])<<8 + (uint32)(buf[3])
		return result, nil
	}
}

//将端口号从IP中删去
func TrimIP(ip string) string {
	i := strings.Index(ip, ":")
	if i >= 0 {
		ip = ip[:i]
	}
	return ip
}
