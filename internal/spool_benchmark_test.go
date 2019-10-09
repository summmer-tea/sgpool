

package internal

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)


type workersp struct {
	ID string
}


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



func BenchmarkSPool(b *testing.B) {

	//创建协程池
	spool := NewSPool(1,  b.N,0,false)

	//提交任务
	for n := 0; n < b.N; n++ {
		np := workersp{ID: strconv.Itoa(n)}
		spool.Commit(&np)
	}

	spool.Release()
}
