package agg

import (
	"errors"
	"fmt"
	"go_test/ordersys/domain/model/entity"
	"go_test/ordersys/domain/model/event"
	valueobject2 "go_test/ordersys/domain/model/vo"
)

// Order 订单聚合根
// 聚合根特性：
// 1. 是聚合的入口点，外部只能通过聚合根访问聚合内部对象
// 2. 维护聚合内部的一致性
// 3. 负责发布领域事件
// 4. 控制聚合边界内的所有操作
type Order struct {
	id              string
	userID          string
	items           []*entity.OrderItem
	shippingAddress *valueobject2.Address
	status          valueobject2.OrderStatus
	totalAmount     *valueobject2.Money
	domainEvents    []event.DomainEvent
}

// NewOrder 创建订单聚合根
// 工厂方法确保订单的有效性
func NewOrder(id, userID string, shippingAddress *valueobject2.Address) (*Order, error) {
	if id == "" {
		return nil, errors.New("订单ID不能为空")
	}
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if shippingAddress == nil {
		return nil, errors.New("收货地址不能为空")
	}

	order := &Order{
		id:              id,
		userID:          userID,
		items:           make([]*entity.OrderItem, 0),
		shippingAddress: shippingAddress,
		status:          valueobject2.OrderStatusPending,
		totalAmount:     nil,
		domainEvents:    make([]event.DomainEvent, 0),
	}

	return order, nil
}

// AddDomainEvent 添加领域事件（供订单工厂在创建完成后挂 OrderCreated 等事件）
func (o *Order) AddDomainEvent(evt event.DomainEvent) {
	o.addDomainEvent(evt)
}

// ID 获取订单ID
func (o *Order) ID() string {
	return o.id
}

// UserID 获取用户ID
func (o *Order) UserID() string {
	return o.userID
}

// Items 获取订单项列表
// 聚合根控制对内部对象的访问
func (o *Order) Items() []*entity.OrderItem {
	return o.items
}

// ShippingAddress 获取收货地址
func (o *Order) ShippingAddress() *valueobject2.Address {
	return o.shippingAddress
}

// Status 获取订单状态
func (o *Order) Status() valueobject2.OrderStatus {
	return o.status
}

// TotalAmount 获取总金额
func (o *Order) TotalAmount() *valueobject2.Money {
	return o.totalAmount
}

// DomainEvents 获取领域事件列表
func (o *Order) DomainEvents() []event.DomainEvent {
	return o.domainEvents
}

// AddItem 添加订单项
// 聚合根维护内部一致性
func (o *Order) AddItem(item *entity.OrderItem) error {
	if item == nil {
		return errors.New("订单项不能为空")
	}

	for _, existingItem := range o.items {
		if existingItem.Product().ID() == item.Product().ID() {
			newQty := existingItem.Quantity() + item.Quantity()
			if err := existingItem.UpdateQuantity(newQty); err != nil {
				return err
			}
			return o.recalculateTotalAmount()
		}
	}

	o.items = append(o.items, item)
	return o.recalculateTotalAmount()
}

// RemoveItem 移除订单项
func (o *Order) RemoveItem(itemID string) error {
	for i, item := range o.items {
		if item.ID() == itemID {
			o.items = append(o.items[:i], o.items[i+1:]...)
			return o.recalculateTotalAmount()
		}
	}
	return fmt.Errorf("订单项不存在: %s", itemID)
}

// UpdateItemQuantity 更新订单项数量
func (o *Order) UpdateItemQuantity(itemID string, quantity int) error {
	for _, item := range o.items {
		if item.ID() == itemID {
			if err := item.UpdateQuantity(quantity); err != nil {
				return err
			}
			return o.recalculateTotalAmount()
		}
	}
	return fmt.Errorf("订单项不存在: %s", itemID)
}

// Pay 支付订单
// 充血模型：业务逻辑在领域对象内部
func (o *Order) Pay(paymentMethod string) error {
	if err := o.status.CanTransitionTo(valueobject2.OrderStatusPaid); err != nil {
		return err
	}

	if o.totalAmount == nil || o.totalAmount.Amount() == 0 {
		return errors.New("订单金额为0，无需支付")
	}

	o.status = valueobject2.OrderStatusPaid

	o.addDomainEvent(event.NewOrderPaidEvent(o.id, o.totalAmount.Amount(), paymentMethod))

	return nil
}

// Ship 发货
func (o *Order) Ship(trackingNumber string) error {
	if err := o.status.CanTransitionTo(valueobject2.OrderStatusShipped); err != nil {
		return err
	}

	o.status = valueobject2.OrderStatusShipped

	o.addDomainEvent(event.NewOrderShippedEvent(o.id, trackingNumber))

	return nil
}

// Deliver 确认送达
func (o *Order) Deliver() error {
	if err := o.status.CanTransitionTo(valueobject2.OrderStatusDelivered); err != nil {
		return err
	}

	o.status = valueobject2.OrderStatusDelivered
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel(reason string) error {
	if err := o.status.CanTransitionTo(valueobject2.OrderStatusCancelled); err != nil {
		return err
	}

	o.status = valueobject2.OrderStatusCancelled

	o.addDomainEvent(event.NewOrderCancelledEvent(o.id, reason))

	return nil
}

// UpdateShippingAddress 更新收货地址
func (o *Order) UpdateShippingAddress(address *valueobject2.Address) error {
	if address == nil {
		return errors.New("收货地址不能为空")
	}

	if o.status != valueobject2.OrderStatusPending {
		return errors.New("只有待支付状态的订单可以修改收货地址")
	}

	o.shippingAddress = address
	return nil
}

// recalculateTotalAmount 重新计算总金额
// 私有方法，聚合根内部使用
func (o *Order) recalculateTotalAmount() error {
	if len(o.items) == 0 {
		o.totalAmount = nil
		return nil
	}

	var total *valueobject2.Money

	for _, item := range o.items {
		itemTotal, err := item.CalculateTotalPrice()
		if err != nil {
			return err
		}

		if total == nil {
			total = itemTotal
		} else {
			total, err = total.Add(itemTotal)
			if err != nil {
				return err
			}
		}
	}

	o.totalAmount = total
	return nil
}

// addDomainEvent 添加领域事件
// 私有方法，聚合根内部使用
func (o *Order) addDomainEvent(evt event.DomainEvent) {
	o.domainEvents = append(o.domainEvents, evt)
}

// ClearDomainEvents 清除领域事件
// 在事件发布后调用
func (o *Order) ClearDomainEvents() {
	o.domainEvents = make([]event.DomainEvent, 0)
}

// Equals 判断两个订单是否相等
func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.id == other.id
}
