package main

import "strings"

// 本模块为库与测试工程；订单系统行为验证请执行：
//
//	go test ./ordersys/...

func main() {
	occupyKey := "xyjh_online_mix:legion:war:occupy:1:26032:10:-1:-1:11919"
	println(strings.Contains(occupyKey, "260326"))
}
