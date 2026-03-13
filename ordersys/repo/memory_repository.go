package repo

import (
	"fmt"
	"go_test/ordersys/domain/model/agg"
	"go_test/ordersys/domain/model/entity"
	"sync"
)

// InMemoryOrderRepository 内存订单仓储实现
// 基础设施层：实现领域层定义的仓储接口
type InMemoryOrderRepository struct {
	orders map[string]*agg.Order
	mu     sync.RWMutex
}

// NewInMemoryOrderRepository 创建内存订单仓储
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*agg.Order),
	}
}

// Save 保存订单
func (r *InMemoryOrderRepository) Save(order *agg.Order) error {
	if order == nil {
		return fmt.Errorf("订单不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID()] = order
	return nil
}

// FindByID 根据ID查找订单
func (r *InMemoryOrderRepository) FindByID(id string) (*agg.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("订单不存在: %s", id)
	}

	return order, nil
}

// FindByUserID 根据用户ID查找订单列表
func (r *InMemoryOrderRepository) FindByUserID(userID string) ([]*agg.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var orders []*agg.Order
	for _, order := range r.orders {
		if order.UserID() == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// FindAll 查找所有订单
func (r *InMemoryOrderRepository) FindAll() ([]*agg.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*agg.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// Delete 删除订单
func (r *InMemoryOrderRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return fmt.Errorf("订单不存在: %s", id)
	}

	delete(r.orders, id)
	return nil
}

// InMemoryProductRepository 内存产品仓储实现
type InMemoryProductRepository struct {
	products map[string]*entity.Product
	mu       sync.RWMutex
}

// NewInMemoryProductRepository 创建内存产品仓储
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*entity.Product),
	}
}

// Save 保存产品
func (r *InMemoryProductRepository) Save(product *entity.Product) error {
	if product == nil {
		return fmt.Errorf("产品不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID()] = product
	return nil
}

// FindByID 根据ID查找产品
func (r *InMemoryProductRepository) FindByID(id string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("产品不存在: %s", id)
	}

	return product, nil
}

// FindAll 查找所有产品
func (r *InMemoryProductRepository) FindAll() ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entity.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}

// Delete 删除产品
func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return fmt.Errorf("产品不存在: %s", id)
	}

	delete(r.products, id)
	return nil
}
