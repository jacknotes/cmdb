package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "github.com/jacknotes/cmdb/apps/book/api"
	_ "github.com/jacknotes/cmdb/apps/host/api"
	_ "github.com/jacknotes/cmdb/apps/resource/api"
)
