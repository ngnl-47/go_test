package service

// 包级别领域服务（无状态），init 时初始化
var (
	DefaultPricingService   *OrderPricingService
	DefaultInventoryService *ProductInventoryService
)

func init() {
	DefaultPricingService = NewOrderPricingService(0.1)
	DefaultInventoryService = NewProductInventoryService()
}
