package vo

import (
	"errors"
	"fmt"
)

// Address 地址值对象
// 值对象特性：描述领域对象的某个属性，无唯一标识
type Address struct {
	province string
	city     string
	district string
	street   string
	zipCode  string
}

// NewAddress 创建地址值对象
func NewAddress(province, city, district, street, zipCode string) (*Address, error) {
	if province == "" || city == "" || district == "" || street == "" {
		return nil, errors.New("地址信息不完整")
	}
	return &Address{
		province: province,
		city:     city,
		district: district,
		street:   street,
		zipCode:  zipCode,
	}, nil
}

// Province 获取省份
func (a *Address) Province() string {
	return a.province
}

// City 获取城市
func (a *Address) City() string {
	return a.city
}

// District 获取区县
func (a *Address) District() string {
	return a.district
}

// Street 获取街道地址
func (a *Address) Street() string {
	return a.street
}

// ZipCode 获取邮政编码
func (a *Address) ZipCode() string {
	return a.zipCode
}

// Equals 判断两个地址是否相等
func (a *Address) Equals(other *Address) bool {
	if other == nil {
		return false
	}
	return a.province == other.province &&
		a.city == other.city &&
		a.district == other.district &&
		a.street == other.street &&
		a.zipCode == other.zipCode
}

// String 返回地址的完整字符串表示
func (a *Address) String() string {
	return fmt.Sprintf("%s%s%s%s", a.province, a.city, a.district, a.street)
}
