package utils

import (
	"sync"
	"time"
)

// 雪花算法生成订单号,一秒钟512个ID生成
func GetSnowflakeId() (id int64) {
	worker := NewWorker1()
	return worker.GetId()
}

// 定义一个woker工作节点所需要的基本参数
type Worker1 struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录时间戳
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// 实例化一个工作节点
func NewWorker1() *Worker1 {
	return &Worker1{timestamp: 0, number: 0}
}
func (w *Worker1) GetId() int64 {
	epoch := int64(1670486957) // 设置为去年今天的时间戳...因为位数变了后,几百年都用不完,,实际可以设置上线日期的
	idlength := uint(9)
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().Unix() // 时间戳由原来的毫秒变成了现在的秒
	if w.timestamp == now {
		w.number++
		if w.number > (-1 ^ (-1 << idlength)) { //此处为最大节点ID,大概是2^9-1 511条,
			for now <= w.timestamp {
				now = time.Now().Unix()
			}
		}
	} else {
		w.number = 0
		w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}
	ID := int64((now-epoch)<<idlength | (int64(1) << 1) | (w.number))
	return ID
}
