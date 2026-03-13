# DDD电商订单系统架构设计文档

## 一、项目概述

本项目是一个基于领域驱动设计(DDD)架构的Go语言电商订单系统示例，完整展示了DDD领域模型的核心特性，包括实体、值对象、聚合根、领域服务、工厂、仓储、领域事件等核心概念。

## 二、架构分层

项目严格遵循DDD分层架构，分为以下几层：

### 2.1 领域层 (Domain Layer)
领域层是系统的核心，包含所有业务逻辑和领域模型。

**目录结构：**
```
ordersys/
├── domain/
│   ├── model/
│   │   ├── valueobject/      # 值对象
│   │   ├── entity/           # 实体
│   │   ├── aggregate/        # 聚合根
│   │   └── event/            # 领域事件
│   ├── service/              # 领域服务
│   ├── factory/              # 工厂
│   └── repository/           # 仓储接口
```

**核心组件：**

#### 1. 值对象 (Value Object)
- **Money**: 金额值对象，封装金额计算逻辑
- **Address**: 地址值对象，表示收货地址
- **OrderStatus**: 订单状态值对象，封装状态转换规则

**特性：**
- 不可变性：所有操作返回新实例
- 无唯一标识：通过属性值判断相等性
- 自验证：构造函数确保对象有效性

**代码示例：**
```go
// Money 金额值对象
type Money struct {
    amount   float64
    currency string
}

// Add 金额相加，返回新的金额值对象
func (m *Money) Add(other *Money) (*Money, error) {
    if m.currency != other.currency {
        return nil, fmt.Errorf("货币单位不匹配: %s vs %s", m.currency, other.currency)
    }
    return NewMoney(m.amount+other.amount, m.currency)
}
```

#### 2. 实体 (Entity)
- **Product**: 产品实体，包含产品信息和库存管理
- **OrderItem**: 订单项实体，关联产品和数量

**特性：**
- 有唯一标识：通过ID判断相等性
- 可变性：通过方法修改内部状态
- 充血模型：业务逻辑在实体内部

**代码示例：**
```go
// Product 产品实体
type Product struct {
    id          string
    name        string
    description string
    price       *valueobject.Money
    stock       int
}

// DecreaseStock 减少库存
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
```

#### 3. 聚合根 (Aggregate Root)
- **Order**: 订单聚合根，管理订单项和订单状态

**特性：**
- 聚合入口：外部只能通过聚合根访问聚合内部对象
- 一致性维护：负责维护聚合内部的一致性
- 事件发布：负责发布领域事件
- 边界控制：控制聚合边界内的所有操作

**代码示例：**
```go
// Order 订单聚合根
type Order struct {
    id              string
    userID          string
    items           []*entity.OrderItem
    shippingAddress *valueobject.Address
    status          valueobject.OrderStatus
    totalAmount     *valueobject.Money
    domainEvents    []event.DomainEvent
}

// AddItem 添加订单项
func (o *Order) AddItem(item *entity.OrderItem) error {
    // 聚合根维护内部一致性
    for _, existingItem := range o.items {
        if existingItem.Product().ID() == item.Product().ID() {
            newQty := existingItem.Quantity() + item.Quantity()
            if err := existingItem.UpdateQuantity(newQty); err != nil {
                return err
            }
            return o.recalculateTotalAmount()
        }
    }
    o.items = append(o.items, item)
    return o.recalculateTotalAmount()
}
```

#### 4. 领域事件 (Domain Event)
- **OrderCreatedEvent**: 订单创建事件
- **OrderPaidEvent**: 订单支付事件
- **OrderShippedEvent**: 订单发货事件
- **OrderCancelledEvent**: 订单取消事件

**特性：**
- 表示领域中发生的重要事情
- 用于解耦领域对象之间的依赖
- 支持事件驱动架构

**代码示例：**
```go
// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    BaseEvent
    OrderID   string
    UserID    string
    TotalAmount float64
}

// EventType 返回事件类型
func (e *OrderCreatedEvent) EventType() string {
    return "OrderCreated"
}
```

#### 5. 领域服务 (Domain Service)
- **OrderPricingService**: 订单定价服务
- **ProductInventoryService**: 产品库存服务

**特性：**
- 不自然属于任何实体或值对象的业务逻辑
- 无状态
- 以领域概念命名

