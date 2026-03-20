package model

import (
	"errors"
	"fmt"
)

// Product 产品实体
// 实体特性：有唯一标识、可变、通过标识判断相等性
type Product struct {
	id          string
	name        string
	description string
	price       *Money
	stock       int
}

// NewProduct 创建产品实体
func NewProduct(id, name, description string, price *Money, stock int) (*Product, error) {
	if id == "" {
		return nil, errors.New("产品ID不能为空")
	}
	if name == "" {
		return nil, errors.New("产品名称不能为空")
	}
	if price == nil {
		return nil, errors.New("产品价格不能为空")
	}
	if stock < 0 {
		return nil, errors.New("库存不能为负数")
	}
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
	}, nil
}

// ID 获取产品ID
func (p *Product) ID() string {
	return p.id
}

// Name 获取产品名称
func (p *Product) Name() string {
	return p.name
}

// Description 获取产品描述
func (p *Product) Description() string {
	return p.description
}

// Price 获取产品价格
func (p *Product) Price() *Money {
	return p.price
}

// Stock 获取库存
func (p *Product) Stock() int {
	return p.stock
}

// DecreaseStock 减少库存
// 实体可变，通过方法修改内部状态
func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("减少数量必须大于0")
	}
	if p.stock < quantity {
		return fmt.Errorf("库存不足，当前库存: %d，需要: %d", p.stock, quantity)
	}
	p.stock -= quantity
	return nil
}

// IncreaseStock 增加库存
func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("增加数量必须大于0")
	}
	p.stock += quantity
	return nil
}

// Equals 判断两个产品是否相等
// 实体通过唯一标识判断相等性
func (p *Product) Equals(other *Product) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}
