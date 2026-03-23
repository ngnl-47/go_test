package tests

import (
	"testing"
	"time"
)

type futureTask struct {
	name string
	fn   func()
}

func (f futureTask) run(t *testing.T) chan int32 {
	//ret := make(chan int32)
	ret := make(chan int32, 1) // 区别在于，带缓冲的chan，不会阻塞子协程，写入chan后，会立即返回
	go func() {
		t.Log(f.name, "started!")
		f.fn()
		ret <- 1
		ret <- 2
		t.Log(f.name, "finished!")
	}()
	return ret
}

// 模拟 future 特性
func TestFutureFeat(t *testing.T) {
	task1 := futureTask{
		name: "task1",
		fn: func() {
			time.Sleep(3 * time.Second)
		},
	}
	// 启动 task1
	task1Ret := task1.run(t)
	// 执行 task2
	t.Log("task2 started!")
	t.Log("task2 finished!")
	// 等待 task1 完成，阻塞
	code := <-task1Ret
	t.Log("task1 code:", code)
	t.Log("all tasks finished!")
}
