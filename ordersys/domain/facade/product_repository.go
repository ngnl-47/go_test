package facade

import "go_test/ordersys/domain/model"

// ProductRepository 产品仓储接口
type ProductRepository interface {
	Save(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	FindAll() ([]*model.Product, error)
	Delete(id string) error
}
