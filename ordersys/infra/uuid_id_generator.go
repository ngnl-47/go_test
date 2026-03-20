package infra

import (
	"fmt"

	"go_test/ordersys/domain/factory"
)

// UUIDGenerator UUID生成器（简化实现，可替换为 google/uuid 等）
type UUIDGenerator struct {
	counter int
}

// NewUUIDGenerator 创建UUID生成器
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{
		counter: 0,
	}
}

// Generate 生成唯一ID
func (g *UUIDGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("ID-%d", g.counter)
}

var _ factory.IDGenerator = (*UUIDGenerator)(nil)
