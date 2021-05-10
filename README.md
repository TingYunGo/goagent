# goagent
TingYun APM3.0 - Go

## 嵌码
### 自动嵌码
1. gin框架自动嵌码
	
	在工程中引入 "github.com/TingYunGo/goagent/frameworks/gin"

	举例: 工程文件夹中添加tingyun.go文件如下
	tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/frameworks/gin"
)
```

2. go语言内置http框架自动嵌码

	在工程中引入 "github.com/TingYunGo/goagent"

	举例: 在工程中添加文件tingyun.go文件如下
	tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent"
)
```

3.  数据库自动嵌码

	在工程中引入 "github.com/TingYunGo/goagent/database"

	举例: 在工程中添加文件tingyun.go文件如下
	tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/database"
)
```

4. redis自动嵌码
	+ 4.1 redigo(github.com/gomodule/redigo)自动嵌码
	
	在工程中引入 "github.com/TingYunGo/goagent/nosql/redigo"

	举例: 在工程中添加文件tingyun.go文件如下
	tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/nosql/redigo"
)
```
	+ 4.2 go-redis(github.com/go-redis/redis)自动嵌码
		+ 4.2.1 gopath模式(非gomodule模式)
		在工程中引入 "github.com/TingYunGo/goagent/nosql/go-redis"
		举例: 在工程中添加文件tingyun.go文件如下
		tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/nosql/go-redis"
)
```

		+ 4.2.2 gomodule模式
			+ 4.2.2.1 go-redis v6(缺省)版本嵌码
			在工程中引入 "github.com/TingYunGo/goagent/nosql/go-redis"
			举例: 在工程中添加文件tingyun.go文件如下
			tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/nosql/go-redis"
)
```
			+ 4.2.2.2 go-redis v7版本嵌码
			在工程中引入 "github.com/TingYunGo/goagent/nosql/go-redis/v7"
			举例: 在工程中添加文件tingyun.go文件如下
			tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/nosql/go-redis/v7"
)
```
			+ 4.2.2.3 go-redis v8版本嵌码
			在工程中引入 "github.com/TingYunGo/goagent/nosql/go-redis/v8"
			举例: 在工程中添加文件tingyun.go文件如下
			tingyun.go:
```
package main
import (
	_ "github.com/TingYunGo/goagent/nosql/go-redis/v8"
)
```


## 配置&运行
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

