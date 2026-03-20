package facade

import "go_test/ordersys/domain/model"

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Save(order *model.Order) error
	FindByID(id string) (*model.Order, error)
	FindByUserID(userID string) ([]*model.Order, error)
	FindAll() ([]*model.Order, error)
	Delete(id string) error
}
