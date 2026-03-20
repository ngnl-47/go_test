package infra

import (
	"go_test/ordersys/domain/facade"
	"go_test/ordersys/domain/factory"
	infrarepo "go_test/ordersys/infra/repo"
)

// 包级别仓储与 ID 生成器，init 时初始化
var (
	DefaultIDGen       factory.IDGenerator
	DefaultOrderRepo   facade.OrderRepository
	DefaultProductRepo facade.ProductRepository
)

func init() {
	DefaultIDGen = NewUUIDGenerator()
	DefaultOrderRepo = infrarepo.NewInMemoryOrderRepository()
	DefaultProductRepo = infrarepo.NewInMemoryProductRepository()
}
