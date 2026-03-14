package factory

import (
	"fmt"

	entity2 "go_test/ordersys/domain/model/entity"
	valueobject2 "go_test/ordersys/domain/model/vo"
)

// CreateProduct 静态工厂方法：创建产品（使用包变量 idGenerator）
func CreateProduct(
	name string,
	description string,
	price *valueobject2.Money,
	stock int,
) (*entity2.Product, error) {
	if idGenerator == nil {
		return nil, fmt.Errorf("工厂未设置 ID 生成器，请调用 factory.SetIDGenerator")
	}
	productID := idGenerator.Generate()
	return entity2.NewProduct(productID, name, description, price, stock)
}
