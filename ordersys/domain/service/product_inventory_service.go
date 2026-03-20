package service

import (
	"errors"

	"go_test/ordersys/domain/model"
)

// ProductInventoryService 产品库存领域服务
type ProductInventoryService struct{}

// NewProductInventoryService 创建产品库存服务
func NewProductInventoryService() *ProductInventoryService {
	return &ProductInventoryService{}
}

// CheckInventory 检查库存是否充足
func (s *ProductInventoryService) CheckInventory(products []*model.Product, quantities map[string]int) error {
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
func (s *ProductInventoryService) ReserveInventory(products []*model.Product, quantities map[string]int) error {
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

// ReleaseInventory 释放库存（如订单取消时）
func (s *ProductInventoryService) ReleaseInventory(products []*model.Product, quantities map[string]int) error {
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
