package service

import (
	"fmt"
	"go_test/ordersys/domain/model/event"
)

// orderCreatedHandler 订单创建事件处理器
type orderCreatedHandler struct{}

func (h *orderCreatedHandler) Handle(evt event.DomainEvent) error {
	e := evt.(*event.OrderCreatedEvent)
	fmt.Printf("📧 发送通知：订单 %s 已创建，用户ID: %s\n", e.OrderID, e.UserID)
	return nil
}

// orderPaidHandler 订单支付事件处理器
type orderPaidHandler struct{}

func (h *orderPaidHandler) Handle(evt event.DomainEvent) error {
	e := evt.(*event.OrderPaidEvent)
	fmt.Printf("💳 支付成功：订单 %s，支付金额: %.2f，支付方式: %s\n",
		e.OrderID, e.PaidAmount, e.PaymentMethod)
	return nil
}

// orderShippedHandler 订单发货事件处理器
type orderShippedHandler struct{}

func (h *orderShippedHandler) Handle(evt event.DomainEvent) error {
	e := evt.(*event.OrderShippedEvent)
	fmt.Printf("🚚 订单发货：订单 %s，物流单号: %s\n", e.OrderID, e.TrackingNumber)
	return nil
}

// orderCancelledHandler 订单取消事件处理器
type orderCancelledHandler struct{}

func (h *orderCancelledHandler) Handle(evt event.DomainEvent) error {
	e := evt.(*event.OrderCancelledEvent)
	fmt.Printf("❌ 订单取消：订单 %s，原因: %s\n", e.OrderID, e.Reason)
	return nil
}
