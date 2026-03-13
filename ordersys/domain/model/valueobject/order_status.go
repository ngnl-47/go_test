package valueobject

import "errors"

// OrderStatus 订单状态值对象
// 使用类型别名和常量表示有限的状态集合
type OrderStatus int

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusPaid
	OrderStatusShipped
	OrderStatusDelivered
	OrderStatusCancelled
)

// String 返回订单状态的字符串表示
func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusDelivered:
		return "已送达"
	case OrderStatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}

// CanTransitionTo 判断是否可以转换到目标状态
// 封装状态转换的业务规则
func (s OrderStatus) CanTransitionTo(target OrderStatus) error {
	switch s {
	case OrderStatusPending:
		if target == OrderStatusPaid || target == OrderStatusCancelled {
			return nil
		}
	case OrderStatusPaid:
		if target == OrderStatusShipped || target == OrderStatusCancelled {
			return nil
		}
	case OrderStatusShipped:
		if target == OrderStatusDelivered {
			return nil
		}
	}
	return errors.New("无效的状态转换")
}
