package main

import (
	"fmt"
	"math/rand"
	"sgpool/internal"
	"time"
)





//======================woker实现===start=====================\\
type workersp struct {
	ID string
}

//要执行的任务列表

var name_slices_sp = []string{"001","002","003","004","005","006","007","008","009"}


func (m *workersp) Task() error {

	//fmt.Println("job:" + m.ID + "runing...")
	timen := rand.Intn(3)
	//fmt.Println(timen,"seconds")
	time.Sleep(time.Second * time.Duration(timen))
	fmt.Println("job:" + m.ID + "over")
	return nil
}


//获取任务id,便于
func (m *workersp) GetTaskID() interface{} {
	return m.ID
}

//======================woker实现===end=====================\\

//并发池实例




//例子演示
func main() {

	//创建协程池
	spool := internal.NewSPool(3, cap(name_slices_sp),2,true)

	//提交任务
	for _, id := range name_slices_sp {
		np := workersp{ID: id}
		spool.Commit(&np)
	}

	spool.Release()
	time.Sleep(time.Second * 1)
}