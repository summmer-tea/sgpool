

package internal

import (
	"strconv"
	"testing"
)


const(
	WORKER_NUM_10000 = 10000
	WORKER_NUM_1000 = 1000
	WORKER_NUM_100 = 100
)

type workersp struct {
	ID string
}


func (m *workersp) Task() error {

	//fmt.Println("job:" + m.ID + "runing...")
	//timen := rand.Intn(3)
	//fmt.Println(timen,"seconds")
	//time.Sleep(time.Second * time.Duration(timen))
	//fmt.Println("job:" + m.ID + "over")
	return nil
}


//获取任务id,便于
func (m *workersp) GetTaskID() interface{} {
	return m.ID
}



func BenchmarkSPool(b *testing.B) {

	//创建协程池
	spool := NewSPool(WORKER_NUM_100,  b.N,0,false)

	//提交任务
	for n := 0; n < b.N; n++ {
		np := workersp{ID: strconv.Itoa(n)}
		spool.Commit(&np)
	}

	spool.Release()
}


func BenchmarkWPool(b *testing.B) {

	//创建协程池
	spool := NewWPool(WORKER_NUM_100,  b.N,0,false)

	//提交任务
	for n := 0; n < b.N; n++ {
		np := workersp{ID: strconv.Itoa(n)}
		spool.Commit(&np)
	}

	spool.Release()
}