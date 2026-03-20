package repo

import (
	"fmt"
	"sync"

	"go_test/ordersys/domain/model"
)

// InMemoryProductRepository 内存产品仓储实现
type InMemoryProductRepository struct {
	products map[string]*model.Product
	mu       sync.RWMutex
}

// NewInMemoryProductRepository 创建内存产品仓储
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*model.Product),
	}
}

// Save 保存产品
func (r *InMemoryProductRepository) Save(product *model.Product) error {
	if product == nil {
		return fmt.Errorf("产品不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID()] = product
	return nil
}

// FindByID 根据ID查找产品
func (r *InMemoryProductRepository) FindByID(id string) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("产品不存在: %s", id)
	}

	return product, nil
}

// FindAll 查找所有产品
func (r *InMemoryProductRepository) FindAll() ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*model.Product, 0, len(r.products))
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
