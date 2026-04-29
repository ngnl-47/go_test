package collection

import (
	"strings"
	"time"

	"git.hortorgames.com/orange/mandarin/v3/generics/math"
)

// EqualityComparer 比较器
type EqualityComparer[T any] func(x T, y T) bool

// Comparer 比较器
type Comparer[T any] func(x T, y T) int

// IRawCount 支持直接获取数量
type IRawCount interface {
	RawCount() int
}

// IEnumerator 枚举器
type IEnumerator[T any] interface {
	Next() (item T, ok bool)
}

// Enumerator 枚举器
type Enumerator[T any] func() (item T, ok bool)

// Next 获取下一个元素
func (e Enumerator[T]) Next() (item T, ok bool) {
	return e()
}

// IEnumerable 可枚举
type IEnumerable[T any] interface {
	GetEnumerator() IEnumerator[T]
}

// Enumerable 可枚举
type Enumerable[T any] func() Enumerator[T]

// GetEnumerator 获取枚举器
func (f Enumerable[T]) GetEnumerator() IEnumerator[T] {
	return f()
}

// TimeComparer 时间比较器
func TimeComparer(x time.Time, y time.Time) int {
	return int(x.Sub(y))
}

// StringComparer 字符串比较器
func StringComparer(x string, y string) int {
	return strings.Compare(x, y)
}

// ComparableEqualityComparer 可比较相等性比较器
func ComparableEqualityComparer[T comparable](x T, y T) bool {
	return x == y
}

// SortableCompare 可排序比较器
func SortableCompare[T math.Ordered](x T, y T) int {
	if x == y {
		return 0
	}
	if x < y {
		return -1
	}
	return 1
}

type Predicate[T any] func(item T) bool

//type ICollection[T any] interface {
//	Add(item T)
//	Clear()
//	Contains(item T) bool
//	CopyTo(array []T, arrayIndex int)
//	Remove(item T) bool
//	Count() int
//}
//
//type IList[T any] interface {
//	ICollection[T]
//	IndexOf(item T) int
//	Insert(index int, item T)
//	RemoveAt(index int)
//	GetAt(index int) T
//	SetAt(index int, item T)
//}

//type IMap[TK comparable, TV any] interface {
//	Count() int
//	ContainsKey(key TK) bool
//	Set(key TK, value TV)
//	TryAdd(key TK, value TV) bool
//	Get(key TK) (value TV)
//	TryGet(key TK) (value TV, ok bool)
//	Remove(key TK)
//	TryRemove(key TK) bool
//	Clear()
//}
