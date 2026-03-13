package valueobject

import (
	"errors"
	"fmt"
)

// Money 金额值对象
// 值对象特性：不可变、通过值比较相等性、无唯一标识
type Money struct {
	amount   float64
	currency string
}

// NewMoney 创建金额值对象
// 工厂方法确保金额的有效性
func NewMoney(amount float64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, errors.New("金额不能为负数")
	}
	if currency == "" {
		return nil, errors.New("货币单位不能为空")
	}
	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount 获取金额
func (m *Money) Amount() float64 {
	return m.amount
}

// Currency 获取货币单位
func (m *Money) Currency() string {
	return m.currency
}

// Add 金额相加，返回新的金额值对象
// 值对象不可变，操作返回新实例
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("货币单位不匹配: %s vs %s", m.currency, other.currency)
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract 金额相减，返回新的金额值对象
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("货币单位不匹配: %s vs %s", m.currency, other.currency)
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

// Multiply 金额乘以倍数，返回新的金额值对象
func (m *Money) Multiply(multiplier float64) (*Money, error) {
	return NewMoney(m.amount*multiplier, m.currency)
}

// Equals 判断两个金额是否相等
// 值对象通过属性值比较相等性
func (m *Money) Equals(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount == other.amount && m.currency == other.currency
}

// String 返回金额的字符串表示
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}
