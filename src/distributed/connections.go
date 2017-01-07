// connections
// 连接管理
// 赵亦平
// 2017.1.4

package main

import (
	"fmt"
	"sort"
)

//每个连接的信息
type ConnectionInfo struct {
	Connected bool   //是否有连接
	InterIP   string //内部连接的IP:Port
	OuterIP   string //对外监听的IP:Port
	Count     uint32 //该工作服务器的客户连接数
}

type ConnectionInfoSlice []ConnectionInfo

/*sort.Interface */
func (c ConnectionInfoSlice) Len() int {
	return len(c)
}

func (c ConnectionInfoSlice) Less(i, j int) bool {
	return c[i].Count < c[j].Count
}

func (c ConnectionInfoSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

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

//深度拷贝,c必须是已经创建对象，并且capacity必须一致
func (c *Connections) Copy(src *Connections) *Connections {
	//根据本项目的特点，不考虑c.count大于src.count的情况
	for i := 0; i < c.count; i++ {
		*(c.At(i)) = *(src.At(i))
	}
	for i := c.count; i < src.count; i++ {
		c.Append(*(src.At(i)))
	}

	c.currentMapIndex = src.currentMapIndex
	c.currentSlice = c.datas[c.currentMapIndex]
	c.count = src.count

	return c
}

//排序
func (c *Connections) Sort() {
	if c.count <= 1 {
		//少于两个元素，不用排序
		return
	}

	tempSlice := make(ConnectionInfoSlice, c.count, c.count)
	for i := 0; i < c.count; i++ {
		tempSlice[i] = *c.At(i)
	}
	sort.Sort(tempSlice)
	for i := 0; i < c.count; i++ {
		*c.At(i) = tempSlice[i]
	}
}