**代码示例：**
```go
// OrderPricingService 订单定价服务
type OrderPricingService struct {
    discountRate float64
}

// CalculateDiscountedTotal 计算折扣后的订单总价
func (s *OrderPricingService) CalculateDiscountedTotal(order *aggregate.Order) (float64, error) {
    if order == nil {
        return 0, errors.New("订单不能为空")
    }
    totalAmount := order.TotalAmount()
    if totalAmount == nil {
        return 0, nil
    }
    discountedTotal := totalAmount.Amount() * (1 - s.discountRate)
    return discountedTotal, nil
}
```

#### 6. 工厂 (Factory)
- **OrderFactory**: 订单工厂
- **ProductFactory**: 产品工厂

**特性：**
- 负责复杂对象的创建
- 封装创建逻辑
- 确保创建的对象处于有效状态

**代码示例：**
```go
// OrderFactory 订单工厂
type OrderFactory struct {
    idGenerator IDGenerator
}

// CreateOrder 创建订单
func (f *OrderFactory) CreateOrder(
    userID string,
    shippingAddress *valueobject.Address,
    productItems []ProductItem,
) (*aggregate.Order, error) {
    orderID := f.idGenerator.Generate()
    order, err := aggregate.NewOrder(orderID, userID, shippingAddress)
    if err != nil {
        return nil, err
    }
    
    for _, item := range productItems {
        itemID := f.idGenerator.Generate()
        orderItem, err := entity.NewOrderItem(itemID, item.Product, item.Quantity)
        if err != nil {
            return nil, err
        }
        if err := order.AddItem(orderItem); err != nil {
            return nil, err
        }
    }
    
    return order, nil
}
```

#### 7. 仓储接口 (Repository Interface)
- **OrderRepository**: 订单仓储接口
- **ProductRepository**: 产品仓储接口

**特性：**
- 领域层定义接口，基础设施层实现
- 提供聚合根的持久化操作
- 隐藏数据访问细节
- 只针对聚合根定义仓储

**代码示例：**
```go
// OrderRepository 订单仓储接口
type OrderRepository interface {
    Save(order *aggregate.Order) error
    FindByID(id string) (*aggregate.Order, error)
    FindByUserID(userID string) ([]*aggregate.Order, error)
    FindAll() ([]*aggregate.Order, error)
    Delete(id string) error
}
```

### 2.2 应用层 (Application Layer)
应用层负责协调领域对象完成业务用例，不包含业务逻辑。

**目录结构：**
```
ordersys/
├── application/
│   ├── dto/                  # 数据传输对象
│   └── service/              # 应用服务
```

**核心组件：**

#### 1. DTO (Data Transfer Object)
- **CreateOrderRequest**: 创建订单请求
- **OrderResponse**: 订单响应
- **ProductResponse**: 产品响应

**特性：**
- 用于应用层与外部交互的数据传输
- 与领域对象分离，避免暴露领域模型细节

**代码示例：**
```go
// CreateOrderRequest 创建订单请求DTO
type CreateOrderRequest struct {
    UserID          string
    Province        string
    City            string
    District        string
    Street          string
    ZipCode         string
    ProductItems    []ProductItemDTO
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
```

#### 2. 应用服务 (Application Service)
- **OrderApplicationService**: 订单应用服务

**特性：**
- 协调领域对象完成业务用例
- 不包含业务逻辑，业务逻辑在领域层
- 负责事务管理和权限控制
- DTO与领域对象之间的转换
- 发布领域事件

**代码示例：**
```go
// OrderApplicationService 订单应用服务
type OrderApplicationService struct {
    orderRepo         repository.OrderRepository
    productRepo       repository.ProductRepository
    orderFactory      *factory.OrderFactory
    productFactory    *factory.ProductFactory
    pricingService    *service.OrderPricingService
    inventoryService  *service.ProductInventoryService
    eventPublisher    event.EventPublisher
}

// CreateOrder 创建订单
func (s *OrderApplicationService) CreateOrder(req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
    // 创建地址值对象
    shippingAddress, err := valueobject.NewAddress(
        req.Province, req.City, req.District, req.Street, req.ZipCode,
    )
    if err != nil {
        return nil, err
    }
    
    // 准备产品项
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
    
    // 使用工厂创建订单
    order, err := s.orderFactory.CreateOrder(req.UserID, shippingAddress, productItems)
    if err != nil {
        return nil, err
    }
    
    // 保存订单
    if err := s.orderRepo.Save(order); err != nil {
        return nil, err
    }
    
    // 发布领域事件
    if err := s.publishDomainEvents(order); err != nil {
        return nil, err
    }
    
    return s.toOrderResponse(order), nil
}
```

