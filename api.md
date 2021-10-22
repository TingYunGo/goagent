# 听云 goagent API
**TingYun APM3.0 - GoAgent**


-------------------
### goagent API
对于自动嵌码不能满足需求的情况,可以使用自定义嵌码(听云 goagent API)方案.
自定义嵌码相关知识:

#### **概念**
 为了方便对程序问题的追踪, 我们将应用程序执行过程拆分出 **事务** 和 **组件** 两个概念.
##### 1. **事务**:
  **事务** 是一次服务请求的处理过程，例如:一次web请求, 或者rpc调用的server端处理过程。
##### 2. **组件**:
  **组件** 是事务处理过程中各个子功能的处理过程,例如: rpc外部调用, 数据库调用, nosql调用, 消息队列访问的生产者和消费者请求,功能/逻辑算法的计算过程,模块封装的过程等.

#### SDK API
| API 函数声明 | 功能 |
|--------------|---------|
| [CreateAction](#CreateAction) <br/> func CreateAction(name string, method string) (*Action, error) | 创建事务对象. <br/> 如果是基于 net/http 的应用,请使用GetAction |
| [GetAction](#GetAction) <br/> func GetAction() *Action| 取当前协程上的事务对象(使用net/http的应用,探针会自动创建事务对象.使用本函数获取) |
|[(*Action).CreateExternalComponent](#Action_CreateExternalComponent)<br/>func (*Action) CreateExternalComponent(url string, method string) *Component|创建一个外部服务访问组件 |
|[(*Action).CreateMQComponent](#Action_CreateMQComponent)<br/>func (*Action) CreateMQComponent(vender string, isConsumer bool, host, queue string) *Component|创建一个消息队列访问组件 |
|[(*Action).CreateMongoComponent](#Action_CreateMongoComponent)<br/>func (*Action) CreateMongoComponent(host, database, collection, op, method string) *Component|创建mongos访问组件 |
|[(*Action).CreateComponent](#Action_CreateComponent)<br/>func (*Action) CreateComponent(method string) *Component| 创建自定义过程监控组件 |
|[(*Action).AddRequestParam](#Action_AddRequestParam)<br/>func (*Action) AddRequestParam(k string, v string)| 添加事务请求参数 |
|[(*Action).AddResponseParam](#Action_AddResponseParam)<br/>func (*Action) AddResponseParam(k string, v string)| 添加事务应答参数 |
|[(*Action).AddCustomParam](#Action_AddCustomParam)<br/>func (*Action) AddCustomParam(k string, v string)| 添加事务自定义参数 |
|[(*Action).GetTxData](#Action_GetTxData)<br/>func (*Action) GetTxData() string|用于跨应用追踪被调用端,获取事务数据|
|[(*Action).SetTrackID](#Action_SetTrackID)<br/>func (*Action) SetTrackID(id string)|用于跨应用追踪被调用端,传递跨应用追踪ID |
|[(*Action).SetName](#Action_SetName)<br/>func (*Action) SetName(name string, method string)|设置事务名, 参数同CreateAction |
|[(*Action).SetHTTPMethod](#Action_SetHTTPMethod)<br/>func (*Action) SetHTTPMethod(httpMethod string)|设置http请求方法类型(GET/POST/PUT/OPTIONS/HEAD) |
|[(*Action).SetURL](#Action_SetURL)<br/>func (*Action) SetURL(name string)|设置事务的URI |
|[(*Action).Ignore](#Action_Ignore)<br/>func (*Action) Ignore()|忽略本次事务数据采集 |
|[(*Action).SetError](#Action_SetError)<br/>func (*Action) SetError(e interface{})|采集事务错误信息 |
|[(*Action).Finish](#Action_Finish)<br/>func (*Action) Finish()|事务采集结束 |
|[(*Action).SetStatusCode](#Action_SetStatusCode)<br/>func (*Action) SetStatusCode(code uint16) int|采集事务状态码 |
|[(*Component).GetAction](#Component_GetAction)<br/>func (*Component) GetAction() *Action|取组件对应事务对象 |
|[(*Component).SetError](#Component_SetError)<br/>func (*Component) SetError(e interface{}, errType string, skipStack int) |采集错误 |
|[(*Component).Finish](#Component_Finish)<br/>func (*Component) Finish()|组件过程结束 |
|[(*Component).CreateTrackID](#Component_CreateTrackID)<br/>func (*Component) CreateTrackID() string|生成跨应用追踪ID |
|[(*Component).SetTxData](#Component_SetTxData)<br/>func (*Component) SetTxData(txData string)|接收被调用端返回的跨应用追踪ID |
|[(*Component).CreateComponent](#Component_CreateComponent)<br/>func (*Component) CreateComponent(method string) *Component |创建组件的子过程(组件再分解) |


<span id="CreateAction">CreateAction</span>
```go
/* 功能 : 创建事务组件
 * 参数 :
 *      name : 对应包/struct名
 *      method : 函数名
 * 返回值 :
 *      (事务对象指针, 错误)
 */
func CreateAction(name string, method string) (*Action, error)
```

<span id="GetAction">GetAction</span>
```go
/* 功能 : 取当前协程对应的事务对象
 * 参数 : 无
 * 返回值 :
 *      事务组件对象指针
 */
GetAction() *Action
```
<span id="Action_CreateExternalComponent">(*Action).CreateExternalComponent</span>
```go
/* 功能 : 创建外部调用组件
 * 参数 : 
 *         url : 外部服务调用url
 *         method : 外部调用过程识别名
 * 返回值 :
 *      外部调用组件指针
 */
func (*Action) CreateExternalComponent(url string, method string) *Component
```

<span id="Action_CreateMQComponent">(*Action).CreateMQComponent</span>
```go
/*
 * 功能 : 创建消息队列组件
 * 参数 :
 *     vender : 消息队列类型(rabbitMQ/Kafka/ActiveMQ)
 *     isConsumer : 是否消费者组件
 *     host : MQ地址
 *     queue : 消息队列名
 * 返回值 :
 *     消息队列组件指针
 */
func (*Action) CreateMQComponent(vender string, isConsumer bool, host, queue string) *Component
```

<span id="Action_CreateMongoComponent">(*Action).CreateMongoComponent</span>
```go
/*
 * 功能 :  创建MongoDB访问组件
 * 参数 :
 *       host : 服务器地址
 *       database : 库名
 *       collection : collection 名
 *       op : mongo访问操作
 *       method : 组件过程识别名
 * 返回值 :
 *       MongoDB访问组件指针
 */
func (*Action) CreateMongoComponent(host, database, collection, op, method string) *Component
```

<span id="Action_CreateComponent">(*Action).CreateComponent</span>
```go
/*
 * 功能 : 创建自定义应用过程
 * 参数 :
 *        method : 组件过程识别名
 * 返回值 :
 *       自定义组件指针
 */
func (*Action) CreateComponent(method string) *Component
```

<span id="Action_AddRequestParam">(*Action).AddRequestParam</span>
```go
/*
 * 功能 : 采集事务请求参数
 * 参数 : 
 *        k : 参数名
 *        v : 参数值
 * 返回值 : 无
 */
func (*Action) AddRequestParam(k string, v string)
```

<span id="Action_AddResponseParam">(*Action).AddResponseParam</span>
```go
/*
 * 功能 : 采集事务响应参数
 * 参数 :
 *        k : 参数名
 *        v : 参数值
 * 返回值 : 无
 */
func (*Action) AddResponseParam(k string, v string)
```

<span id="Action_AddCustomParam">(*Action).AddCustomParam</span>
```go
/*
 * 功能 : 采集自定义参数
 * 参数 :
 *        k : 参数名
 *        v : 参数值
 * 返回值 : 无
 */
func (*Action) AddCustomParam(k string, v string)
```

<span id="Action_GetTxData">(*Action).GetTxData</span>
```go
/*
 * 功能 : 取事务执行性能数据 (跨应用追踪: 被调用端执行)
 * 参数 : 无
 * 返回值 : 
 *         事务性能数据
 */
func (*Action) GetTxData() string
```

<span id="Action_SetTrackID">(*Action).SetTrackID</span>
```go
/*
 * 功能 : 写入跨应用追踪数据 (跨应用追踪: 被调用端执行)
 * 参数 :
 *        id : 由调用端(*Component).CreateTrackID生成的, 调用过程携带到被调用端的字符串.
 * 返回值 : 无
 */
func (*Action) SetTrackID(id string)
```

<span id="Action_SetName">(*Action).SetName</span>
```go
/*
 * 功能 : 重新设置事务名
 * 参数 :
 *      name : 定位更准确的包名/类(结构)名
 *      method : 函数名
 * 返回值 : 无
 */
func (*Action) SetName(name string, method string)
```

<span id="Action_SetHTTPMethod">(*Action).SetHTTPMethod</span>
```go
/*
 * 功能 : 采集 http 请求方法
 * 参数 :
 *        httpMethod: GET/POST/HEAD/OPTIONS/PUT
 * 返回值 : 无
 */
func (*Action) SetHTTPMethod(httpMethod string)
```

<span id="Action_SetURL">(*Action).SetURL</span>
```go
/*
 * 功能 : 采集事务 URI
 * 参数 :
 *          name : 请求的URI
 * 返回值 : 无
 */
func (*Action) SetURL(name string)
```

<span id="Action_Ignore">(*Action).Ignore</span>
```go
/*
 * 功能 : 放弃本次采集的事务性能数据
 * 参数 :  无
 * 返回值 : 无
 */
func (*Action) Ignore()
```

<span id="Action_SetError">(*Action).SetError</span>
```go
/*
 * 功能 : 采集错误数据, 抓取调用栈
 * 参数 : 
 *        e : error 对象
 * 返回值 : 无
 */
func (*Action) SetError(e interface{})
```

<span id="Action_Finish">(*Action).Finish</span>
```go
/*
 * 功能 : 事务数据采集结束
 * 参数 :  无
 * 返回值 : 无
 */
func (*Action) Finish()
```

<span id="Action_SetStatusCode">(*Action).SetStatusCode</span>
```go
/*
 * 功能 : 采集事务状态码
 * 参数 : 
 *      code : 事务应答状态码
 * 返回值 : 
 *        整数, 成功为0, 失败为非0
 */
func (*Action) SetStatusCode(code uint16) int
```

<span id="Component_GetAction">(*Component).GetAction</span>
```go
/*
 * 功能 : 取组件关联的事务对象
 * 参数 : 无
 * 返回值 :
 *        事务对象
 */
func (*Component) GetAction() *Action
```

<span id="Component_SetError">(*Component).SetError</span>
```go
/*
 * 功能 : 采集组件错误
 * 参数 :
 *       e : error 对象
 *       errorType : 错误相关类型
 *       skipStack : 跳过采集的调用栈个数, 通常写0
 * 返回值 : 无
 */
func (*Component) SetError(e interface{}, errType string, skipStack int)
```

<span id="Component_Finish">(*Component).Finish</span>
```go
/*
 * 功能 : 组件性能数据采集结束
 * 参数 : 无
 * 返回值 : 无
 */
func (*Component) Finish()
```

<span id="Component_CreateTrackID">(*Component).CreateTrackID</span>
```go
/*
 * 功能 : 创建跨应用追踪ID字符串 (跨应用追踪: 调用端执行)
 * 参数 : 无
 * 返回值 :
 *        跨应用追踪ID信息(通过调用请求发送到被调用端)
 */
func (*Component) CreateTrackID() string
```

<span id="Component_SetTxData">(*Component).SetTxData</span>
```go
/*
* 功能 : 写入被调用端事务执行性能数据 (跨应用追踪: 调用端执行). 将rpc调用返回时携带的事务数据写入 
 * 参数 :  
 *        txData : 由被调用端 (*Action).GetTxData 生成的数据
 * 返回值 : 无
 */
func (*Component) SetTxData(txData string)
```

<span id="Component_CreateComponent">(*Component).CreateComponent</span>
```go
/*
 * 功能 : 组件再细分子组件
 * 参数 :
 *       子组件方法名
 * 返回值 :
 *       组件对象
 */
func (*Component) CreateComponent(method string) *Component
```

### API使用范例
  
  API手动嵌码请参考范例程序: [api_general](https://github.com/TingYunGo/goagent_examples/tree/main/api_general) 和 [api_webapp](https://github.com/TingYunGo/goagent_examples/tree/main/api_webapp)


## Code License
听云 goagent 使用 [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) 协议发布.

