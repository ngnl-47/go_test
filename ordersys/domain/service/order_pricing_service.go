package service

import (
	"errors"

	"go_test/ordersys/domain/model/agg"
)

// OrderPricingService 订单定价领域服务
type OrderPricingService struct {
	discountRate float64
}

// NewOrderPricingService 创建订单定价服务
func NewOrderPricingService(discountRate float64) *OrderPricingService {
	return &OrderPricingService{discountRate: discountRate}
}

// CalculateDiscountedTotal 计算折扣后的订单总价
func (s *OrderPricingService) CalculateDiscountedTotal(order *agg.Order) (float64, error) {
	if order == nil {
		return 0, errors.New("订单不能为空")
	}
	totalAmount := order.TotalAmount()
	if totalAmount == nil {
		return 0, nil
	}
	return totalAmount.Amount() * (1 - s.discountRate), nil
}
