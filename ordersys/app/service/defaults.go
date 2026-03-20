package service

import (
	"go_test/ordersys/domain/factory"
	domainservice "go_test/ordersys/domain/service"
	"go_test/ordersys/infra"
)

// OrderAppSvc 包级别订单应用服务，init 时组装
var OrderAppSvc *OrderAppService

func init() {
	factory.SetIDGenerator(infra.DefaultIDGen)
	OrderAppSvc = NewOrderAppService(
		infra.DefaultOrderRepo,
		infra.DefaultProductRepo,
		domainservice.DefaultPricingService,
		domainservice.DefaultInventoryService,
	)
}
