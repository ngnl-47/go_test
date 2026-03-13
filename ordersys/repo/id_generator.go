package repo

import (
	"fmt"
	"go_test/ordersys/domain/factory"
)

// UUIDGenerator UUID生成器
// 基础设施层：实现领域层定义的IDGenerator接口
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
// 简化实现，实际项目可以使用google/uuid等库
func (g *UUIDGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("ID-%d", g.counter)
}

var _ factory.IDGenerator = (*UUIDGenerator)(nil)
