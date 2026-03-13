package repo

import (
	"go_test/ordersys/domain/model/agg"
	"go_test/ordersys/domain/model/entity"
)

// OrderRepository 订单仓储接口
// 仓储模式特性：
// 1. 领域层定义接口，基础设施层实现
// 2. 提供聚合根的持久化操作
// 3. 隐藏数据访问细节
// 4. 只针对聚合根定义仓储
type OrderRepository interface {
	Save(order *agg.Order) error
	FindByID(id string) (*agg.Order, error)
	FindByUserID(userID string) ([]*agg.Order, error)
	FindAll() ([]*agg.Order, error)
	Delete(id string) error
}

// ProductRepository 产品仓储接口
type ProductRepository interface {
	Save(product *entity.Product) error
	FindByID(id string) (*entity.Product, error)
	FindAll() ([]*entity.Product, error)
	Delete(id string) error
}
