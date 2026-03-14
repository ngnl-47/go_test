package repo

import (
	domainrepo "go_test/ordersys/domain/repo"
	"go_test/ordersys/domain/factory"
)

// 包级别仓储与 ID 生成器，init 时初始化
var (
	DefaultIDGen       factory.IDGenerator
	DefaultOrderRepo   domainrepo.OrderRepository
	DefaultProductRepo domainrepo.ProductRepository
)

func init() {
	DefaultIDGen = NewUUIDGenerator()
	DefaultOrderRepo = NewInMemoryOrderRepository()
	DefaultProductRepo = NewInMemoryProductRepository()
}
