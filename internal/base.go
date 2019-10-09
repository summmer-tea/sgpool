package internal

import (
	"fmt"
	"strconv"
)

//协程池接口
//type PoolInterface interface {
//	Commit(w WorkerInterface)
//}



type job func()

//任务接口
type WorkerInterface interface {
	Task() error
	GetTaskID() interface{}
}



//浮点保留2位小数
func  MathDecimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}