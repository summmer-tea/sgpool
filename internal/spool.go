package internal

import (
	"fmt"
	"sync"
	"time"
)

///共享协程池
type SPool struct {
	jobQueue    chan WorkerInterface
	wg          sync.WaitGroup
	workerNum   int
	TotalNum    int
	CounterOk   int
	CounterFail int
	CounterOut  int
	mutexFail   sync.Mutex
	mutexOk     sync.Mutex
	mutexOut    sync.Mutex
	TimeStart   int64
	TimeOut     int
	Debug       bool
}

// 协程池
func NewSPool(workerNum int, totalNum int, timeout int, debug bool) *SPool {
	p := SPool{
		workerNum: workerNum,
		TotalNum:  totalNum,
		jobQueue:  make(chan WorkerInterface,workerNum),
		TimeOut:   timeout,
		Debug:     debug,
	}
	//任务开始时间记录
	p.TimeStart = time.Now().Unix()

	p.dispatch()
	return &p
}

func (p *SPool) dispatch() {
	p.wg.Add(p.workerNum)
	for i := 0; i < p.workerNum; i++ {
		//共享协程池
		go func() {
			// 收尾工作 容灾
			defer func() {

				p.wg.Done()
				if err := recover(); err != nil {
					fmt.Println("task run err", err)
				}
			}()

			for w := range p.jobQueue {
				if p.TimeOut > 0 {
					timeout_ch := make(chan interface{})

					go func(wr WorkerInterface) { p.runtaskTimeout(wr, timeout_ch) }(w)

					for {
						select {
						case <-timeout_ch:
							goto forend
						case <-time.After(time.Duration(p.TimeOut) * time.Second):
							p.CountOut()
							fmt.Println(w.GetTaskID(), "timeout")
							goto forend
						}
					}
				forend:
				} else {
					p.runtask(w)
				}
			}

		}()
	}
}

func (p *SPool) runtaskTimeout(wr WorkerInterface, timeout_ch chan interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("task error", err)
		}
	}()
	//执行job里的任务
	err := wr.Task()

	timeout_ch <- "ok"

	if err != nil {
		p.CountFail()
		panic(err)
	} else {
		p.CountOk()
	}

}

func (p *SPool) runtask(wr WorkerInterface) {
	//执行job里的任务
	err := wr.Task()
	if err != nil {
		p.CountFail()
		panic(err)
	} else {
		p.CountOk()
	}

}

// 提交任务
func (p *SPool) Commit(w WorkerInterface) {
	p.jobQueue <- w
}

// 等待组 关闭channel
func (p *SPool) Release() {
	close(p.jobQueue)
	p.wg.Wait()

	if p.Debug {
		p.Runtimelog()
	}

}

//计数器-执行成功
func (p *SPool) CountOk() {
	if p.Debug {
		p.mutexOk.Lock()
		//runtime.Gosched()
		p.CounterOk++
		p.mutexOk.Unlock()
	}


}

//计数器-超时
func (p *SPool) CountOut() {
	if p.Debug {
		p.mutexOut.Lock()
		//runtime.Gosched()
		p.CounterOut++
		p.mutexOut.Unlock()
	}
}

//计数器-失败
func (p *SPool) CountFail() {
	if p.Debug {
		p.mutexFail.Lock()
		//runtime.Gosched()
		p.CounterFail++
		p.mutexFail.Unlock()
	}

}

// log
func (p *SPool) Runtimelog() {
	if p.Debug {
		ttime := MathDecimal(float64(time.Now().Unix() - p.TimeStart))
		trange := MathDecimal(float64(p.TotalNum) / ttime)
		if p.CounterOk > 0 || p.CounterFail > 0 {
			if p.TimeOut > 0 {
				fmt.Println(fmt.Sprintln("runtime:total|fail|timeout:", p.TotalNum, "|", p.CounterFail, "|", p.CounterOut, "", "消耗时间:(", ttime, "s)", "平均:(", trange, "次/s)"))
			} else {
				fmt.Println(fmt.Sprintln("runtime:total|fail:", p.TotalNum, "|", p.CounterFail, "", "消耗时间:(", ttime, "s)", "平均:(", trange, "次/s)"))
			}
		}
	}
}
