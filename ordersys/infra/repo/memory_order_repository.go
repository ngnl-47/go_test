package repo

import (
	"fmt"
	"sync"

	"go_test/ordersys/domain/model"
)

// InMemoryOrderRepository 内存订单仓储实现
type InMemoryOrderRepository struct {
	orders map[string]*model.Order
	mu     sync.RWMutex
}

// NewInMemoryOrderRepository 创建内存订单仓储
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

// Save 保存订单
func (r *InMemoryOrderRepository) Save(order *model.Order) error {
	if order == nil {
		return fmt.Errorf("订单不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID()] = order
	return nil
}

// FindByID 根据ID查找订单
func (r *InMemoryOrderRepository) FindByID(id string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("订单不存在: %s", id)
	}

	return order, nil
}

// FindByUserID 根据用户ID查找订单列表
func (r *InMemoryOrderRepository) FindByUserID(userID string) ([]*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var orders []*model.Order
	for _, order := range r.orders {
		if order.UserID() == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// FindAll 查找所有订单
func (r *InMemoryOrderRepository) FindAll() ([]*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*model.Order, 0, len(r.orders))
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
