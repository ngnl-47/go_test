package collection

const (
	queue_MinimumGrow     int = 4
	queue_ShrinkThreshold int = 0x20
	queue_GrowFactor      int = 200
	queue_DefaultCapacity int = 4
)

// Queue 是一个队列
type Queue[T any] struct {
	// Fields
	_array   []T
	_head    int
	_tail    int
	_size    int
	_version int
	comparer EqualityComparer[T]
}

// Methods
func NewQueue[T comparable](capacity int) *Queue[T] {
	if capacity < 0 {
		panic("ArgumentOutOfRange_NeedNonNegNumRequired")
	}
	q := &Queue[T]{}
	q._array = make([]T, capacity)
	return NewQueueComparable(capacity, ComparableEqualityComparer[T])
}

// NewQueueComparable 创建一个可比较的队列
func NewQueueComparable[T any](capacity int, comparer EqualityComparer[T]) *Queue[T] {
	if capacity < 0 {
		panic("ArgumentOutOfRange_NeedNonNegNumRequired")
	}
	q := &Queue[T]{}
	q._array = make([]T, capacity)
	q.comparer = comparer
	return q
}

// Clear 清空队列
func (q *Queue[T]) Clear() {
	var empty T
	if q._head < q._tail {
		start := q._head
		end := start + q._size
		for i := start; i < end; i++ {
			q._array[i] = empty
		}
	} else {
		start := q._head
		end := start + len(q._array) - q._head
		for i := start; i < end; i++ {
			q._array[i] = empty
		}
		for i := start; i < q._tail; i++ {
			q._array[i] = empty
		}
	}
	q._head = 0
	q._tail = 0
	q._size = 0
	q._version++
}

// Contains 检查队列是否包含指定元素
func (q *Queue[T]) Contains(item T) bool {
	index := q._head
	num2 := q._size
	for num2 > 0 {
		num2--
		if q.comparer(q._array[index], item) {
			return true
		}
		index = (index + 1) % len(q._array)
	}
	return false
}

// CopyTo 复制队列到数组
func (q *Queue[T]) CopyTo(array []T, arrayIndex int) {
	if array == nil {
		panic("Array is nil")
	}
	if (arrayIndex < 0) || (arrayIndex > len(array)) {
		panic("ArgumentOutOfRange_Index")
	}
	arrayLen := len(array)
	if (arrayLen - arrayIndex) < q._size {
		panic("Argument_InvalidOffLen")
	}
	numToCopy := 0
	if (arrayLen - arrayIndex) < q._size {
		numToCopy = arrayLen - arrayIndex
	} else {
		numToCopy = q._size
	}
	if numToCopy == 0 {
		return
	}
	firstPart := 0
	if (len(q._array) - q._head) < numToCopy {
		firstPart = len(q._array) - q._head
	} else {
		firstPart = numToCopy
	}
	copy(array[arrayIndex:], q._array[q._head:q._head+firstPart])
	numToCopy -= firstPart
	if numToCopy > 0 {
		copy(array[(arrayIndex+len(q._array))-q._head:], q._array[:numToCopy])
	}
}

// Dequeue 出队
func (q *Queue[T]) Dequeue() T {
	if q._size == 0 {
		panic("InvalidOperation_EmptyQueue")
	}
	removed := q._array[q._head]
	var empty T
	q._array[q._head] = empty
	q._head = (q._head + 1) % len(q._array)
	q._size--
	q._version++
	return removed
}

// Enqueue 入队
func (q *Queue[T]) Enqueue(item T) {
	if q._size == len(q._array) {
		newcapacity := len(q._array) * queue_GrowFactor / 100
		if newcapacity < len(q._array)+queue_MinimumGrow {
			newcapacity = len(q._array) + queue_MinimumGrow
		}
		q.SetCapacity(newcapacity)
	}
	q._array[q._tail] = item
	q._tail = (q._tail + 1) % len(q._array)
	q._size++
	q._version++
}

// getElement 获取指定索引的元素
func (q *Queue[T]) getElement(i int) T {
	return q._array[(q._head+i)%len(q._array)]
}

// Peek 获取队列顶部的元素
func (q *Queue[T]) Peek() T {
	if q._size == 0 {
		panic("InvalidOperation_EmptyQueue")
	}
	return q._array[q._head]
}

// SetCapacity 设置队列容量
func (q *Queue[T]) SetCapacity(capacity int) {
	newarray := make([]T, capacity)
	if q._size > 0 {
		if q._head < q._tail {
			copy(newarray[:q._size], q._array[q._head:])
		} else {
			copy(newarray, q._array[q._head:])
			copy(newarray[len(q._array)-q._head:], q._array[:q._tail])
		}
	}
	q._array = newarray
	q._head = 0
	if q._size == capacity {
		q._tail = 0
	} else {
		q._tail = q._size
	}
	q._version++
}

// ToArray 转换为数组
func (q *Queue[T]) ToArray() []T {
	arr := make([]T, q._size)
	if q._size == 0 {
		return arr
	}
	if q._head < q._tail {
		copy(arr, q._array[q._head:])
	} else {
		copy(arr, q._array[q._head:])
		copy(arr[len(q._array)-q._head:], q._array[:q._tail])
	}
	return arr
}

// TrimExcess 修剪队列
func (q *Queue[T]) TrimExcess() {
	threshold := len(q._array) * 9 / 10
	if q._size < threshold {
		q.SetCapacity(q._size)
	}
}

// Count 返回队列中的元素数量
func (q *Queue[T]) Count() int {
	return q._size
}

// Range 遍历队列中的所有元素
func (q *Queue[T]) Range(f func(value T) bool) {
	idx := 0
	for idx < q._size {
		if f(q.getElement(idx)) {
			break
		}
		idx++
	}
}

// GetEnumerator 获取队列的枚举器
func (q *Queue[T]) GetEnumerator() IEnumerator[T] {
	version := q._version
	index := -1
	return Enumerator[T](func() (item T, ok bool) {
		if version != q._version {
			panic("InvalidOperation_EnumFailedVersion")
		}
		if index == -2 {
			return
		}
		index++
		if index == q._size {
			index = -2
			return
		}
		item = q.getElement(index)
		ok = true
		return
	})
}
