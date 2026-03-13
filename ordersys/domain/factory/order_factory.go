package factory

import (
	"fmt"
	"go_test/ordersys/domain/model/aggregate"
	entity2 "go_test/ordersys/domain/model/entity"
	valueobject2 "go_test/ordersys/domain/model/valueobject"
)

// OrderFactory 订单工厂
// 工厂模式特性：
// 1. 负责复杂对象的创建
// 2. 封装创建逻辑
// 3. 确保创建的对象处于有效状态
type OrderFactory struct {
	idGenerator IDGenerator
}

// IDGenerator ID生成器接口
// 依赖倒置原则：领域层定义接口，基础设施层实现
type IDGenerator interface {
	Generate() string
}

// NewOrderFactory 创建订单工厂
func NewOrderFactory(idGenerator IDGenerator) *OrderFactory {
	return &OrderFactory{
		idGenerator: idGenerator,
	}
}

// CreateOrder 创建订单
// 工厂方法：封装订单创建的复杂逻辑
func (f *OrderFactory) CreateOrder(
	userID string,
	shippingAddress *valueobject2.Address,
	productItems []ProductItem,
) (*aggregate.Order, error) {
	if len(productItems) == 0 {
		return nil, fmt.Errorf("订单项不能为空")
	}

	orderID := f.idGenerator.Generate()
	order, err := aggregate.NewOrder(orderID, userID, shippingAddress)
	if err != nil {
		return nil, err
	}

	for _, item := range productItems {
		itemID := f.idGenerator.Generate()
		orderItem, err := entity2.NewOrderItem(itemID, item.Product, item.Quantity)
		if err != nil {
			return nil, err
		}

		if err := order.AddItem(orderItem); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// ProductItem 产品项
// 工厂方法参数
type ProductItem struct {
	Product  *entity2.Product
	Quantity int
}

// ProductFactory 产品工厂
type ProductFactory struct {
	idGenerator IDGenerator
}

// NewProductFactory 创建产品工厂
func NewProductFactory(idGenerator IDGenerator) *ProductFactory {
	return &ProductFactory{
		idGenerator: idGenerator,
	}
}

// CreateProduct 创建产品
func (f *ProductFactory) CreateProduct(
	name string,
	description string,
	price *valueobject2.Money,
	stock int,
) (*entity2.Product, error) {
	productID := f.idGenerator.Generate()
	return entity2.NewProduct(productID, name, description, price, stock)
}
