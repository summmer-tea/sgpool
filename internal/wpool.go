package internal

import (
	"fmt"
	"sync"
	"time"
)

//非共享协程池
type WPool struct {
	//任务传递
	jobQueue chan WorkerInterface
	//并发数控制
	limitChan   chan interface{}
	wg          sync.WaitGroup
	TotalNum    int
	CounterOk   int
	CounterFail int
	CounterOut int
	mutexFail   sync.Mutex
	mutexOk     sync.Mutex
	mutexOut     sync.Mutex
	TimeStart   int64
	TimeOut     int
	Debug       bool

}

// 协程池

func NewWPool(workerNum int, totalNum int,timeout int,debug bool) *WPool {

	p := WPool{
		TotalNum:  totalNum,
		jobQueue:  make(chan WorkerInterface),
		limitChan: make(chan interface{}, workerNum),
		TimeOut:timeout,
		Debug:debug,
	}
	p.wg.Add(totalNum)
	p.dispatch()
	return &p
}

// 提交任务
func (p *WPool) Commit(w WorkerInterface) {

	p.limitChan <- "ok"
	p.jobQueue <- w

}

// 控制最大并发数
func (p *WPool) dispatch() {
	//任务开始时间记录
	p.TimeStart = time.Now().Unix()

	//新起一个协程
	go func() {
		for w := range p.jobQueue {

			go func(wr WorkerInterface) {

				// 收尾工作 容灾
				defer func() {
					p.wg.Done()
					<-p.limitChan
					if err := recover(); err != nil {
						fmt.Println("task run err", err)
					}
				}()


				if p.TimeOut > 0{
					//增加超时任务统计
					timeout_ch := make(chan interface{})

					go func() { p.runtaskTimeout(wr, timeout_ch) }()

					for {
						select {
						case <-timeout_ch:
							//任务执行结束
							return
						case <-time.After(time.Duration(p.TimeOut) * time.Second):
							//打印超时的任务id
							p.CountOut()
							if p.Debug {
								fmt.Println(wr.GetTaskID(), "timeout")
							}
							return
						}
					}
				}else{
					p.runtask(wr)
				}


			}(w)

		}
	}()

}

func (p *WPool) runtaskTimeout(wr WorkerInterface, timeout_ch chan interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("task run err", err)
		}
	}()

	//执行job里的任务
	err := wr.Task()

	timeout_ch <- "ok"

	if err == nil {
		p.CountOk()
	} else {
		p.CountFail()
		panic(err)
	}
}

func (p *WPool) runtask(wr WorkerInterface) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("task run err", err)
		}
	}()

	//执行job里的任务
	err := wr.Task()

	if err == nil {
		p.CountOk()
	} else {
		p.CountFail()
		panic(err)
	}
}

// 等待组 关闭channel
func (p *WPool) Release() {
	p.wg.Wait()
	close(p.jobQueue)
	close(p.limitChan)

	if p.Debug{
		p.Runtimelog()
	}

}

//计数器-执行成功
func (p *WPool) CountOk() {
	p.mutexOk.Lock()
	//runtime.Gosched()
	p.CounterOk++
	p.mutexOk.Unlock()

}

//计数器-超时
func (p *WPool) CountOut() {
	if p.Debug{
		p.mutexOut.Lock()
		//runtime.Gosched()
		p.CounterOut++
		p.mutexOut.Unlock()
	}
}

//计数器-失败
func (p *WPool) CountFail() {
	if p.Debug{
		p.mutexFail.Lock()
		//runtime.Gosched()
		p.CounterFail++
		p.mutexFail.Unlock()
	}


}

// log
func (p *WPool) Runtimelog() {
	if p.Debug{
		ttime := MathDecimal(float64(time.Now().Unix() - p.TimeStart))
		trange := MathDecimal(float64(p.TotalNum) / ttime)
		if p.CounterOk > 0 || p.CounterFail > 0 {
			if p.TimeOut >0 {
				fmt.Println(fmt.Sprintln("runtime:total|fail|timeout:", p.TotalNum, "|",p.CounterFail,"|", p.CounterOut, "", "消耗时间:(", ttime, "s)", "平均:(", trange, "次/s)"))
			}else{
				fmt.Println(fmt.Sprintln("runtime:total|fail:", p.TotalNum, "|", p.CounterFail, "", "消耗时间:(", ttime, "s)", "平均:(", trange, "次/s)"))
			}
		}
	}
}
