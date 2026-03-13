package event

import "sync"

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event DomainEvent) error
}

// EventPublisher 事件发布器接口
// 领域层定义接口，基础设施层实现
type EventPublisher interface {
	Publish(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler)
}

// InMemoryEventPublisher 内存事件发布器
// 简单实现，用于演示
type InMemoryEventPublisher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventPublisher 创建内存事件发布器
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件
func (p *InMemoryEventPublisher) Publish(event DomainEvent) error {
	p.mu.RLock()
	handlers, exists := p.handlers[event.EventType()]
	p.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅事件
func (p *InMemoryEventPublisher) Subscribe(eventType string, handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[eventType] = append(p.handlers[eventType], handler)
}
