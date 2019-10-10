package main

import (
	"fmt"
	"math/rand"
	"sgpool/internal"
	"time"
)

//======================woker实现===start=====================\\
type worker struct {
	ID string
}

//要执行的任务列表

var name_slices = []string{"001", "002", "003", "004", "005", "006", "007", "008", "009", "010", "011"}

func (m *worker) Task() error {

	fmt.Println("job:" + m.ID + "runing...")
	timen := rand.Intn(3)
	//fmt.Println(timen,"seconds")
	time.Sleep(time.Second * time.Duration(timen))
	fmt.Println("job:" + m.ID + "over")
	return nil
}

//获取任务id,便于
func (m *worker) GetTaskID() interface{} {
	return m.ID
}

//======================woker实现===end=====================\\

//并发池实例
var wpool *internal.WPool

//例子演示
func main() {

	//创建协程池
	wpool = internal.NewWPool(100, cap(name_slices), 0, true)

	//提交任务
	for _, id := range name_slices {
		np := worker{ID: id}
		wpool.Commit(&np)
	}

	wpool.Release()

}
