package dto

// CreateOrderRequest 创建订单请求DTO
// DTO特性：用于应用层与外部交互的数据传输
type CreateOrderRequest struct {
	UserID       string
	Province     string
	City         string
	District     string
	Street       string
	ZipCode      string
	ProductItems []ProductItemDTO
}

// ProductItemDTO 产品项DTO
type ProductItemDTO struct {
	ProductID string
	Quantity  int
}

// OrderResponse 订单响应DTO
type OrderResponse struct {
	OrderID         string
	UserID          string
	Status          string
	TotalAmount     float64
	Currency        string
	ShippingAddress string
	Items           []OrderItemResponse
}

// OrderItemResponse 订单项响应DTO
type OrderItemResponse struct {
	ItemID      string
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

// ProductResponse 产品响应DTO
type ProductResponse struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Currency    string
	Stock       int
}
