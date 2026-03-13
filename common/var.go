package common

import "github.com/asaskevich/EventBus"

// GlobalEvent 全局事件总线
var GlobalEvent EventBus.Bus = EventBus.New()