### 2.3 基础设施层 (Infrastructure Layer)
基础设施层提供技术支持，实现领域层定义的接口。

**目录结构：**
```
ordersys/
├── infrastructure/
│   ├── persistence/          # 持久化实现
│   └── id_generator.go       # ID生成器
```

**核心组件：**

#### 1. 仓储实现
- **InMemoryOrderRepository**: 内存订单仓储实现
- **InMemoryProductRepository**: 内存产品仓储实现

**特性：**
- 实现领域层定义的仓储接口
- 提供数据持久化能力
- 可替换为数据库实现

**代码示例：**
```go
// InMemoryOrderRepository 内存订单仓储实现
type InMemoryOrderRepository struct {
    orders map[string]*aggregate.Order
    mu     sync.RWMutex
}

// Save 保存订单
func (r *InMemoryOrderRepository) Save(order *aggregate.Order) error {
    if order == nil {
        return fmt.Errorf("订单不能为空")
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    r.orders[order.ID()] = order
    return nil
}

// FindByID 根据ID查找订单
func (r *InMemoryOrderRepository) FindByID(id string) (*aggregate.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    order, exists := r.orders[id]
    if !exists {
        return nil, fmt.Errorf("订单不存在: %s", id)
    }
    return order, nil
}
```

#### 2. ID生成器
- **UUIDGenerator**: UUID生成器

**特性：**
- 实现领域层定义的IDGenerator接口
- 提供唯一ID生成能力

**代码示例：**
```go
// UUIDGenerator UUID生成器
type UUIDGenerator struct {
    counter int
}

// Generate 生成唯一ID
func (g *UUIDGenerator) Generate() string {
    g.counter++
    return fmt.Sprintf("ID-%d", g.counter)
}
```

## 三、DDD核心概念实现

### 3.1 充血模型
本项目采用充血模型，业务逻辑在领域对象内部：

```go
// 订单聚合根包含丰富的业务方法
func (o *Order) Pay(paymentMethod string) error {
    // 业务逻辑在领域对象内部
    if err := o.status.CanTransitionTo(valueobject.OrderStatusPaid); err != nil {
        return err
    }
    
    if o.totalAmount == nil || o.totalAmount.Amount() == 0 {
        return errors.New("订单金额为0，无需支付")
    }
    
    o.status = valueobject.OrderStatusPaid
    
    // 发布领域事件
    o.addDomainEvent(&event.OrderPaidEvent{
        BaseEvent:     event.NewBaseEvent(),
        OrderID:       o.id,
        PaidAmount:    o.totalAmount.Amount(),
        PaymentMethod: paymentMethod,
    })
    
    return nil
}
```

**优势：**
- 业务逻辑集中，易于理解和维护
- 领域对象自包含，高内聚
- 避免贫血模型，业务逻辑分散的问题

### 3.2 聚合边界
订单聚合根控制聚合边界：

```go
// 外部只能通过聚合根访问聚合内部对象
func (o *Order) AddItem(item *entity.OrderItem) error {
    // 聚合根维护内部一致性
    for _, existingItem := range o.items {
        if existingItem.Product().ID() == item.Product().ID() {
            // 合并相同产品的订单项
            newQty := existingItem.Quantity() + item.Quantity()
            if err := existingItem.UpdateQuantity(newQty); err != nil {
                return err
            }
            return o.recalculateTotalAmount()
        }
    }
    o.items = append(o.items, item)
    return o.recalculateTotalAmount()
}
```

**聚合设计原则：**
- 订单是聚合根，订单项是聚合内部实体
- 外部不能直接操作订单项，必须通过订单
- 订单负责维护订单项的一致性
- 订单负责发布领域事件

### 3.3 领域事件
领域事件用于解耦领域对象之间的依赖：

