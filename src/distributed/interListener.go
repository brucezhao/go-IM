// 内部端口监听
// 赵亦平
// 2016.12.18

package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//每个连接的信息
type ConnectionInfo struct {
	Connected bool //是否有连接
	IP        string
	Port      string //监听的端口
	Count     uint32 //该工作服务器的客户连接数
}

type InterListener struct {
	whiteList          map[string]bool  //白名单，用map是因为检索快
	timeout            int              //超时
	port               string           //监听端口
	serverCount        int              //服务顺的初始数量
	workList           []ConnectionInfo //用于分配IP的服务器信息切片
	realList           []ConnectionInfo //保存了真实情况的服务器信息切片
	maxCount           uint32           //当前所有服务器的最大连接数
	needLockWorkBuffer bool
	mutexReal          sync.Mutex
	mutexWork          sync.Mutex
	syncBufferInterval int
	bufferIndex        int      //当前workList的索引
	currentServerCount int      //当前工作服务器的数量
	chanIndex          chan int //读写WorkList的chan
}

//构造函数
func NewInterListener(cfg *Config) *InterListener {
	var interListener InterListener

	interListener.whiteList = make(map[string]bool, len(cfg.WhiteList))
	for _, v := range cfg.WhiteList {
		interListener.whiteList[v] = true
	}

	interListener.timeout = cfg.Timeout
	interListener.port = cfg.InterPort
	interListener.needLockWorkBuffer = cfg.NeedLockWorkBuffer
	interListener.serverCount = cfg.InterServerCount
	interListener.syncBufferInterval = cfg.SyncBufferInterval
	interListener.workList = make([]ConnectionInfo, interListener.serverCount)
	interListener.realList = make([]ConnectionInfo, interListener.serverCount)
	interListener.bufferIndex = -1
	interListener.currentServerCount = 0
	interListener.chanIndex = make(chan int, 1)

	//启动同步buffer的协程
	go interListener.syncBuffer()
	//启动IP分配协程
	go interListener.distributeWorkIndex()

	return &interListener
}

//开始监听
func (il *InterListener) Listen() {
	listener, err := net.Listen("tcp", il.port)
	if err != nil {
		g_loger.log(nil, C_LOGLEVEL_ERROR, "监听内部端口失败：", err)
		return
	}
	//退出函数
	defer func() {
		recover()
		listener.Close()
	}()

	g_loger.log(nil, C_LOGLEVEL_RUN, "开始监听内部端口：", il.port)

	var sIP string
	for {
		conn, err := listener.Accept()
		if err != nil {
			g_loger.log(nil, C_LOGLEVEL_ERROR, "内部端口出现错误：", err)
			continue
		}
		//判断是否在白名单中
		sIP = conn.RemoteAddr().String()
		sIP = TrimIP(sIP)
		if _, ok := il.whiteList[sIP]; ok {
			go il.handler(conn)
		} else {
			g_loger.log(&conn, C_LOGLEVEL_WARNING, "IP不在白名单中，拒绝连接。")
			conn.Close()
		}
	}
}

//返回一个可用的工作服务器IP
//在没有加锁的情况下，这个函数可能会存在读写的冲突上，导致分配的结果不是很均衡，
//但是这种不均衡的影响不大，而且每隔一分钟是会同步真实的连接数，所以读写冲突产
//生的后果没有什么太大影响
func (il *InterListener) GetServerIP() string {
	g_loger.log(nil, C_LOGLEVEL_RUN, "分配前：")
	g_loger.log(nil, C_LOGLEVEL_RUN, "Real:", il.getServerListInfo(&il.realList))
	g_loger.log(nil, C_LOGLEVEL_RUN, "Work:", il.getServerListInfo(&il.workList))

	if len(il.workList) == 0 {
		return ""
	}

	/*
		if il.needLockWorkBuffer {
			il.mutexWork.Lock()
			defer il.mutexWork.Unlock()
		}

		bi := il.bufferIndex
		bufferLen := len(il.workList)

		if bi < 0 {
			bi = 0
		} else if bi >= bufferLen {
			bi = 0
		}

		//找到一个连接数小于最大连接数的服务器
		var count uint32
		for ; bi < bufferLen; bi++ {
			if il.workList[bi].Connected {
				count = il.workList[bi].Count
				if count < il.maxCount {
					break
				}
			}
		}
		//如果都没有找到，则将maxcount加10,重新分配
		if bi >= bufferLen {
			bi = 0
			count = il.workList[bi].Count
			il.maxCount += 10
		}
		count++

		il.workList[bi].Count = count
	*/
	var bi int = -1
	select {
	case bi = <-il.chanIndex:
	case <-time.After(time.Second):
		//超时退出
		return ""
	}
	if (bi < 0) || (bi >= len(il.workList)) {
		return ""
	}

	g_loger.log(nil, C_LOGLEVEL_RUN, "Index:", bi, ", Result:", il.workList[bi].IP+":"+il.workList[bi].Port)
	g_loger.log(nil, C_LOGLEVEL_RUN, "分配后：")
	g_loger.log(nil, C_LOGLEVEL_RUN, "Real:", il.getServerListInfo(&il.realList))
	g_loger.log(nil, C_LOGLEVEL_RUN, "Work:", il.getServerListInfo(&il.workList))

	return il.workList[bi].IP + ":" + il.workList[bi].Port
}

