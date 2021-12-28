# 听云 goagent
**TingYun APM3.0 - GoAgent**

## 简介
  听云 goagent  能为您的 Go 语言应用程序提供运行时状态监控功能。
  
  它能帮助您追踪事务请求，外部调用，数据调用，nosql调用(redis, mongodb等) 以及自定义过程的性能数据和错误等运行状态信息。

  为方便理解,下面这部分简单介绍Go语言的特点,如果您熟悉Go语言,请转到 [听云 goagent 嵌码](#tingyun_goagent) 开始阅读。

### Go 语言简介
   与c和c++语言类似, Go 语言是编译类型的语言。

#### **Go语言程序如何运行?**
   Go语言源程序经过Go编译器编译,生成独立的二进制ELF(linux)/EXE(windows) 格式的可执行文件，运行时不再需要Go源程序参与。

#### **Go语言如何使用第三方模块?**
  Go语言通过import语句引入第三方模块
  示例代码:

  ```go
  package main
  import (
  	tingyun "github.com/TingYunGo/goagent"
  )
  func RoutineID() int64 {
  	return tingyun.GetGID()
  }
  ```
  其中: tingyun 为引入模块的别名。

  如果只是引入模块,并不直接使用模块的方法,那么引入模块的别名使用下划线 _ 代替, 否则Go编译器提示错误:

  ```go
  package main
  import (
  	_ "github.com/TingYunGo/goagent"
  )
  ```

#### **Go语言第三方模块如何下载安装到本地?**
  Go语言的所有第三方模块都是源码形式发布到git服务器上的。
  第三方模块的下载安装分为 **GOPATH** 模式和 **GOMOD** 模式两种情况:
  
 * **GOPATH模式** : Go语言版本低于1.11,或者禁用了GOMOD (设置环境变量GO111MODULE=off)模式时，处于GOPATH模式。
     <br/>这种模式下, Go编译器检测GOPATH环境变量,在GOPATH环境变量指定的每个路径下根据名字查找第三方库。GOPATH未设置时, 缺省值是当前用户的跟路径下的 go 文件夹。<br/>
     GOPATH模式第三方模块安装有两种方法:
	+ 自动安装:
	  <br/>使用命令 go get <第三方库路径>, 举例, 安装 tingyun goagent:
	  
	  ```bash
	  $ go get github.com/TingYunGo/goagent
	  ```
	  命令执行后,会自动下载模块到GOPATH下的src/github.com/TingYunGo/goagent 下, 并且,递归下载该模块依赖的模块到相应目录。
	  <br/>
	+ 手动安装:
	  <br/> 在GOPATH下手动创建路径src/github.com/TingYunGo/goagent , 并将代码复制到该文件夹下。

 * **GOMOD模式** : Go语言版本大于等于1.11, 并且设置环境变量GO111MODULE=on时， 处于GOMOD模式。
    <br/>GOMOD 模式下，在Go应用的根路径下需要go.mod文件，主要内容为应用名,Go语言版本和依赖包及版本。<br/>
    范例:
    ```
    module http_example
    go 1.12
    require (
    	github.com/TingYunGo/goagent v0.7.8
        github.com/golang/protobuf v1.5.2 // indirect
    )
    ```
    其中: 
    - http_example 是应用的名字。
    - go 1.12 : 是要求go版本不低于 v1.12 。
    - require:  这部分指定依赖的第三方模块及对应的版本。
    
    GOMOD模式下依赖包的下载：
    - 使用命令：
      ```bash 
      go mod tidy
      ```
      这个命令将自动检查当前应用的依赖并下载所有依赖包，并且校验依赖包的hash值,写入到go.sum文件。

#### **Go语言应用如何编译?**
   进入应用源码路径,执行 go build 命令，即生成应用的可执行文件。

<span id="tingyun_goagent">

## 听云 goagent 嵌码

###  听云 goagent 是什么?
  - 听云 goagent是一个Go语言第三方模块, 发布根路径是: github.com/TingYunGo/goagent
  - 听云 goagent 支持 amd64架构处理器 的linux环境， go1.9 到 最新的 go1.17.x Go语言版本。
  - 听云 goagent 提供自动嵌码和自定义嵌码(Go API)两种嵌码机制。
  - 自动嵌码支持列表参考  [框架支持列表](#tingyun_agent_frame)  和  [组件支持列表](#tingyun_agent_component) 

### 听云 goagent 如何安装?
  <br/>与所有第三方模块的安装方式相同。
  - GOPATH模式下安装: 
    ```bash
    $ go get github.com/TingYunGo/goagent
    ```
  - GOMOD模式下安装: 
    <br/>在应用文件夹下执行:
    ```
    $ go mod tidy
    ```

### 听云 goagent如何使用(嵌码)?
  根据应用使用http框架的不同, 需要import 不同的路径。
  
  我们以一个使用内置http框架举的简单例子说明如何嵌码: <br/>
  源文件: main.go 代码如下:
  
  ```go
  package main
  import (
  	"encoding/json"
  	"net/http"
  )
  func main() {
  	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  		header := w.Header()
  		header.Set("Cache-Control", "no-cache")
  		header.Set("Content-Type", "application/json; charset=utf-8")
  		w.WriteHeader(http.StatusOK)
  		b, _ := json.Marshal(map[string]interface{}{
  			"status": "success",
  			"URI":    r.URL.RawQuery,
  		})
  		w.Write(b)
  	})
  	http.ListenAndServe(":3000", nil)
  }
  ```
  要对如上这个应用嵌码, 在源文件同级目录下,创建tingyun.go文件,内容如下:

  ```
  package main
  import (
  	_ "github.com/TingYunGo/goagent"
  )
  ```

  全部嵌码工作即完成。

  嵌码说明: 
   <br/>此例应用使用了内置http框架, 对于内置http框架, 需要 引入 github.com/TingYunGo/goagent 。

### 如何确定嵌码要使用哪个/哪些 import 模块路径?
  下边整理了 goagent 支持的 框架/组件支持列表, 通过查表获得引用依赖。

#### 获取当前应用的依赖模块:
  * GOMOD 方式: 查看go.mod文件,或者使用命令:
  ```bash
  	$ go mod graph
  ```

  * GOPATH 方式: 编译时使用 -a -v 参数:
  ```
  	$ go build -a -v
  ```

<span id="tingyun_agent_frame">

### 听云探针框架支持列表:

| 框架 | 听云探针嵌码 import 模块路径 | 支持版本 |
|-----|---------------------------|--------|
| net/http <br/> 内置http框架 | github.com/TingYunGo/goagent | go1.9 ~ go1.17.x |
| github.com/gin-gonic/gin <br/>gin框架 | github.com/TingYunGo/goagent/frameworks/gin | gin v1.3.0 ~ gin v1.7.4 |
| github.com/astaxie/beego <br/>beego框架: GOPATH模式  | github.com/TingYunGo/goagent/frameworks/beego/path/astaxie | beego v1.12.0 ~ beego v2.0.0-beta |
| github.com/beego/beego <br/>beego框架: GOPATH模式  | github.com/TingYunGo/goagent/frameworks/beego/path | beego v1.12.0 ~ beego v2.0.1 |
| github.com/beego/beego <br/> beego框架v1: GOMOD模式  | github.com/TingYunGo/goagent/frameworks/beego | beego v1.12.0 ~ beego v1.12.3 |
| github.com/beego/beego/v2 <br/> beego框架v2: GOMOD模式  | github.com/TingYunGo/goagent/frameworks/beego/v2 | beego v2.0.0 ~ beego v2.0.1 |
| github.com/labstack/echo <br/>echo 框架 GOPATH模式 | github.com/TingYunGo/goagent/frameworks/echo | echo v3.3.10 ~ echo v4.6.1 |
| github.com/labstack/echo/v4 <br/>echo 框架 V4 GOMOD模式 | github.com/TingYunGo/goagent/frameworks/echo/v4 | echo v4.0.0 ~ echo v4.6.1 |
| github.com/kataras/iris/v12 <br/> iris 框架 v12.1.x | github.com/TingYunGo/goagent/frameworks/iris/v12 | iris v12.1.0 ~ iris v12.1.8 |
| github.com/kataras/iris/v12 <br/>iris 框架 v12.2  | github.com/TingYunGo/goagent/frameworks/iris/v12/2 | iris v12.2.0-alpha ~ iris v12.2.0-alpha3 |

<span id="tingyun_agent_component">

### 听云探针组件支持列表

| 组件 | 听云探针嵌码 import 模块路径 | 支持版本 |
|-----|---------------------------|---------|
| database/sql <br/> 数据库 | github.com/TingYunGo/goagent/database | go1.9 ~ go1.17.x<br/>驱动列表:<br/>mssql: github.com/denisenkom/go-mssqldb v0.9.0 ~ v0.11.0 <br/>mysql: github.com/go-sql-driver/mysql v1.0.0 ~ v1.6.0 <br/> postgresql: github.com/lib/pq v1.0.0 ~ v1.10.3 <br> sqlite: github.com/mattn/go-sqlite3 v1.0.0 ~ v1.14.8 |
| github.com/gomodule/redigo <br/> redis: redigo | github.com/TingYunGo/goagent/nosql/redigo | v1.7.0 ~ v1.8.5 |
| github.com/go-redis/redis <br/> redis: go-redis, GOPATH模式 | github.com/TingYunGo/goagent/nosql/go-redis | v6.10.0 ~ v8.11.4 |
| github.com/go-redis/redis <br/> redis: go-redis default, GOMOD模式 | github.com/TingYunGo/goagent/nosql/go-redis | v6.10.0 ~ v8.11.4 |
| github.com/go-redis/redis/v7 <br/> redis: go-redis v7, GOMOD模式 | github.com/TingYunGo/goagent/nosql/go-redis/v7 | v7.0.0 ~ v7.4.1 |
| github.com/go-redis/redis/v8 <br/> redis: go-redis v8, GOMOD模式 | github.com/TingYunGo/goagent/nosql/go-redis/v8 | v8.0.0 ~ v8.11.4 |
| Go.mongodb.org/mongo-driver/mongo <br/> mongodb | github.com/TingYunGo/goagent/nosql/mongodb | v1.1.0 ~ v1.7.3 |

### 嵌码实例演示
以开源项目 photoprism 为例 : 
项目地址 https://github.com/photoprism/photoprism
* 步骤1. 首先克隆项目:
  ```bash
  $ git clone https://github.com/photoprism/photoprism.git
  ```

* 步骤2. 确定项目使用哪些框架和库:
  
  进入项目文件夹,查看 go.mod文件:
  
  ```bash
  $ cd photoprism
  $ cat go.mod
  ```
   我们会看到,此项目使用了gin框架, 数据库支持: 
   
   ```
   postgresql:(github.com/lib/pq)
   mysql:(github.com/go-sql-driver/mysql)
   sqlite:(github.com/mattn/go-sqlite3)
   ```

* 步骤3: 查表确定引用路径,添加源码:
   
   查看框架支持列表, 我们的嵌码操作需要引用两个路径:
   
   ```
   github.com/TingYunGo/goagent/frameworks/gin
   github.com/TingYunGo/goagent/database
   ```
   
   在代码目录 internal/photoprism下创建 tingyun.go文件,内容如下:
   ```go
   package photoprism
   import (
   	_ "github.com/TingYunGo/goagent/database"
   	_ "github.com/TingYunGo/goagent/frameworks/gin"
   )
   ```

* 步骤4. 执行 go mod tidy, 编译:
  
  在项目的go.mod文件所在的文件夹下,执行:

  ```bash
  $ go mod tidy
  $ make
  ```
  以上4个步骤完成后,项目编译和探针嵌码工作就全部完成。

### 常用Go语言第三方组件支持的数据库版本列表

| 组件 | 支持版本 |
|------|---------|
| mssql: github.com/denisenkom/go-mssqldb | SQL Server 2008 SP3 + |
| mysql: github.com/go-sql-driver/mysql | mysql 5.5+ |
|postgresql: github.com/lib/pq | postgresql 9.6+ |
| sqlite: github.com/mattn/go-sqlite3 | sqlite 3.8.5 ~ 3.36.0 |
| redis: redigo github.com/gomodule/redigo | redis 4.0.0+ |
| redis: go-redis github.com/go-redis/redis  |  redis 4.0.0+ |
| mongodb: go.mongodb.org/mongo-driver/mongo | mongodb 2.6.1+ |

## 配置&运行
### 已经嵌码的应用, 听云agent 部分如何配置?

  为保证应用程序的最大限度的安全,缺省不配置的情况下, 嵌码逻辑是处于禁用状态。

#### 指定配置文件:
  要启用agent监控,需要设置环境变量: TINGYUN_GO_APP_CONFIG 指定配置文件路径。
  ```bash
  $ export TINGYUN_GO_APP_CONFIG=/[configfilepath]/tingyun.conf
  $ /your app path/appname [your app args]
  ```

**注意！！！**
 **必须在应用程序启动前设置环境变量，在应用程序里设置环境变量是无效的。**

#### 听云agent 配置项: 
参考配置 tingyun.conf:
```conf

######## 应用配置项 ########

# 设置应用名字符串类型, 可选项, 缺省取进程名做应用名 
# app_name = "My Go App"

# Agent启用标志, bool 类型, 可选项, 缺省为 true
# agent_enabled = true


######## 授权 / 服务器配置项 #########

# 授权序列码, 字符串类型, 必选项, 不能为空 
#license_key = "999-999-999"

# collector服务器地址, 多个地址用逗号分隔, 必选项, 不能为空
#collector.address = "collector_ip:port"

# 向collector服务器发送请求是否启用ssl(https), 缺省值 false
# ssl = false


######## 日志配置项 ########

# 日志文件路径, 必选项. 置空或不配置此项时, 无日志输出
agent_log_file = agent.log

# 日志输出级别设置, 可设置级别: debug, info, error, off
# debug: 输出最多日志
#  info: 输出 info级别和 error级别日志
# error: 仅输出 error级别日志
#   off: 关闭日志输出
# 缺省值: info
agent_log_level = info

# 审计模式, bool类型, 缺省值 false; 配合日志级别, 控制日志的输出
# 设置为 false 时, 部分审计模式日志不输出
audit_mode = true

# 日志文件大小, 整数, 单位 MB, 缺省值10
# 日志文件大小超过此阈值时, 将创建新的日志文件, 旧的日志将依次更新日志文件名
agent_log_file_size = 10

# 保留日志文件个数, 整数, 缺省值3
# 日志文件个数超过此阈值将从最早的文件开始删除
agent_log_file_count = 3


######## 内存控制阈值 ########

# 事务数据采集对象在内存缓冲队列中存放的最大数量, 缺省值10000
# 超过此阈值意味着当前并发数过高, 后台工作协程的处理能力不足以消费完采集到的数据
# 为防止数据积压导致的应用内存无限制增长, 当事务缓冲队列超过此阈值时, 不再采集事务性能数据 
action_cache_max = 10000

# 每次向collector发送数据包含的最大事务数量, 缺省值5000 
action_report_max = 5000

# 向collector发送的数据缓冲队列长度, 缺省值10 
# 发送队列长度超过此阈值意味网络缓慢或者collector处理能力不足 
# 为防止数据积压导致的内存无限制增长, 当发送队列长度超过此阈值时, 新的发送请求将被丢弃 
report_queue_count = 5

# 每个事务可采集的最大组件调用次数, 缺省值3000
# 事务采集的组件次数超过此阈值后,此事务的采集过程将不再采集新的组件调用数据
agent_component_max = 3000

# sql语句的最大字节长度,缺省值5000
# 超过此阈值,采集的sql将被截断为阈值设定长度
agent_sql_size_max = 5000
```

#### 环境变量支持
  某些特殊场合,有可能需要通过环境变量灵活控制配置项。作为配置文件的补充,听云agent也提供了相应支持。
  听云agent支持的环境变量:
  ```bash
  agent_enabled  # 启用探针选项, 取值: true/false; 缺省值: true
  ```
  ```bash
  audit_mode  # 审计模式, 审计模式日志开启选项, 取值: true/false; 缺省值: false
  ```
  ```bash
  agent_log_level  # 日志级别设置, 取值: debug/info/error/off; 缺省值: info
  ```
  ```bash
  agent_log_file  # 日志文件路径
  ```
  ```bash
  license_key  # 授权序列码
  ```
  ```bash
  collectors  # collector服务器地址,多个地址以逗号分隔
  ```
  ```bash
  agent_init_delay  # 探针延时初始化时间,整数,单位秒, 缺省值1.  说明: 如果探针初始化过早,可能在应用开始listen之前初始化, 这种情况下探针抓不到应用listen的端口. 增加初始化延时以解决此问题.
  ```
  ```bash
  TINGYUN_GO_APP_NAME  # 应用名称
  ```

## 自动嵌码演示范例
  使用听云 goagent嵌码,请参考部分范例程序 [听云 goagent演示](https://github.com/TingYunGo/goagent_examples)
## 自定义嵌码(听云 goagent API)
  自定义嵌码部分请参考阅读 [听云 goagent API](https://github.com/TingYunGo/goagent/blob/master/api.md)

## Code License
听云 goagent 使用 [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) 协议发布.


## 交叉编译
### OSX 系统

引入听云探针后,MAC系统环境下的交叉编译(linux,amd64), 需要安装交叉编译工具(C编译器及LIBC库),如果未安装相关工具,可参考使用如下包(使用MUSL C):

```
$ brew install FiloSottile/musl-cross/musl-cross
```

Go语言项目 MUSL C的交叉编译命令请参考:

```
# CC: C语言交叉编译器
# GOARCH: 编译后的可执行文件运行的CPU架构
# GOOS: 编译后的可执行文件运行的操作系统
# "-extldflags -static" : 静态链接
$ CC=x86_64-linux-musl-gcc GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static"
```