```go
// 订单创建时发布领域事件
order.addDomainEvent(&event.OrderCreatedEvent{
    BaseEvent: event.NewBaseEvent(),
    OrderID:   id,
    UserID:    userID,
})

// 应用服务负责发布事件
func (s *OrderApplicationService) publishDomainEvents(order *aggregate.Order) error {
    events := order.DomainEvents()
    for _, evt := range events {
        if err := s.eventPublisher.Publish(evt); err != nil {
            return err
        }
    }
    order.ClearDomainEvents()
    return nil
}
```

**事件处理器示例：**
```go
// OrderCreatedEventHandler 订单创建事件处理器
type OrderCreatedEventHandler struct{}

func (h *OrderCreatedEventHandler) Handle(evt event.DomainEvent) error {
    orderCreatedEvt := evt.(*event.OrderCreatedEvent)
    fmt.Printf("📧 发送通知：订单 %s 已创建，用户ID: %s\n", 
        orderCreatedEvt.OrderID, orderCreatedEvt.UserID)
    return nil
}
```

### 3.4 仓储模式
仓储接口在领域层定义，实现在基础设施层：

```go
// 领域层定义接口
type OrderRepository interface {
    Save(order *aggregate.Order) error
    FindByID(id string) (*aggregate.Order, error)
    FindByUserID(userID string) ([]*aggregate.Order, error)
    FindAll() ([]*aggregate.Order, error)
    Delete(id string) error
}

// 基础设施层实现
type InMemoryOrderRepository struct {
    orders map[string]*aggregate.Order
    mu     sync.RWMutex
}
```

**仓储设计原则：**
- 只为聚合根定义仓储
- 隐藏数据访问细节
- 提供类似集合的接口
- 支持事务管理

### 3.5 工厂模式
工厂负责复杂对象的创建：

```go
func (f *OrderFactory) CreateOrder(
    userID string,
    shippingAddress *valueobject.Address,
    productItems []ProductItem,
) (*aggregate.Order, error) {
    // 封装订单创建的复杂逻辑
    orderID := f.idGenerator.Generate()
    order, err := aggregate.NewOrder(orderID, userID, shippingAddress)
    if err != nil {
        return nil, err
    }
    
    for _, item := range productItems {
        itemID := f.idGenerator.Generate()
        orderItem, err := entity.NewOrderItem(itemID, item.Product, item.Quantity)
        if err != nil {
            return nil, err
        }
        if err := order.AddItem(orderItem); err != nil {
            return nil, err
        }
    }
    
    return order, nil
}
```

**工厂设计原则：**
- 封装复杂对象的创建逻辑
- 确保创建的对象处于有效状态
- 隐藏创建细节
- 支持依赖注入

## 四、业务流程示例

### 4.1 创建订单流程
```
1. 应用服务接收创建订单请求 (DTO)
   ↓
2. 创建地址值对象
   ↓
3. 从仓储加载产品实体
   ↓
4. 通过工厂创建订单聚合根
   ↓
5. 保存订单到仓储
   ↓
6. 发布订单创建事件
   ↓
7. 返回订单响应 (DTO)
```

### 4.2 支付订单流程
```
1. 应用服务接收支付请求
   ↓
2. 从仓储加载订单聚合根
   ↓
3. 调用订单的Pay方法（业务逻辑在领域对象内部）
   ↓
4. 保存订单到仓储
   ↓
5. 发布订单支付事件
   ↓
6. 返回支付结果
```

### 4.3 发货流程
```
1. 应用服务接收发货请求
   ↓
2. 从仓储加载订单聚合根
   ↓
3. 调用订单的Ship方法
   ↓
4. 保存订单到仓储
   ↓
5. 发布订单发货事件
   ↓
6. 返回发货结果
```

### 4.4 完整业务流程演示

