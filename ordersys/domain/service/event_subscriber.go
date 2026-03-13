package service

import (
	"go_test/common"
	"go_test/ordersys/domain/model/event"
)

func init() {
	registerOrderEventHandlers()
}

// registerOrderEventHandlers 向全局事件总线注册订单领域事件回调
func registerOrderEventHandlers() {
	created := &orderCreatedHandler{}
	paid := &orderPaidHandler{}
	shipped := &orderShippedHandler{}
	cancelled := &orderCancelledHandler{}
	wrap := func(h func(event.DomainEvent) error) func(interface{}) {
		return func(evt interface{}) { _ = h(evt.(event.DomainEvent)) }
	}
	_ = common.GlobalEvent.Subscribe("OrderCreated", wrap(created.Handle))
	_ = common.GlobalEvent.Subscribe("OrderPaid", wrap(paid.Handle))
	_ = common.GlobalEvent.Subscribe("OrderShipped", wrap(shipped.Handle))
	_ = common.GlobalEvent.Subscribe("OrderCancelled", wrap(cancelled.Handle))
}
