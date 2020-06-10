//流控
package main

import "log"

//  限流结构体   最大并发连接数 ，  桶
type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

//ConnLimiter构造函数
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

//取连接
func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation")
		return false
	}

	cl.bucket <- 1 //只需要传一个值进通道
	return true
}

//释放连接
func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket //只需要把Chanel的值传出去即可
	log.Printf("New conection coming: %d", c)
}
