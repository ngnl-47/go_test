# 使用 asaskevich/EventBus 替代自研领域事件实现 — 评估报告

## 一、当前实现概览

| 组件 | 位置 | 职责 |
|------|------|------|
| `DomainEvent` 接口 + 具体事件体 | `domain/model/event/domain_event.go` | 领域事件定义（OrderCreated、OrderPaid 等） |
| `EventHandler`、`EventPublisher` 接口 | `domain/model/event/event_publisher.go` | 领域层抽象：发布/订阅契约 |
| `InMemoryEventPublisher` | 同上 | 内存版发布器实现（约 50 行） |

当前用法：按事件类型字符串订阅，发布时把 `DomainEvent` 交给对应 `EventHandler.Handle(evt)`，同步执行，任一 handler 返回 error 则中止并返回。

---

## 二、asaskevich/EventBus 能力对照

| 能力 | 当前实现 | EventBus |
|------|----------|----------|
| 按 topic 订阅 | ✅ `Subscribe(eventType, handler)` | ✅ `Subscribe(topic, fn)` |
| 发布事件 | ✅ `Publish(evt)`，按 `evt.EventType()` 路由 | ✅ `Publish(topic, args...)`，topic 自定 |
| 同步执行 | ✅ 顺序执行 handlers | ✅ 默认同步 |
| Handler 返回 error | ✅ 支持，并短路后续 handler | ❌ 回调为 `interface{}`，无返回值 |
| 异步发布 | ❌ 无 | ✅ `SubscribeAsync` + `WaitAsync` |
| 退订 / 单次订阅 | ❌ 无 | ✅ `Unsubscribe`、`SubscribeOnce` |

结论：**可以替代**。用 EventBus 做“总线”，在**基础设施层**做一层薄适配，实现现有 `EventPublisher` 接口，即可在不改领域层、应用层契约的前提下替换掉 `InMemoryEventPublisher`。

---

## 三、替换方案与注意点

### 3.1 建议做法：保留领域抽象 + EventBus 做实现

- **保留**（领域层）  
  - `DomainEvent` 及所有具体事件类型（`OrderCreatedEvent` 等）  
  - `EventPublisher` 接口（以及可选保留 `EventHandler` 接口，便于现有 handler 复用）

- **删除**（基础设施/实现）  
  - `InMemoryEventPublisher` 的完整实现

- **新增**（基础设施层）  
  - 一个适配器：实现 `event.EventPublisher`，内部持有一个 `EventBus`，  
    - `Subscribe(eventType, handler)` → 用 EventBus 的 `Subscribe(eventType, 包装函数)`，包装函数内将参数断言为 `DomainEvent` 后调用 `handler.Handle(evt)`  
    - `Publish(evt)` → `bus.Publish(evt.EventType(), evt)`  
  - 若需保留“第一个 handler 报错即返回”的语义，可在适配器内顺序执行已注册的 handler（或对 EventBus 的同步 Publish 在包装函数里收集/返回 error）。

这样：

- 领域层、应用层仍只依赖 `event.EventPublisher` 和 `event.DomainEvent`，**不依赖** EventBus。
- 后续若要换成 Kafka/RabbitMQ 等，只需再写一个实现 `EventPublisher` 的适配器，符合 DIP。

### 3.2 行为差异与处理

| 点 | 说明 | 建议 |
|----|------|------|
| Handler 错误返回 | EventBus 的 `Publish` 不返回 error，回调签名任意 | 在适配器的包装函数里调用 `handler.Handle(evt)`，若返回 error 则记录并决定是否短路/聚合后由适配器的 `Publish` 返回 |
| 并发安全 | 当前 `InMemoryEventPublisher` 用 `sync.RWMutex` 保护 map | EventBus 自身有并发保护，适配器无需再锁 |
| 异步能力 | 当前为同步 | 若将来需要“发事件不阻塞主流程”，可用 EventBus 的 `SubscribeAsync` + `WaitAsync` 在适配器内实现 |

### 3.3 可简化与可扩展

- **可简化**：删除自研的 `InMemoryEventPublisher` 和其内部 map/锁逻辑，减少自维护代码。
- **可扩展**：  
  - 需要异步、单次订阅、退订时，可直接在适配器里委托给 EventBus。  
  - 需要跨进程时，可继续用当前 `EventPublisher` 抽象，在基础设施层做 Kafka/RPC 等实现，而不动领域事件定义。

---

## 四、结论与建议

- **是否可以用 EventBus 替代当前领域事件“发布/订阅”的独立实现？**  
  **可以。** 用“保留领域接口 + EventBus 做基础设施层实现”的方式，既能替代自研内存发布器，又便于后续维护和扩展。

- **推荐做法**：  
  1. 保留 `domain/model/event/domain_event.go` 和 `event_publisher.go` 中的**接口与领域事件定义**。  
  2. 删除 `InMemoryEventPublisher` 实现。  
  3. 在 `infrastructure`（或 `infrastructure/event`）中新增 **EventBus 适配器**，实现 `event.EventPublisher`，内部用 `github.com/asaskevich/EventBus` 做 topic 订阅与发布。  
  4. 在 `main.go` 或组装处将 `EventPublisher` 依赖改为该适配器。

这样既用上了三方 EventBus、简化了自研维护，又保持了 DDD 中“领域事件接口在领域层、实现可替换”的分层与可维护性。
