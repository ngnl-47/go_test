package service

import (
	"fmt"

	"go_test/common"
	"go_test/ordersys/app/dto"
	"go_test/ordersys/domain/facade"
	"go_test/ordersys/domain/factory"
	"go_test/ordersys/domain/model"
	domainservice "go_test/ordersys/domain/service"
)

// OrderAppService 订单应用服务
// 应用服务特性：协调领域对象完成业务用例，不承担核心业务规则。
type OrderAppService struct {
	orderRepo        facade.OrderRepository
	productRepo      facade.ProductRepository
	pricingService   *domainservice.OrderPricingService
	inventoryService *domainservice.ProductInventoryService
}

// NewOrderAppService 创建订单应用服务
func NewOrderAppService(
	orderRepo facade.OrderRepository,
	productRepo facade.ProductRepository,
	pricingService *domainservice.OrderPricingService,
	inventoryService *domainservice.ProductInventoryService,
) *OrderAppService {
	return &OrderAppService{
		orderRepo:        orderRepo,
		productRepo:      productRepo,
		pricingService:   pricingService,
		inventoryService: inventoryService,
	}
}

// CreateOrder 创建订单
func (s *OrderAppService) CreateOrder(req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	shippingAddress, err := model.NewAddress(
		req.Province,
		req.City,
		req.District,
		req.Street,
		req.ZipCode,
	)
	if err != nil {
		return nil, err
	}

	productItems := make([]factory.ProductItem, 0, len(req.ProductItems))
	for _, item := range req.ProductItems {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("产品不存在: %s", item.ProductID)
		}
		productItems = append(productItems, factory.ProductItem{
			Product:  product,
			Quantity: item.Quantity,
		})
	}

	order, err := factory.CreateOrder(req.UserID, shippingAddress, productItems)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return nil, err
	}

	if err := s.publishDomainEvents(order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// PayOrder 支付订单
func (s *OrderAppService) PayOrder(orderID, paymentMethod string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if err := order.Pay(paymentMethod); err != nil {
		return err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return err
	}

	return s.publishDomainEvents(order)
}

// ShipOrder 发货
func (s *OrderAppService) ShipOrder(orderID, trackingNumber string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if err := order.Ship(trackingNumber); err != nil {
		return err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return err
	}

	return s.publishDomainEvents(order)
}

// CancelOrder 取消订单
func (s *OrderAppService) CancelOrder(orderID, reason string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(reason); err != nil {
		return err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return err
	}

	return s.publishDomainEvents(order)
}

// GetOrder 查询订单
func (s *OrderAppService) GetOrder(orderID string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetUserOrders 查询用户订单列表
func (s *OrderAppService) GetUserOrders(userID string) ([]*dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, s.toOrderResponse(order))
	}

	return responses, nil
}

// CreateProduct 创建产品
func (s *OrderAppService) CreateProduct(
	name string,
	description string,
	price float64,
	currency string,
	stock int,
) (*dto.ProductResponse, error) {
	money, err := model.NewMoney(price, currency)
	if err != nil {
		return nil, err
	}

	product, err := factory.CreateProduct(name, description, money, stock)
	if err != nil {
		return nil, err
	}

	if err := s.productRepo.Save(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// publishDomainEvents 发布领域事件到全局总线 common.GlobalEvent
func (s *OrderAppService) publishDomainEvents(order *model.Order) error {
	for _, evt := range order.DomainEvents() {
		common.GlobalEvent.Publish(evt.EventType(), evt)
	}
	order.ClearDomainEvents()
	return nil
}

// toOrderResponse 将订单聚合根转换为DTO
func (s *OrderAppService) toOrderResponse(order *model.Order) *dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items()))
	for _, item := range order.Items() {
		itemTotal, _ := item.CalculateTotalPrice()
		items = append(items, dto.OrderItemResponse{
			ItemID:      item.ID(),
			ProductID:   item.Product().ID(),
			ProductName: item.Product().Name(),
			Quantity:    item.Quantity(),
			UnitPrice:   item.UnitPrice().Amount(),
			TotalPrice:  itemTotal.Amount(),
		})
	}

	var totalAmount float64
	var currency string
	if order.TotalAmount() != nil {
		totalAmount = order.TotalAmount().Amount()
		currency = order.TotalAmount().Currency()
	}

	return &dto.OrderResponse{
		OrderID:         order.ID(),
		UserID:          order.UserID(),
		Status:          order.Status().String(),
		TotalAmount:     totalAmount,
		Currency:        currency,
		ShippingAddress: order.ShippingAddress().String(),
		Items:           items,
	}
}

// toProductResponse 将产品实体转换为DTO
func (s *OrderAppService) toProductResponse(product *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ProductID:   product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       product.Price().Amount(),
		Currency:    product.Price().Currency(),
		Stock:       product.Stock(),
	}
}
