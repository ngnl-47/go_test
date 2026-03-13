package service

import (
	"fmt"
	"go_test/common"
	"go_test/ordersys/application/dto"
	"go_test/ordersys/domain/factory"
	"go_test/ordersys/domain/model/aggregate"
	"go_test/ordersys/domain/model/entity"
	valueobject2 "go_test/ordersys/domain/model/valueobject"
	"go_test/ordersys/domain/repository"
	"go_test/ordersys/domain/service"
)

// OrderApplicationService 订单应用服务
// 应用服务特性：
// 1. 协调领域对象完成业务用例
// 2. 不包含业务逻辑，业务逻辑在领域层
// 3. 负责事务管理和权限控制
// 4. DTO与领域对象之间的转换
type OrderApplicationService struct {
	orderRepo        repository.OrderRepository
	productRepo      repository.ProductRepository
	orderFactory     *factory.OrderFactory
	productFactory   *factory.ProductFactory
	pricingService   *service.OrderPricingService
	inventoryService *service.ProductInventoryService
}

// NewOrderApplicationService 创建订单应用服务
func NewOrderApplicationService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	orderFactory *factory.OrderFactory,
	productFactory *factory.ProductFactory,
	pricingService *service.OrderPricingService,
	inventoryService *service.ProductInventoryService,
) *OrderApplicationService {
	return &OrderApplicationService{
		orderRepo:        orderRepo,
		productRepo:      productRepo,
		orderFactory:     orderFactory,
		productFactory:   productFactory,
		pricingService:   pricingService,
		inventoryService: inventoryService,
	}
}

// CreateOrder 创建订单
// 应用服务方法：协调领域对象完成业务用例
func (s *OrderApplicationService) CreateOrder(req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	shippingAddress, err := valueobject2.NewAddress(
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

	order, err := s.orderFactory.CreateOrder(req.UserID, shippingAddress, productItems)
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
func (s *OrderApplicationService) PayOrder(orderID, paymentMethod string) error {
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
func (s *OrderApplicationService) ShipOrder(orderID, trackingNumber string) error {
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
func (s *OrderApplicationService) CancelOrder(orderID, reason string) error {
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
func (s *OrderApplicationService) GetOrder(orderID string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetUserOrders 查询用户订单列表
func (s *OrderApplicationService) GetUserOrders(userID string) ([]*dto.OrderResponse, error) {
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
func (s *OrderApplicationService) CreateProduct(
	name string,
	description string,
	price float64,
	currency string,
	stock int,
) (*dto.ProductResponse, error) {
	money, err := valueobject2.NewMoney(price, currency)
	if err != nil {
		return nil, err
	}

	product, err := s.productFactory.CreateProduct(name, description, money, stock)
	if err != nil {
		return nil, err
	}

	if err := s.productRepo.Save(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// publishDomainEvents 发布领域事件到全局总线 common.GlobalEvent
func (s *OrderApplicationService) publishDomainEvents(order *aggregate.Order) error {
	for _, evt := range order.DomainEvents() {
		common.GlobalEvent.Publish(evt.EventType(), evt)
	}
	order.ClearDomainEvents()
	return nil
}

// toOrderResponse 将订单聚合根转换为DTO
func (s *OrderApplicationService) toOrderResponse(order *aggregate.Order) *dto.OrderResponse {
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
func (s *OrderApplicationService) toProductResponse(product *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ProductID:   product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       product.Price().Amount(),
		Currency:    product.Price().Currency(),
		Stock:       product.Stock(),
	}
}
