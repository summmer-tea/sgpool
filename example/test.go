package main

import (
	//"fmt"
	//"time"
	"time"
)

func main2(){


	timeout_ch := make(chan interface{})


	go func(){
		for {
			select {
			case <-timeout_ch:
				//任务执行结束
				return
			case <-time.After(time.Duration(3 * time.Second)):
				//打印超时的任务id
				//p.CountFail()
				//if p.Debug {
				//	fmt.Println(wr.GetTaskID(), "timeout")
				//}
				return
			}
		}
	}()

	time.Sleep(time.Second *10)
}