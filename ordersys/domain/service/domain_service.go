package service

import (
	"errors"
	"go_test/ordersys/domain/model/aggregate"
	"go_test/ordersys/domain/model/entity"
)

// OrderPricingService 订单定价服务
// 领域服务特性：
// 1. 不自然属于任何实体或值对象的业务逻辑
// 2. 无状态
// 3. 以领域概念命名
type OrderPricingService struct {
	discountRate float64
}

// NewOrderPricingService 创建订单定价服务
func NewOrderPricingService(discountRate float64) *OrderPricingService {
	return &OrderPricingService{
		discountRate: discountRate,
	}
}

// CalculateDiscountedTotal 计算折扣后的订单总价
// 领域服务方法：跨多个实体的业务逻辑
func (s *OrderPricingService) CalculateDiscountedTotal(order *aggregate.Order) (float64, error) {
	if order == nil {
		return 0, errors.New("订单不能为空")
	}

	totalAmount := order.TotalAmount()
	if totalAmount == nil {
		return 0, nil
	}

	discountedTotal := totalAmount.Amount() * (1 - s.discountRate)
	return discountedTotal, nil
}

// ProductInventoryService 产品库存服务
// 领域服务：管理产品库存相关的业务逻辑
type ProductInventoryService struct{}

// NewProductInventoryService 创建产品库存服务
func NewProductInventoryService() *ProductInventoryService {
	return &ProductInventoryService{}
}

// CheckInventory 检查库存是否充足
// 领域服务方法：涉及多个产品的库存检查
func (s *ProductInventoryService) CheckInventory(products []*entity.Product, quantities map[string]int) error {
	for _, product := range products {
		requiredQty, exists := quantities[product.ID()]
		if !exists {
			continue
		}

		if product.Stock() < requiredQty {
			return errors.New("库存不足")
		}
	}
	return nil
}

// ReserveInventory 预留库存
// 领域服务方法：批量预留库存
func (s *ProductInventoryService) ReserveInventory(products []*entity.Product, quantities map[string]int) error {
	for _, product := range products {
		requiredQty, exists := quantities[product.ID()]
		if !exists {
			continue
		}

		if err := product.DecreaseStock(requiredQty); err != nil {
			return err
		}
	}
	return nil
}

// ReleaseInventory 释放库存
// 领域服务方法：批量释放库存（如订单取消时）
func (s *ProductInventoryService) ReleaseInventory(products []*entity.Product, quantities map[string]int) error {
	for _, product := range products {
		releaseQty, exists := quantities[product.ID()]
		if !exists {
			continue
		}

		if err := product.IncreaseStock(releaseQty); err != nil {
			return err
		}
	}
	return nil
}