```go
func main() {
    // 初始化基础设施
    idGenerator := infrastructure.NewUUIDGenerator()
    orderRepo := persistence.NewInMemoryOrderRepository()
    productRepo := persistence.NewInMemoryProductRepository()
    eventPublisher := event.NewInMemoryEventPublisher()
    
    // 订阅事件
    eventPublisher.Subscribe("OrderCreated", &OrderCreatedEventHandler{})
    eventPublisher.Subscribe("OrderPaid", &OrderPaidEventHandler{})
    eventPublisher.Subscribe("OrderShipped", &OrderShippedEventHandler{})
    
    // 创建应用服务
    orderAppService := service.NewOrderApplicationService(...)
    
    // 创建产品
    product1, _ := orderAppService.CreateProduct("iPhone 15 Pro", "最新款苹果手机", 7999.00, "CNY", 100)
    product2, _ := orderAppService.CreateProduct("MacBook Pro", "专业笔记本电脑", 14999.00, "CNY", 50)
    
    // 创建订单
    order, _ := orderAppService.CreateOrder(&dto.CreateOrderRequest{
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
    })
    
    // 支付订单
    orderAppService.PayOrder(order.OrderID, "支付宝")
    
    // 发货
    orderAppService.ShipOrder(order.OrderID, "SF1234567890")
}
```

## 五、设计原则

### 5.1 依赖倒置原则 (DIP)
- 领域层定义接口，基础设施层实现
- 领域层不依赖具体实现
- 高层模块不依赖低层模块，都依赖抽象

**示例：**
```go
// 领域层定义接口
type IDGenerator interface {
    Generate() string
}

// 基础设施层实现
type UUIDGenerator struct {
    counter int
}

func (g *UUIDGenerator) Generate() string {
    g.counter++
    return fmt.Sprintf("ID-%d", g.counter)
}
```

### 5.2 单一职责原则 (SRP)
- 每个领域对象只负责一个业务概念
- 应用服务只负责协调，不包含业务逻辑
- 仓储只负责持久化

**示例：**
```go
// 订单聚合根只负责订单相关的业务逻辑
type Order struct {
    id              string
    userID          string
    items           []*entity.OrderItem
    shippingAddress *valueobject.Address
    status          valueobject.OrderStatus
    totalAmount     *valueobject.Money
    domainEvents    []event.DomainEvent
}

// 仓储只负责持久化
type OrderRepository interface {
    Save(order *aggregate.Order) error
    FindByID(id string) (*aggregate.Order, error)
}
```

### 5.3 开闭原则 (OCP)
- 通过领域事件扩展功能，而不修改现有代码
- 通过仓储接口切换持久化实现
- 通过事件处理器扩展业务逻辑

**示例：**
```go
// 通过事件处理器扩展功能，不修改订单代码
type OrderCreatedEventHandler struct{}

func (h *OrderCreatedEventHandler) Handle(evt event.DomainEvent) error {
    // 可以在这里添加新的业务逻辑，如发送短信、邮件等
    // 不需要修改订单聚合根的代码
    return nil
}
```

### 5.4 接口隔离原则 (ISP)
- 仓储接口只包含必要的方法
- 领域服务接口细粒度设计
- 客户端不应该依赖它不需要的接口

**示例：**
```go
// 仓储接口只包含必要的方法
type OrderRepository interface {
    Save(order *aggregate.Order) error
    FindByID(id string) (*aggregate.Order, error)
    FindByUserID(userID string) ([]*aggregate.Order, error)
    FindAll() ([]*aggregate.Order, error)
    Delete(id string) error
}
```

## 六、最佳实践

### 6.1 错误处理
- 领域对象构造函数进行自验证
- 业务方法返回错误，不使用panic
- 应用服务统一处理错误

**示例：**
```go
// 构造函数自验证
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

// 业务方法返回错误
func (o *Order) Pay(paymentMethod string) error {
    if err := o.status.CanTransitionTo(valueobject.OrderStatusPaid); err != nil {
        return err
    }
    // ...
    return nil
}
```

### 6.2 并发安全
- 仓储实现使用读写锁保护共享资源
- 领域对象本身无状态，天然并发安全
- 使用sync.RWMutex保护并发访问

**示例：**
```go
type InMemoryOrderRepository struct {
    orders map[string]*aggregate.Order
    mu     sync.RWMutex
}

func (r *InMemoryOrderRepository) Save(order *aggregate.Order) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.orders[order.ID()] = order
    return nil
}

func (r *InMemoryOrderRepository) FindByID(id string) (*aggregate.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    order, exists := r.orders[id]
    if !exists {
        return nil, fmt.Errorf("订单不存在: %s", id)
    }
    return order, nil
}
```

