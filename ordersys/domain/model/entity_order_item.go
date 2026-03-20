package model

import "errors"

// OrderItem 订单项实体
// 实体特性：有唯一标识（在订单范围内）、可变、包含业务逻辑
type OrderItem struct {
	id        string
	product   *Product
	quantity  int
	unitPrice *Money
}

// NewOrderItem 创建订单项实体
func NewOrderItem(id string, product *Product, quantity int) (*OrderItem, error) {
	if id == "" {
		return nil, errors.New("订单项ID不能为空")
	}
	if product == nil {
		return nil, errors.New("产品不能为空")
	}
	if quantity <= 0 {
		return nil, errors.New("数量必须大于0")
	}
	return &OrderItem{
		id:        id,
		product:   product,
		quantity:  quantity,
		unitPrice: product.Price(),
	}, nil
}

// ID 获取订单项ID
func (oi *OrderItem) ID() string {
	return oi.id
}

// Product 获取产品
func (oi *OrderItem) Product() *Product {
	return oi.product
}

// Quantity 获取数量
func (oi *OrderItem) Quantity() int {
	return oi.quantity
}

// UnitPrice 获取单价
func (oi *OrderItem) UnitPrice() *Money {
	return oi.unitPrice
}

// UpdateQuantity 更新数量
func (oi *OrderItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("数量必须大于0")
	}
	oi.quantity = quantity
	return nil
}

// CalculateTotalPrice 计算订单项总价
// 充血模型：业务逻辑在领域对象内部
func (oi *OrderItem) CalculateTotalPrice() (*Money, error) {
	return oi.unitPrice.Multiply(float64(oi.quantity))
}

// Equals 判断两个订单项是否相等
func (oi *OrderItem) Equals(other *OrderItem) bool {
	if other == nil {
		return false
	}
	return oi.id == other.id
}
