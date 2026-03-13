package event

import (
	"go_test/common"
	domainevent "go_test/ordersys/domain/model/event"
)

func init() {
	registerOrderEventHandlers()
}

func wrap(h func(domainevent.DomainEvent) error) func(interface{}) {
	return func(evt interface{}) { _ = h(evt.(domainevent.DomainEvent)) }
}

// registerOrderEventHandlers 向全局事件总线注册订单领域事件回调（应用层消费）
func registerOrderEventHandlers() {
	_ = common.GlobalEvent.Subscribe("OrderCreated", wrap((&orderCreatedHandler{}).Handle))
	_ = common.GlobalEvent.Subscribe("OrderPaid", wrap((&orderPaidHandler{}).Handle))
	_ = common.GlobalEvent.Subscribe("OrderShipped", wrap((&orderShippedHandler{}).Handle))
	_ = common.GlobalEvent.Subscribe("OrderCancelled", wrap((&orderCancelledHandler{}).Handle))
}
