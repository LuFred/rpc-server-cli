package pool

import (
	"fmt"
	"runtime/debug"
)

func init() {
	cap = 1024
}

type f func()

type GoPool struct {
	JobChannel chan f
	capacity   int
	quit       chan bool
}

var (
	workerPool GoPool //go工作对象
	cap        int    //线程数
)

func newWorker() {
	workerPool = GoPool{
		JobChannel: make(chan f, 10),
		capacity:   cap,
		quit:       make(chan bool)}
}
func run() {
	go func() {
		for {
			select {
			case fun := <-workerPool.JobChannel:
				go func() {
					defer func() {
						if r := recover(); r != nil {
							s := debug.Stack()
							fmt.Errorf("[goroutine panic] error:%s  stack :%s", r, string(s[:]))
						}
					}()
					fun()
				}()
			case <-workerPool.quit:
				break
			}
		}
	}()
}

func Go(fc f) {
	workerPool.JobChannel <- fc
}
func InitGoPool(capacity int) {
	cap = capacity
	newWorker()
	run()

}
