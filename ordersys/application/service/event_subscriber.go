package service

import (
	"fmt"

	"go_test/common"
	domainevent "go_test/ordersys/domain/model/event"
)

func init() {
	registerOrderEventHandlers()
}

func wrap(fn func(domainevent.DomainEvent) error) func(interface{}) {
	return func(evt interface{}) { _ = fn(evt.(domainevent.DomainEvent)) }
}

// registerOrderEventHandlers 向全局事件总线注册订单领域事件回调（应用层消费）
func registerOrderEventHandlers() {
	_ = common.GlobalEvent.Subscribe("OrderCreated", wrap(onOrderCreated))
	_ = common.GlobalEvent.Subscribe("OrderPaid", wrap(onOrderPaid))
	_ = common.GlobalEvent.Subscribe("OrderShipped", wrap(onOrderShipped))
	_ = common.GlobalEvent.Subscribe("OrderCancelled", wrap(onOrderCancelled))
}

// 订单领域事件回调（应用层：发通知等侧效应）
func onOrderCreated(evt domainevent.DomainEvent) error {
	e := evt.(*domainevent.OrderCreatedEvent)
	fmt.Printf("📧 发送通知：订单 %s 已创建，用户ID: %s\n", e.OrderID, e.UserID)
	return nil
}

func onOrderPaid(evt domainevent.DomainEvent) error {
	e := evt.(*domainevent.OrderPaidEvent)
	fmt.Printf("💳 支付成功：订单 %s，支付金额: %.2f，支付方式: %s\n",
		e.OrderID, e.PaidAmount, e.PaymentMethod)
	return nil
}

func onOrderShipped(evt domainevent.DomainEvent) error {
	e := evt.(*domainevent.OrderShippedEvent)
	fmt.Printf("🚚 订单发货：订单 %s，物流单号: %s\n", e.OrderID, e.TrackingNumber)
	return nil
}

func onOrderCancelled(evt domainevent.DomainEvent) error {
	e := evt.(*domainevent.OrderCancelledEvent)
	fmt.Printf("❌ 订单取消：订单 %s，原因: %s\n", e.OrderID, e.Reason)
	return nil
}
