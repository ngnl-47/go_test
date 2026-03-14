package service

import (
	"go_test/ordersys/domain/factory"
	domainservice "go_test/ordersys/domain/service"
	infrarepo "go_test/ordersys/repo"
)

// 包级别订单应用服务，init 时组装；静态工厂使用 repo.DefaultIDGen 保证 ID 唯一
var OrderAppService *OrderApplicationService

func init() {
	factory.SetIDGenerator(infrarepo.DefaultIDGen)
	OrderAppService = NewOrderApplicationService(
		infrarepo.DefaultOrderRepo,
		infrarepo.DefaultProductRepo,
		domainservice.DefaultPricingService,
		domainservice.DefaultInventoryService,
	)
}
