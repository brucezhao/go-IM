// 读配置文件单元
// 赵亦平
// 2016.12.17

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const C_DEFAULT_INTERPORT = ":12170"  //默认的内部端口
const C_DEFAULT_OUTERPORT = ":12171"  //默认的外部端口
const C_DEFAULT_INTERSERVERCOUNT = 10 //默认内部工作服务器数量

type Config struct {
	InterPort        string   //内部端口，面向工作服务器的端口
	OuterPort        string   //外部端口，面向客户端的端口
	WhiteList        []string //白名单，允许连接的工作服务器的IP列表
	BlackList        []string //黑名单，不允许连接的外部客户端的IP列表
	InterServerCount int      //InterServerCount只是用来初始化服务器数组的值，后期会动态调整，
	//但是为了尽量减少内存分配，初期根据实际情况，设置一个比较合适的值
	LogLevel int //日志级别
}

//构造函数
func NewConfig() (*Config, error) {
	buf, err := ioutil.ReadFile(os.Args[0] + ".json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(buf, &cfg)

	if cfg.InterPort == "" {
		cfg.InterPort = C_DEFAULT_INTERPORT
	}
	if cfg.OuterPort == "" {
		cfg.OuterPort = C_DEFAULT_OUTERPORT
	}
	if cfg.InterServerCount == 0 {
		cfg.InterServerCount = C_DEFAULT_INTERSERVERCOUNT
	}

	return &cfg, nil
}