//链接处理函数
func (il *InterListener) handler(conn net.Conn) {
	il.increaseCurrentServerCount(1)
	g_loger.log(&conn, C_LOGLEVEL_RUN, "工作服务器已连接，服务器数量为：", il.currentServerCount)

	//首先获取可用的链接信息数据结构体
	realIndex := il.getRealIndex()
	pInfo := &(il.realList[realIndex])
	pInfo.Connected = true
	pInfo.IP = TrimIP(conn.RemoteAddr().String())
	pInfo.Port = ""
	pInfo.Count = 0

	//结束处理函数
	defer func() {
		recover()
		conn.Close()
		//将链接信息设置为未连接
		pInfo.Connected = false
		pInfo.IP = ""
		pInfo.Port = ""
		pInfo.Count = 0
		il.increaseCurrentServerCount(-1)
		g_loger.log(&conn, C_LOGLEVEL_RUN, "断开连接，服务器数量为：", il.currentServerCount)
	}()

	//设置超时
	SetDeadLine(conn, il.timeout*2)

	var num uint32
	var err error
	//接收缓冲区，因为只需要接收一个uint32的连接数，所以是4个字节
	readBuffer := make([]byte, 4)

	//首先获取工作服务器发过来的端口号
	num, err = ReadUint32(conn, readBuffer)
	if err != nil {
		g_loger.log(&conn, C_LOGLEVEL_ERROR, err)
		return
	}
	g_loger.log(&conn, C_LOGLEVEL_RUN, "工作服务器端口：", num)
	pInfo.Port = fmt.Sprint(num)

	for {
		//读工作服务器发过来的连接数
		num, err = ReadUint32(conn, readBuffer)
		//g_loger.log(&conn, C_LOGLEVEL_RUN, num, err)
		if err == nil {
			//g_loger.log(&conn, C_LOGLEVEL_RUN, "工作服务器连接数：", num)
			pInfo.Count = num
			//重新设置超时
			SetDeadLine(conn, il.timeout*2)

		} else {
			//出现通讯错误，断开连接，io.EOF是客户端断开时的消息
			if err != io.EOF {
				g_loger.log(&conn, C_LOGLEVEL_ERROR, "收到非法数据:", err)
			}
			return
		}
	}
}

//返回一个可用的realBuffer的索引值
func (il *InterListener) getRealIndex() int {
	il.mutexReal.Lock()
	defer il.mutexReal.Unlock()

	for k, v := range il.realList {
		if v.Connected == false {
			return k
		}
	}
	//如果已经没有可用的buffer，则新增一个
	info := ConnectionInfo{}
	il.realList = append(il.realList, info)

	return len(il.realList)
}

//在workList中找出合适的服务器，并将索引值写入chanIndex中
//该函数由goroutine调用
func (il *InterListener) distributeWorkIndex() {
	var maxCount uint32
	var listLen int
	var index int = 0

	for {
		//如果还没有工作列表，则休眠等待
		if il.currentServerCount <= 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		//g_loger.log(nil, C_LOGLEVEL_RUN, "分配IP之大循环")
		listLen = len(il.workList)
		if index == -1 {
			//说明workList所有的连接数都已经大于maxCount了
			maxCount += 10
		} else {
			maxCount = il.maxCount
		}
		index = -1
		var pv *ConnectionInfo = nil
		for i := 0; i < listLen; i++ {
			pv = &(il.workList[i])
			if pv.Connected && (pv.Count < maxCount) {
				pv.Count++
				//发给chan，如果没有客户端来取，则阻塞等待
				index = i
				il.chanIndex <- i
			}
		}
	}
}

//将realBuffer的数据同步到workBuffer中
func (il *InterListener) syncBuffer() {
	for {
		var maxCount uint32 = 0

		//		g_loger.log(nil, C_LOGLEVEL_RUN, "同步前数据：")
		//		g_loger.log(nil, C_LOGLEVEL_RUN, "Real:", il.getServerListInfo(&il.realList))
		//		g_loger.log(nil, C_LOGLEVEL_RUN, "Work:", il.getServerListInfo(&il.workList))

		for _, v := range il.realList {
			if v.Connected && (v.Count > maxCount) {
				maxCount = v.Count
			}
		}

		if il.needLockWorkBuffer {
			il.mutexWork.Lock()
		}

		//il.workList = il.realList[:]
		il.maxCount = maxCount
		for i := 0; i < len(il.workList); i++ {
			il.workList[i] = il.realList[i]
		}
		for i := len(il.workList); i < len(il.realList); i++ {
			il.workList = append(il.workList, il.realList[i])
		}

		if il.needLockWorkBuffer {
			defer il.mutexWork.Unlock()
		}

		//		g_loger.log(nil, C_LOGLEVEL_RUN, "同步后数据：")
		//		g_loger.log(nil, C_LOGLEVEL_RUN, "Real:", il.getServerListInfo(&il.realList))
		//		g_loger.log(nil, C_LOGLEVEL_RUN, "Work:", il.getServerListInfo(&il.workList))

		//每隔1分钟同步一次
		d := time.Duration(int64(il.syncBufferInterval) * int64(time.Second))
		time.Sleep(d)
	}
}

//增减当前服务器数量
func (il *InterListener) increaseCurrentServerCount(val int) {
	il.mutexReal.Lock()
	defer il.mutexReal.Unlock()

	il.currentServerCount += val
}

//记录服务器的信息
func (il *InterListener) getServerListInfo(list *[]ConnectionInfo) string {
	var str string = ""

	for _, v := range *list {
		if v.Connected {
			str += fmt.Sprintf("%s=%d,", v.Port, v.Count)
		}
	}
	return str
}
