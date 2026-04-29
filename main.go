package main

import (
	"fmt"
	"runtime"
	"time"
)

func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("当前堆内存占用: %d KB, 活跃 Timer 数量: %d\n", m.Alloc/1024, runtime.NumGoroutine())
}

func main() {
	ch := make(chan int)

	// 启动一个生产者，飞快地发消息，让 select 永远选不到超时分支
	go func() {
		for {
			ch <- 1
			time.Sleep(time.Millisecond) // 每毫秒发一个
		}
	}()

	fmt.Println("--- 开始测试：错误示范 (time.After) ---")
	for i := 0; i < 10000; i++ {
		select {
		case <-ch:
			// 处理数据
		case <-time.After(1 * time.Hour): // ⚠️ 每次循环都创建但不触发
			return
		}
		if i%2000 == 0 {
			printMem()
		}
	}

	// 此时你会发现内存持续上涨，因为那 10000 个 1小时的定时器全挂在内存里
	fmt.Println("测试结束，你可以看到内存并没有被回收，因为定时器还没到期。")
}
