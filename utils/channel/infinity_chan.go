package channel

import (
	"git.hortorgames.com/orange/mandarin/v3/generics/collection"
)

type InfinityChan[T any] struct {
	in     chan T
	out    chan T
	buffer *collection.Queue[T]
	notify chan struct{}
}

func NewInfinityChan[T any](chanCap int, bufferCap int) *InfinityChan[T] {
	c := &InfinityChan[T]{
		in:     make(chan T, 8),
		out:    make(chan T, chanCap),
		notify: make(chan struct{}, 1),
		buffer: collection.NewQueueComparable[T](bufferCap, func(x T, y T) bool {
			return false
		}),
	}
	go c.run()
	return c
}

// Send 添加，如果chan已关闭会panic
func (c *InfinityChan[T]) Send(v T) {
	c.in <- v
}

// TrySend 尝试添加，如果chan已关闭返回false
func (c *InfinityChan[T]) TrySend(v T) (ok bool) {
	defer func() {
		recover()
	}()
	c.in <- v
	ok = true
	return
}

func (c *InfinityChan[T]) Recv() (v T, ok bool) {
	v, ok = <-c.out
	if ok {
		c.Notify()
	}
	return
}

// 直接使用out需要在拿到数据后调用Notify
func (c *InfinityChan[T]) Out() <-chan T {
	return c.out
}

func (c *InfinityChan[T]) trySendToOut(v T) bool {
	select {
	case c.out <- v:
		return true
	default:
		return false
	}
}

func (c *InfinityChan[T]) Notify() {
	defer func() {
		recover()
	}()
	select {
	case c.notify <- struct{}{}:
	default:
	}
}

func (c *InfinityChan[T]) run() {
	closed := false
	for !closed {
		select {
		case v, ok := <-c.in:
			if !ok {
				closed = true
				break
			}
			// 如果队列为空，尝试直接写入out
			if c.buffer.Count() == 0 && c.trySendToOut(v) {
				continue
			}
			// 如果队列非空或out已满,写入队列
			c.buffer.Enqueue(v)
			c.Notify()
		case <-c.notify:
			// 尝试将队列中的数据写入out
			for c.buffer.Count() > 0 && c.trySendToOut(c.buffer.Peek()) {
				c.buffer.Dequeue()
			}
		}
	}
	// close后将队列中剩余的数据写入out
	for c.buffer.Count() > 0 {
		c.out <- c.buffer.Dequeue()
	}
	c.buffer = nil
	close(c.out)
	close(c.notify)
}

// Close 关闭
// dropLeft 是否丢弃未消费的数据，false时，需保证有其他消费者消费完剩余数据，否则会有内存泄漏
func (c *InfinityChan[T]) Close(dropLeft bool) {
	close(c.in)
	if dropLeft {
		GoPooled(func() {
			for {
				if _, ok := c.Recv(); !ok {
					return
				}
			}
		})
	}
}

var GoPooled func(func())
