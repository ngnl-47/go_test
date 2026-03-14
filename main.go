package main

import (
	"fmt"
	"log"

	"go_test/ordersys/application/dto"
	appservice "go_test/ordersys/application/service"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  DDD电商订单系统示例")
	fmt.Println("========================================")
	fmt.Println()

	orderAppService := appservice.OrderAppService

	fmt.Println("步骤1: 创建产品")
	fmt.Println("----------------------------------------")
	product1, err := orderAppService.CreateProduct("iPhone 15 Pro", "最新款苹果手机", 7999.00, "CNY", 100)
	if err != nil {
		log.Fatalf("创建产品失败: %v", err)
	}
	fmt.Printf("✓ 产品创建成功: %s (ID: %s, 价格: %.2f %s, 库存: %d)\n",
		product1.Name, product1.ProductID, product1.Price, product1.Currency, product1.Stock)

	product2, err := orderAppService.CreateProduct("MacBook Pro", "专业笔记本电脑", 14999.00, "CNY", 50)
	if err != nil {
		log.Fatalf("创建产品失败: %v", err)
	}
	fmt.Printf("✓ 产品创建成功: %s (ID: %s, 价格: %.2f %s, 库存: %d)\n",
		product2.Name, product2.ProductID, product2.Price, product2.Currency, product2.Stock)
	fmt.Println()

	fmt.Println("步骤2: 创建订单")
	fmt.Println("----------------------------------------")
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

	order, err := orderAppService.CreateOrder(createOrderReq)
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	fmt.Printf("✓ 订单创建成功\n")
	fmt.Printf("  订单ID: %s\n", order.OrderID)
	fmt.Printf("  用户ID: %s\n", order.UserID)
	fmt.Printf("  订单状态: %s\n", order.Status)
	fmt.Printf("  总金额: %.2f %s\n", order.TotalAmount, order.Currency)
	fmt.Printf("  收货地址: %s\n", order.ShippingAddress)
	fmt.Println("  订单项:")
	for _, item := range order.Items {
		fmt.Printf("    - %s x%d, 单价: %.2f, 小计: %.2f\n",
			item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice)
	}
	fmt.Println()

	fmt.Println("步骤3: 支付订单")
	fmt.Println("----------------------------------------")
	err = orderAppService.PayOrder(order.OrderID, "支付宝")
	if err != nil {
		log.Fatalf("支付订单失败: %v", err)
	}
	fmt.Printf("✓ 订单支付成功\n")

	paidOrder, err := orderAppService.GetOrder(order.OrderID)
	if err != nil {
		log.Fatalf("查询订单失败: %v", err)
	}
	fmt.Printf("  订单状态: %s\n", paidOrder.Status)
	fmt.Println()

	fmt.Println("步骤4: 发货")
	fmt.Println("----------------------------------------")
	err = orderAppService.ShipOrder(order.OrderID, "SF1234567890")
	if err != nil {
		log.Fatalf("发货失败: %v", err)
	}
	fmt.Printf("✓ 订单发货成功\n")

	shippedOrder, err := orderAppService.GetOrder(order.OrderID)
	if err != nil {
		log.Fatalf("查询订单失败: %v", err)
	}
	fmt.Printf("  订单状态: %s\n", shippedOrder.Status)
	fmt.Println()

	fmt.Println("步骤5: 查询用户订单列表")
	fmt.Println("----------------------------------------")
	userOrders, err := orderAppService.GetUserOrders("USER-001")
	if err != nil {
		log.Fatalf("查询用户订单失败: %v", err)
	}
	fmt.Printf("✓ 用户 USER-001 共有 %d 个订单\n", len(userOrders))
	for _, o := range userOrders {
		fmt.Printf("  - 订单ID: %s, 状态: %s, 金额: %.2f\n", o.OrderID, o.Status, o.TotalAmount)
	}
	fmt.Println()

	fmt.Println("步骤6: 创建并取消订单")
	fmt.Println("----------------------------------------")
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

	cancelOrder, err := orderAppService.CreateOrder(cancelOrderReq)
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	fmt.Printf("✓ 订单创建成功: %s\n", cancelOrder.OrderID)

	err = orderAppService.CancelOrder(cancelOrder.OrderID, "用户主动取消")
	if err != nil {
		log.Fatalf("取消订单失败: %v", err)
	}
	fmt.Printf("✓ 订单取消成功\n")

	cancelledOrder, err := orderAppService.GetOrder(cancelOrder.OrderID)
	if err != nil {
		log.Fatalf("查询订单失败: %v", err)
	}
	fmt.Printf("  订单状态: %s\n", cancelledOrder.Status)
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("  示例执行完成")
	fmt.Println("========================================")
}
