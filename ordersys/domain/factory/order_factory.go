package factory

import (
	"fmt"

	"go_test/ordersys/domain/model"
)

// ProductItem 订单项入参
type ProductItem struct {
	Product  *model.Product
	Quantity int
}

// CreateOrder 静态工厂方法：创建订单（使用包变量 idGenerator）
func CreateOrder(
	userID string,
	shippingAddress *model.Address,
	productItems []ProductItem,
) (*model.Order, error) {
	if len(productItems) == 0 {
		return nil, fmt.Errorf("订单项不能为空")
	}
	if idGenerator == nil {
		return nil, fmt.Errorf("工厂未设置 ID 生成器，请调用 factory.SetIDGenerator")
	}

	orderID := idGenerator.Generate()
	order, err := model.NewOrder(orderID, userID, shippingAddress)
	if err != nil {
		return nil, err
	}

	for _, item := range productItems {
		itemID := idGenerator.Generate()
		orderItem, err := model.NewOrderItem(itemID, item.Product, item.Quantity)
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
	order.AddDomainEvent(model.NewOrderCreatedEvent(orderID, userID, totalAmount))
	return order, nil
}
