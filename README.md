# goagent
TingYun APM3.0 - Go

### 嵌码
1. 抓取gin框架数据
	在工程中添加文件tingyun.go:

tingyun.go
```
package main

import (
	_ "github.com/TingYunGo/goagent/frameworks/gin"
)
```

2. 抓取go语言内置http框架数据

在工程中添加文件tingyun.go:

tingyun.go
```
package main

import (
	_ "github.com/TingYunGo/goagent"
)
```

2.  抓取数据库性能数据

在工程中添加文件tingyun.go

tingyun.go
```
package main

import (
	_ "github.com/TingYunGo/goagent/database"
)
```

2. 抓取redis性能数据

在工程中添加文件tingyun.go

tingyun.go (使用redigo驱动)
```
package main

import (
	_ "github.com/TingYunGo/goagent/nosql/redigo"
)
```

或者 (使用go-redis驱动)
```
package main

import (
	/* go mod 方式 缺省 v6, 或者禁用 go mod 时, 启用这行代码 */
	_ "github.com/TingYunGo/goagent/nosql/go-redis"
	
	/* 启用 go mod 方式, 使用 go-redis/v7 版本时, 启用这行代码 */
	_ "github.com/TingYunGo/goagent/nosql/go-redis/v7"
	/* 启用 go mod 方式, 使用 go-redis/v8 版本时, 启用这行代码 */
	_ "github.com/TingYunGo/goagent/nosql/go-redis/v8"
)
```

### 配置&运行
已经嵌码的程序需要通过环境变量指定配置文件路径
```
export TINGYUN_GO_APP_CONFIG=`pwd`/tingyun.json
```
配置文件格式
```
{
  "nbs.app_name" : "替换为您的应用名称",
  "nbs.license_key" : "替换为您的license",
  "nbs.host" : "替换为collector的ip:port",  
  "nbs.agent_enabled" : true,  
  "nbs.log_file_name" : "agent.log",
  "nbs.audit" : false,
  "nbs.max_log_count": 5,
  "nbs.max_log_size": 10,
  "nbs.action_cache_max" : 10000,
  "nbs.ssl" : false,
  "nbs.savecount" : 5
}
```

