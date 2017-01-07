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
const C_DEFAULT_TIMEOUT = 60          //默认的超时时间，单位秒
const C_DEFAULT_CHANBUFFER = 10       //IP地址分配时的缓冲区
//const C_DEFAULT_SYNCINTERVAL = 60     //默认同步列表的间隔，单位秒

type Config struct {
	InterPort        string   //内部端口，面向工作服务器的端口
	OuterPort        string   //外部端口，面向客户端的端口
	WhiteList        []string //白名单，允许连接的工作服务器的IP列表
	BlackList        []string //黑名单，不允许连接的外部客户端的IP列表
	InterServerCount int      //InterServerCount只是用来初始化服务器数组的值，后期会动态调整，
	//但是为了尽量减少内存分配，初期根据实际情况，设置一个比较合适的值
	LogLevel int //日志级别
	Timeout  int //链接超时时间
	//SyncBufferInterval int  //同步真实列表与工作列表的时间间隔，秒为单位
	//	NeedLockWorkBuffer bool //该参数影响分配给客户端IP的机制，加锁降低性能，但是服务器更均衡，不加锁提升性能，但服务器的均衡性可能会降低，默认不加锁
	ChanBuffer int
}

//构造函数
func NewConfig() *Config {
	var cfg Config

	buf, err := ioutil.ReadFile(os.Args[0] + ".json")
	if err == nil {
		err = json.Unmarshal(buf, &cfg)
	}

	if cfg.InterPort == "" {
		cfg.InterPort = C_DEFAULT_INTERPORT
	}
	if cfg.OuterPort == "" {
		cfg.OuterPort = C_DEFAULT_OUTERPORT
	}
	if cfg.InterServerCount == 0 {
		cfg.InterServerCount = C_DEFAULT_INTERSERVERCOUNT
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = C_DEFAULT_TIMEOUT
	}
	//	if cfg.SyncBufferInterval == 0 {
	//		cfg.SyncBufferInterval = C_DEFAULT_SYNCINTERVAL
	//	}
	if cfg.ChanBuffer == 0 {
		cfg.ChanBuffer = C_DEFAULT_CHANBUFFER
	}

	return &cfg
}
