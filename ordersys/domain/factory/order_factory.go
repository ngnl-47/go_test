package factory

import (
	"fmt"

	"go_test/ordersys/domain/model/agg"
	"go_test/ordersys/domain/model/event"
	entity2 "go_test/ordersys/domain/model/entity"
	valueobject2 "go_test/ordersys/domain/model/vo"
)

// ProductItem 订单项入参
type ProductItem struct {
	Product  *entity2.Product
	Quantity int
}

// CreateOrder 静态工厂方法：创建订单（使用包变量 idGenerator）
func CreateOrder(
	userID string,
	shippingAddress *valueobject2.Address,
	productItems []ProductItem,
) (*agg.Order, error) {
	if len(productItems) == 0 {
		return nil, fmt.Errorf("订单项不能为空")
	}
	if idGenerator == nil {
		return nil, fmt.Errorf("工厂未设置 ID 生成器，请调用 factory.SetIDGenerator")
	}

	orderID := idGenerator.Generate()
	order, err := agg.NewOrder(orderID, userID, shippingAddress)
	if err != nil {
		return nil, err
	}

	for _, item := range productItems {
		itemID := idGenerator.Generate()
		orderItem, err := entity2.NewOrderItem(itemID, item.Product, item.Quantity)
		if err != nil {
			return nil, err
		}
		if err := order.AddItem(orderItem); err != nil {
			return nil, err
		}
	}

	var totalAmount float64
	if order.TotalAmount() != nil {
		totalAmount = order.TotalAmount().Amount()
	}
	order.AddDomainEvent(event.NewOrderCreatedEvent(orderID, userID, totalAmount))
	return order, nil
}
