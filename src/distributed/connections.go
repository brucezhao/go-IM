// connections
// 连接管理
// 赵亦平
// 2017.1.4

package main

import (
	"fmt"
)

//每个连接的信息
type ConnectionInfo struct {
	Connected bool //是否有连接
	IP        string
	Port      string //监听的端口
	Count     uint32 //该工作服务器的客户连接数
}

type ConnectionInfoSlice []ConnectionInfo

type Connections struct {
	datas           map[int]*ConnectionInfoSlice
	currentMapIndex int
	capacity        int
	currentSlice    *ConnectionInfoSlice
	count           int
}

//构造函数
func NewConnections(capacity int) *Connections {
	connections := Connections{nil, 0, capacity, nil, 0}

	connections.datas = make(map[int]*ConnectionInfoSlice)
	data := make(ConnectionInfoSlice, 0, capacity)
	connections.datas[connections.count] = &data
	connections.currentSlice = &data

	return &connections
}

//返回可打印内容
func (c *Connections) String() string {
	var cs ConnectionInfoSlice = nil
	var sResult = "["
	var pCi *ConnectionInfo = nil

	for i := 0; i < c.currentMapIndex+1; i++ {
		cs = *c.datas[i]
		for l := 0; l < len(cs); l++ {
			pCi = &cs[l]
			sResult += fmt.Sprint(*pCi)
		}
	}

	sResult += "]"
	return sResult
}

//添加一个链接
func (c *Connections) Append(ci ConnectionInfo) *Connections {
	if len(*c.currentSlice) == cap(*c.currentSlice) {
		data := make(ConnectionInfoSlice, 0, c.capacity)
		c.currentSlice = &data
		c.currentMapIndex++
		c.datas[c.currentMapIndex] = &data

	}
	*c.currentSlice = append(*c.currentSlice, ci)
	c.count++

	return c
}

//链接信息总数，包括未连接的
func (c *Connections) Count() int {
	return c.count
}

//按索引取值
func (c *Connections) At(index int) *ConnectionInfo {
	//首先确定是哪个slice
	i1 := index / c.capacity
	if i1 > c.currentMapIndex {
		return nil
	}

	//再在找到的slice中取值
	i2 := index % c.capacity
	if (i1 == c.currentMapIndex) && (i2 >= len(*c.currentSlice)) {
		return nil
	}

	pCi := &(*(c.datas[i1]))[i2]
	return pCi
}

//深度拷贝
func (c *Connections) Copy(src *Connections) *Connections {
	//首先清空原有内容
	for _, v := range c.datas {
		*v = nil
	}
	c.datas = nil

	c.datas = make(map[int]*ConnectionInfoSlice)
	c.capacity = src.capacity

	var pCs *ConnectionInfoSlice = nil
	for i := 0; i <= src.currentMapIndex; i++ {
		pCs = src.datas[i]
		data := make(ConnectionInfoSlice, len(*pCs), c.capacity)
		copy(data, *pCs)
		//fmt.Println(*pCs)
		//c.datas = append(c.datas, &data)
		c.datas[i] = &data

		if i == src.currentMapIndex {
			c.currentSlice = &data
		}
	}

	c.currentMapIndex = src.currentMapIndex
	c.count = src.count

	return c
}
