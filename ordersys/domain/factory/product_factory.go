package factory

import (
	"fmt"

	"go_test/ordersys/domain/model"
)

// CreateProduct 静态工厂方法：创建产品（使用包变量 idGenerator）
func CreateProduct(
	name string,
	description string,
	price *model.Money,
	stock int,
) (*model.Product, error) {
	if idGenerator == nil {
		return nil, fmt.Errorf("工厂未设置 ID 生成器，请调用 factory.SetIDGenerator")
	}
	productID := idGenerator.Generate()
	return model.NewProduct(productID, name, description, price, stock)
}