### 6.3 测试友好
- 领域对象可独立测试
- 通过接口注入依赖，便于mock
- 使用依赖注入提高可测试性

**示例：**
```go
// 领域对象可独立测试
func TestOrder_Pay(t *testing.T) {
    address, _ := valueobject.NewAddress("广东省", "深圳市", "南山区", "科技园路100号", "518000")
    order, _ := aggregate.NewOrder("ORDER-001", "USER-001", address)
    
    // 添加订单项
    money, _ := valueobject.NewMoney(100.0, "CNY")
    product, _ := entity.NewProduct("PRODUCT-001", "测试产品", "描述", money, 10)
    orderItem, _ := entity.NewOrderItem("ITEM-001", product, 1)
    order.AddItem(orderItem)
    
    // 测试支付
    err := order.Pay("支付宝")
    assert.NoError(t, err)
    assert.Equal(t, valueobject.OrderStatusPaid, order.Status())
}
```

### 6.4 代码风格
- 遵循Go语言命名规范
- 使用有意义的变量名和方法名
- 添加详细的注释说明

## 七、扩展建议

### 7.1 持久化扩展
- 实现基于数据库的仓储（如MySQL、PostgreSQL）
- 使用ORM框架（如GORM）简化实现
- 支持事务管理

**示例：**
```go
type MySQLOrderRepository struct {
    db *gorm.DB
}

func (r *MySQLOrderRepository) Save(order *aggregate.Order) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 保存订单
        // 保存订单项
        // 保证事务一致性
        return nil
    })
}
```

### 7.2 事件扩展
- 实现基于消息队列的事件发布器（如Kafka、RabbitMQ）
- 实现事件溯源（Event Sourcing）
- 支持事件重放

**示例：**
```go
type KafkaEventPublisher struct {
    producer sarama.SyncProducer
    topic    string
}

func (p *KafkaEventPublisher) Publish(event event.DomainEvent) error {
    message := &sarama.ProducerMessage{
        Topic: p.topic,
        Value: sarama.ByteEncoder(eventToJSON(event)),
    }
    _, _, err := p.producer.SendMessage(message)
    return err
}
```

### 7.3 应用层扩展
- 添加事务管理
- 添加权限控制
- 添加日志记录
- 添加缓存

**示例：**
```go
type OrderApplicationService struct {
    orderRepo       repository.OrderRepository
    transactionMgr  TransactionManager
    logger          Logger
    cache           Cache
}

func (s *OrderApplicationService) CreateOrder(req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
    // 记录日志
    s.logger.Info("Creating order", "userID", req.UserID)
    
    // 检查缓存
    if cached := s.cache.Get(req.UserID); cached != nil {
        return cached.(*dto.OrderResponse), nil
    }
    
    // 事务管理
    return s.transactionMgr.ExecuteInTransaction(func() (interface{}, error) {
        // 业务逻辑
        order, err := s.orderFactory.CreateOrder(...)
        if err != nil {
            return nil, err
        }
        
        if err := s.orderRepo.Save(order); err != nil {
            return nil, err
        }
        
        return s.toOrderResponse(order), nil
    })
}
```

### 7.4 接口层扩展
- 添加REST API
- 添加GraphQL API
- 添加gRPC API

**示例：**
```go
// REST API
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req dto.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    order, err := h.orderService.CreateOrder(&req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, order)
}
```

## 八、项目文件结构

```
go_test/
├── ordersys/                          # DDD订单系统
│   ├── domain/                        # 领域层
│   │   ├── model/
│   │   │   ├── valueobject/           # 值对象
│   │   │   │   ├── money.go           # 金额值对象
│   │   │   │   ├── address.go         # 地址值对象
│   │   │   │   └── order_status.go    # 订单状态值对象
│   │   │   ├── entity/                # 实体
│   │   │   │   ├── product.go         # 产品实体
│   │   │   │   └── order_item.go      # 订单项实体
│   │   │   ├── aggregate/             # 聚合根
│   │   │   │   └── order.go           # 订单聚合根
│   │   │   └── event/                 # 领域事件
│   │   │       ├── domain_event.go    # 领域事件定义
│   │   │       └── event_publisher.go # 事件发布器
│   │   ├── service/                   # 领域服务
│   │   │   └── domain_service.go      # 领域服务实现
│   │   ├── factory/                   # 工厂
│   │   │   └── order_factory.go       # 订单工厂
│   │   └── repository/                # 仓储接口
│   │       └── repository.go          # 仓储接口定义
│   ├── application/                   # 应用层
│   │   ├── dto/                       # 数据传输对象
│   │   │   └── order_dto.go           # 订单DTO
│   │   └── service/                   # 应用服务
│   │       └── order_service.go       # 订单应用服务
│   └── infrastructure/                # 基础设施层
│       ├── persistence/               # 持久化实现
│       │   └── memory_repository.go   # 内存仓储实现
│       └── id_generator.go            # ID生成器
├── main.go                            # 示例主程序
├── go.mod                             # Go模块定义
└── ARCHITECTURE.md                    # 架构设计文档
```

