package factory

// IDGenerator ID 生成器接口（领域层定义，基础设施层实现）
type IDGenerator interface {
	Generate() string
}

// 包级别 ID 生成器，由应用层在 init 时注入，保证全局唯一
var idGenerator IDGenerator

// SetIDGenerator 设置静态工厂使用的 ID 生成器，应在应用启动时调用一次
func SetIDGenerator(g IDGenerator) {
	idGenerator = g
}
