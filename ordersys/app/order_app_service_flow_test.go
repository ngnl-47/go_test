package app_test

import (
	"testing"

	"go_test/ordersys/app/dto"
	appservice "go_test/ordersys/app/service"
	domainservice "go_test/ordersys/domain/service"
	infrarepo "go_test/ordersys/infra/repo"
)

func newIsolatedOrderAppService() *appservice.OrderAppService {
	return appservice.NewOrderAppService(
		infrarepo.NewInMemoryOrderRepository(),
		infrarepo.NewInMemoryProductRepository(),
		domainservice.DefaultPricingService,
		domainservice.DefaultInventoryService,
	)
}

// TestOrderAppService_EcommerceFlow 覆盖原 main 演示：创建产品、下单、支付、发货、用户订单列表、取消订单。
func TestOrderAppService_EcommerceFlow(t *testing.T) {
	svc := newIsolatedOrderAppService()

	product1, err := svc.CreateProduct("iPhone 15 Pro", "最新款苹果手机", 7999.00, "CNY", 100)
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}
	if product1.Name != "iPhone 15 Pro" || product1.Stock != 100 || product1.Currency != "CNY" {
		t.Fatalf("product1: %+v", product1)
	}

	product2, err := svc.CreateProduct("MacBook Pro", "专业笔记本电脑", 14999.00, "CNY", 50)
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}
	if product2.Name != "MacBook Pro" {
		t.Fatalf("product2: %+v", product2)
	}

	createOrderReq := &dto.CreateOrderRequest{
		UserID:   "USER-001",
		Province: "广东省",
		City:     "深圳市",
		District: "南山区",
		Street:   "科技园路100号",
		ZipCode:  "518000",
		ProductItems: []dto.ProductItemDTO{
			{ProductID: product1.ProductID, Quantity: 2},
			{ProductID: product2.ProductID, Quantity: 1},
		},
	}

	order, err := svc.CreateOrder(createOrderReq)
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if order.UserID != "USER-001" || order.Status != "待支付" {
		t.Fatalf("order after create: %+v", order)
	}
	wantTotal := 7999.0*2 + 14999.0
	if order.TotalAmount != wantTotal || order.Currency != "CNY" {
		t.Fatalf("order total: got %.2f %s want %.2f CNY", order.TotalAmount, order.Currency, wantTotal)
	}
	if len(order.Items) != 2 {
		t.Fatalf("order items: %d", len(order.Items))
	}

	if err := svc.PayOrder(order.OrderID, "支付宝"); err != nil {
		t.Fatalf("PayOrder: %v", err)
	}
	paid, err := svc.GetOrder(order.OrderID)
	if err != nil || paid.Status != "已支付" {
		t.Fatalf("after pay: err=%v status=%q", err, paid.Status)
	}

	if err := svc.ShipOrder(order.OrderID, "SF1234567890"); err != nil {
		t.Fatalf("ShipOrder: %v", err)
	}
	shipped, err := svc.GetOrder(order.OrderID)
	if err != nil || shipped.Status != "已发货" {
		t.Fatalf("after ship: err=%v status=%q", err, shipped.Status)
	}

	userOrders, err := svc.GetUserOrders("USER-001")
	if err != nil {
		t.Fatalf("GetUserOrders: %v", err)
	}
	if len(userOrders) != 1 || userOrders[0].OrderID != order.OrderID {
		t.Fatalf("GetUserOrders: %+v", userOrders)
	}

	cancelOrderReq := &dto.CreateOrderRequest{
		UserID:   "USER-002",
		Province: "北京市",
		City:     "北京市",
		District: "朝阳区",
		Street:   "建国路88号",
		ZipCode:  "100000",
		ProductItems: []dto.ProductItemDTO{
			{ProductID: product1.ProductID, Quantity: 1},
		},
	}
	cancelOrder, err := svc.CreateOrder(cancelOrderReq)
	if err != nil {
		t.Fatalf("CreateOrder cancel path: %v", err)
	}
	if err := svc.CancelOrder(cancelOrder.OrderID, "用户主动取消"); err != nil {
		t.Fatalf("CancelOrder: %v", err)
	}
	cancelled, err := svc.GetOrder(cancelOrder.OrderID)
	if err != nil || cancelled.Status != "已取消" {
		t.Fatalf("after cancel: err=%v status=%q", err, cancelled.Status)
	}
}