## 九、运行示例

### 9.1 运行命令
```bash
go run main.go
```

### 9.2 运行结果
```
========================================
  DDD电商订单系统示例
========================================

步骤1: 创建产品
----------------------------------------
✓ 产品创建成功: iPhone 15 Pro (ID: ID-1, 价格: 7999.00 CNY, 库存: 100)
✓ 产品创建成功: MacBook Pro (ID: ID-2, 价格: 14999.00 CNY, 库存: 50)

步骤2: 创建订单
----------------------------------------
📧 发送通知：订单 ID-3 已创建，用户ID: USER-001
✓ 订单创建成功
  订单ID: ID-3
  用户ID: USER-001
  订单状态: 待支付
  总金额: 30997.00 CNY
  收货地址: 广东省深圳市南山区科技园路100号
  订单项:
    - iPhone 15 Pro x2, 单价: 7999.00, 小计: 15998.00
    - MacBook Pro x1, 单价: 14999.00, 小计: 14999.00

步骤3: 支付订单
----------------------------------------
💳 支付成功：订单 ID-3，支付金额: 30997.00，支付方式: 支付宝
✓ 订单支付成功
  订单状态: 已支付

步骤4: 发货
----------------------------------------
🚚 订单发货：订单 ID-3，物流单号: SF1234567890
✓ 订单发货成功
  订单状态: 已发货

步骤5: 查询用户订单列表
----------------------------------------
✓ 用户 USER-001 共有 1 个订单
  - 订单ID: ID-3, 状态: 已发货, 金额: 30997.00

步骤6: 创建并取消订单
----------------------------------------
📧 发送通知：订单 ID-6 已创建，用户ID: USER-002
✓ 订单创建成功: ID-6
❌ 订单取消：订单 ID-6，原因: 用户主动取消
✓ 订单取消成功
  订单状态: 已取消

========================================
  示例执行完成
========================================
```

## 十、总结

本项目完整展示了DDD在Go语言项目中的应用方式，包括：

### 10.1 核心成果
1. **领域模型设计**：实体、值对象、聚合根的合理划分
2. **业务逻辑封装**：充血模型，业务逻辑在领域对象内部
3. **分层架构**：领域层、应用层、基础设施层清晰分离
4. **设计模式应用**：工厂、仓储、领域事件等模式的实践
5. **最佳实践**：错误处理、并发安全、测试友好等

### 10.2 关键特性
- ✅ **充血模型**：业务逻辑在领域对象内部，避免贫血模型
- ✅ **聚合边界**：聚合根控制聚合内部对象，维护一致性
- ✅ **领域事件**：解耦领域对象，支持事件驱动架构
- ✅ **依赖倒置**：领域层定义接口，基础设施层实现
- ✅ **工厂模式**：封装复杂对象的创建逻辑
- ✅ **仓储模式**：隐藏数据访问细节，提供持久化能力

### 10.3 学习价值
通过本示例，可以深入理解：
- DDD的核心概念和实践方法
- 如何在Go语言中实现DDD架构
- 如何设计领域模型和业务逻辑
- 如何组织代码结构和分层架构
- 如何应用设计模式和最佳实践

### 10.4 实际应用
本示例可以作为实际项目的参考模板，适用于：
- 电商系统
- 订单管理系统
- 库存管理系统
- 其他复杂业务系统

通过本示例的学习和实践，可以更好地理解和应用DDD，提高软件设计的质量和可维护性。
