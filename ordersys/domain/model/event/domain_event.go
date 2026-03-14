package event

import "time"

// DomainEvent 领域事件接口
// 所有领域事件都需要实现此接口
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// BaseEvent 基础事件结构
// 提供事件的公共属性和方法
type BaseEvent struct {
	occurredAt time.Time
}

// OccurredAt 事件发生时间
func (e *BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// NewBaseEvent 创建基础事件
func NewBaseEvent() BaseEvent {
	return BaseEvent{
		occurredAt: time.Now(),
	}
}

// OrderCreatedEvent 订单创建事件
// 领域事件：表示领域中发生的重要事情
type OrderCreatedEvent struct {
	BaseEvent
	OrderID     string
	UserID      string
	TotalAmount float64
}

// EventType 返回事件类型
func (e *OrderCreatedEvent) EventType() string {
	return "OrderCreated"
}

// NewOrderCreatedEvent 创建订单创建事件（领域层 event 包创建，agg/factory 均可调用，无循环依赖）
func NewOrderCreatedEvent(orderID, userID string, totalAmount float64) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent:   NewBaseEvent(),
		OrderID:     orderID,
		UserID:      userID,
		TotalAmount: totalAmount,
	}
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	BaseEvent
	OrderID       string
	PaidAmount    float64
	PaymentMethod string
}

// EventType 返回事件类型
func (e *OrderPaidEvent) EventType() string {
	return "OrderPaid"
}

func NewOrderPaidEvent(orderID string, paidAmount float64, paymentMethod string) *OrderPaidEvent {
	return &OrderPaidEvent{
		BaseEvent:     NewBaseEvent(),
		OrderID:       orderID,
		PaidAmount:    paidAmount,
		PaymentMethod: paymentMethod,
	}
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	BaseEvent
	OrderID        string
	TrackingNumber string
}

// EventType 返回事件类型
func (e *OrderShippedEvent) EventType() string {
	return "OrderShipped"
}

func NewOrderShippedEvent(orderID, trackingNumber string) *OrderShippedEvent {
	return &OrderShippedEvent{
		BaseEvent:      NewBaseEvent(),
		OrderID:        orderID,
		TrackingNumber: trackingNumber,
	}
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	BaseEvent
	OrderID string
	Reason  string
}

// EventType 返回事件类型
func (e *OrderCancelledEvent) EventType() string {
	return "OrderCancelled"
}

func NewOrderCancelledEvent(orderID, reason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		BaseEvent: NewBaseEvent(),
		OrderID:   orderID,
		Reason:    reason,
	}
}

// ProductStockChangedEvent 产品库存变更事件
type ProductStockChangedEvent struct {
	BaseEvent
	ProductID    string
	OldStock     int
	NewStock     int
	ChangeReason string
}

// EventType 返回事件类型
func (e *ProductStockChangedEvent) EventType() string {
	return "ProductStockChanged"
}

func NewProductStockChangedEvent(productID string, oldStock, newStock int, changeReason string) *ProductStockChangedEvent {
	return &ProductStockChangedEvent{
		BaseEvent:    NewBaseEvent(),
		ProductID:    productID,
		OldStock:     oldStock,
		NewStock:     newStock,
		ChangeReason: changeReason,
	}
}
