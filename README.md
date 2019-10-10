# sgpool

#### 介绍
协程池

#### 架构
提供两种协程池,可复用/动态创建销毁，debug模式下可以对任务的超时,失败进行统计
#### 安装
```go
go get https://gitee.com/xiawucha365/sgpool
```

#### 使用
具体的worker 需要实现 WorkerInterface 接口
```go
type WorkerInterface interface {
	Task() error
	GetTaskID() interface{}
}
```
##### 常驻复用协程池
```go
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
//例子演示
func main() {

	//创建协程池
    //timeout=0 关闭超时统计 debug=true 打开模式
	spool := internal.NewSPool(3, cap(name_slices_sp),2,true)

	//提交任务
	for _, id := range name_slices_sp {
		np := workersp{ID: id}
		spool.Commit(&np)
	}

	spool.Release()
	time.Sleep(time.Second * 1)
}
```
##### 动态创建销毁协程池
```go
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
    //timeout=0 关闭超时统计 debug=true 打开模式
	wpool = internal.NewWPool(100, cap(name_slices), 0, true)

	//提交任务
	for _, id := range name_slices {
		np := worker{ID: id}
		wpool.Commit(&np)
	}

	wpool.Release()

}
```
