// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/TingYunGo/goagent/libs/pool"
	"github.com/TingYunGo/goagent/libs/tystring"
)

/*
 * 组件类型定义
 */
const (
	ComponentDefault    = 0
	ComponentDefaultDB  = 32
	ComponentMysql      = 33
	ComponentPostgreSQL = 34
	ComponentMSSQL      = 35
	ComponentSQLite     = 36
	ComponentMongo      = 48
	ComponentMemCache   = 49
	ComponentRedis      = 50
	ComponentMQC        = 56
	ComponentMQP        = 57
	ComponentExternal   = 64
	componentUnused     = 255
)

var dbNameMap = [32]string{0: "Database", 1: "MySQL", 2: "PostgreSql", 3: "MSSQL", 4: "SQLite", 16: "MongoDB", 17: "MemCache", 18: "Redis"}

const (
	actionUsing    = 1
	actionFinished = 2
	actionUnused   = 0
)

//Action : 事务对象
type Action struct {
	category       string
	url            string
	path           string
	method         string
	httpMethod     string
	trackID        string
	actionID       string
	cache          pool.SerialReadPool
	errors         pool.Pool
	requestParams  map[string]string
	responseParams map[string]string
	customParams   map[string]string
	time           timeRange
	root           *Component
	current        *Component
	callbacks      []interface{}
	tracerIDMaker  int32
	statusCode     uint16
	stateUsed      uint8
	trackEnable    bool
}

//CreateExternalComponent : 创建Web Service性能分解组件
//参数:
//    url    : 调用Web Service的url,格式: http(s)://host/uri, 例如 http://www.baidu.com/
//    method : 发起这个Web Service调用的类名.方法名, 例如 http.Get
func (a *Action) CreateExternalComponent(url string, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	protocol := ""
	if tystring.SubString(url, 0, 8) == "https://" {
		protocol = "https"
	} else if tystring.SubString(url, 0, 7) == "http://" {
		protocol = "http"
	}
	c := &Component{
		action:         a,
		instance:       url,
		method:         method,
		protocol:       protocol,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		exID:           false,
		callStack:      nil,
		time:           timeRange{time.Now(), -1},
		_type:          ComponentExternal,
	}
	a.cache.Put(c)
	return c
}

// OnEnd : 注册一个在事务结束时执行的回调函数
func (a *Action) OnEnd(cb func()) {
	if a.callbacks == nil {
		a.callbacks = []interface{}{cb}
	} else {
		a.callbacks = append(a.callbacks, cb)
	}
}

// CreateMQComponent : 创建一个消息队列组件
//   vender : mq类型: kafka/rabbit MQ/ActiveMQ
func (a *Action) CreateMQComponent(vender string, isConsumer bool, host, queue string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	var mqType uint8 = ComponentMQP
	if isConsumer {
		mqType = ComponentMQC
	}

	c := &Component{
		action:         a,
		method:         GetCallerName(1),
		protocol:       "",
		vender:         vender,
		instance:       getValidString(host, "NULL") + "/" + getValidString(queue, "NULL"),
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		exID:           false,
		callStack:      nil,
		time:           timeRange{time.Now(), -1},
		_type:          mqType,
	}
	a.cache.Put(c)
	return c
}

func getValidString(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// CreateRedisComponent : 创建一个Redis数据库访问组件
func (a *Action) CreateRedisComponent(host, cmd, key, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	key = getValidString(key, "NULL")
	c := &Component{
		action:         a,
		method:         method,
		instance:       getValidString(host, "NULL") + "/" + key,
		table:          key,
		op:             cmd,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		exID:           false,
		callStack:      nil,
		time:           timeRange{time.Now(), -1},
		_type:          ComponentRedis,
	}
	a.cache.Put(c)
	return c
}

//CreateDBComponent 创建数据库或NOSQL性能分解组件
//参数:
//    dbType : 组件类型 (ComponentMysql, ComponentPostgreSQL, ComponentMongo, ComponentMemCache, ComponentRedis)
//    host   : 主机地址，可空
//    dbname : 数据库名称，可空
//    table  : 数据库表名
//    op     : 操作类型, 关系型数据库("SELECT", "INSERT", "UPDATE", "DELETE" ...), NOSQL("GET", "SET" ...)
//    method : 发起这个数据库调用的类名.方法名, 例如 db.query redis.get
func (a *Action) CreateDBComponent(dbType uint8, host string, dbname string, table string, op string, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	nameID := dbType - ComponentDefaultDB
	if nameID < 0 || nameID >= 32 {
		return nil
	}
	c := &Component{
		action:         a,
		method:         method,
		vender:         getValidString(dbNameMap[nameID], "UnDefDatabase"),
		instance:       getValidString(host, "NULL") + "/" + getValidString(dbname, "NULL"),
		table:          getValidString(table, "NULL"),
		op:             op,
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		_type:          dbType,
	}
	a.cache.Put(c)
	return c
}

// CreateMongoComponent 创建 Mongo 组件
func (a *Action) CreateMongoComponent(host, database, collection, op, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	c := &Component{
		action:         a,
		method:         method,
		vender:         getValidString(dbNameMap[ComponentMongo-ComponentDefaultDB], "MongoDB"),
		instance:       getValidString(host, "NULL") + "/" + getValidString(database, "NULL"),
		table:          getValidString(collection, "NULL"),
		op:             op,
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		_type:          ComponentMongo,
	}
	a.cache.Put(c)
	return c
}

//CreateSQLComponent : 以 SQL语句创建一个数据库组件
func (a *Action) CreateSQLComponent(dbType uint8, host string, dbname string, sql string, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	nameID := dbType - ComponentDefaultDB
	if nameID < 0 || nameID >= 32 {
		return nil
	}
	c := &Component{
		action:         a,
		method:         method,
		op:             sql,
		vender:         getValidString(dbNameMap[nameID], "UnDefDatabase"),
		instance:       getValidString(host, "NULL") + "/" + getValidString(dbname, "NULL"),
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		_type:          dbType,
	}
	a.cache.Put(c)
	return c
}

//CreateComponent : 创建函数/方法类型的组件
//参数
//    method : 类名.方法名, 例如 main.user.login
func (a *Action) CreateComponent(method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing {
		return nil
	}
	c := &Component{
		action:         a,
		method:         url.QueryEscape(method),
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		_type:          ComponentDefault,
	}
	c.parent = a.current
	a.current = c
	a.cache.Put(c)
	return c
}

//AddRequestParam : 添加请求参数
func (a *Action) AddRequestParam(k string, v string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.requestParams[k] = v
}

//AddResponseParam : 添加应答参数
func (a *Action) AddResponseParam(k string, v string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.responseParams[k] = v
}

//AddCustomParam : 添加自定义参数
func (a *Action) AddCustomParam(k string, v string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.customParams[k] = v
}

//GetTxData : 跨应用追踪接口,用于被调用端,获取当前事务的执行性能信息,通过http头或者自定义协议传回调用端
//
//返回值: 事务的性能数据
func (a *Action) GetTxData() string {
	if a == nil || a.stateUsed != actionUsing || len(a.trackID) == 0 {
		return ""
	}
	if !readServerConfigBool(configServerConfigBoolActionTracerEnabled, false) {
		return ""
	}
	secID := app.configs.server.CStrings.Read(configServerStringTingyunIDSecret, "")
	if len(secID) == 0 {
		return ""
	}

	currTime := time.Now()
	currentDuration := a.time.GetDuration(currTime)
	res := map[string]interface{}{
		"id":       secID,
		"tname":    a.GetName(),
		"tid":      a.getTransactionID(),
		"rid":      a.unicID(),
		"duration": currentDuration / time.Millisecond,
	}

	if jsonByte, err := json.Marshal(res); err == nil {
		return string(jsonByte)
	}
	return ""
}

//SetTrackID : 跨应用追踪接口,用于被调用端,保存调用端传递过来的跨应用追踪id
//
//参数: 跨应用追踪id
func (a *Action) SetTrackID(id string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	if readServerConfigBool(configServerConfigBoolActionTracerEnabled, false) {
		a.trackID = id
	}
}
func formatActionName(instance string, method string, isTransaction bool) string {
	if len(instance) == 0 {
		classEnd := strings.LastIndex(method, ".")
		if classEnd != -1 {
			instance = method[0:classEnd]
			method = method[classEnd+1:]
		}
	}
	// mlen := len(method)
	// if mlen > 1 && method[0:1] == "/" {
	// 	method = method[1:mlen]
	// }
	preName := "WebAction/"
	if isTransaction {
		preName = "Transaction/"
	}
	return preName + instance + "/" + method
	//	return preName + url.QueryEscape(instance) + "/" + url.QueryEscape(method)
}

//SetName : 设置HTTP请求的友好名称
//参数:
//    instance   : 分类, 例如 loginController
//    method : 方法, 例如 POST
func (a *Action) SetName(name string, method string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	if name == "URI" {
		a.path = method
	} else if name == "CLASS" {
		a.root.classname = method
	} else {
		a.method = method
		a.root.method = method
	}
	a.category = name
}

// SetHTTPMethod : 设置 HTTP请求方法名
func (a *Action) SetHTTPMethod(httpMethod string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.httpMethod = httpMethod
}
func (a *Action) GetMethod() string {
	if a == nil || a.stateUsed != actionUsing {
		return ""
	}
	return a.httpMethod
}

// GetName : 取事务名字
func (a *Action) GetName() string {
	if a == nil {
		return ""
	}
	path := a.path
	if a.category != "URI" {
		path = a.method
	}
	return formatActionName(a.category, path, !strings.Contains(a.trackID, ";n="))
}

// GetURL : 取事务的 URL
func (a *Action) GetURL() string {
	if a == nil {
		return ""
	}
	return a.url
}

//SetURL : 设置事务的url
func (a *Action) SetURL(name string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.url = name
}

//Slow : 检测本次事务是否为慢请求
//返回值: 当HTTP请求性能超出阈值时为true, 否则为false
func (a *Action) Slow() bool {
	if a == nil {
		return false
	}
	if !readServerConfigBool(configServerConfigBoolActionTracerEnabled, true) {
		return false
	}
	if a.stateUsed == actionUnused {
		return false
	}
	threshold := readServerConfigInt(configServerConfigIntegerActionTracerActionThreshold, 500)
	if a.stateUsed == actionUsing {
		return time.Now().Sub(a.time.begin) >= time.Duration(threshold)*time.Millisecond
	} else if a.stateUsed == actionFinished {
		return a.time.duration >= time.Duration(threshold)*time.Millisecond
	}
	return false
}

//Ignore : 忽略本次事务的性能数据
func (a *Action) Ignore() {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.destroy()
}

//HasError : 是否发生过错误或异常
func (a *Action) HasError() bool {
	if a == nil || a.stateUsed != actionUsing {
		return false
	}
	return a.errors.Size() > 0
}

//SetError : 事务发生错误或异常时调用,记录事务的运行时错误信息
func (a *Action) SetError(e interface{}) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.setError(e, "RUNTIME_ERROR", 1)
}

//Finish : 事务结束时调用
//HTTP请求时长 = Finish时刻 - CreateAction时刻
func (a *Action) Finish() {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	if a.callbacks != nil {
		i := len(a.callbacks)
		for i > 0 {
			i--
			cb := a.callbacks[i].(func())
			cb()
			a.callbacks[i] = nil
		}
		a.callbacks = nil
	}
	a.stateUsed = actionFinished
	if a.statusCode == 0 {
		a.statusCode = 200
	}
	appendAction(a)
}

//SetStatusCode : 正常返回0
//无效的Action 返回1
//如果状态码是被忽略的错误,返回2
//错误的状态码,返回3
func (a *Action) SetStatusCode(code uint16) int {
	if a == nil || a.stateUsed != actionUsing {
		return 1
	}
	if a.statusCode == 0 {
		a.statusCode = code
	}
	if code >= 400 && code != 401 { //401认证失败，非错误代码
		ignoredStatusCodes := app.configs.serverArrays.Read(configServerConfigIArrayIgnoredStatusCodes, nil)
		if ignoredStatusCodes != nil {
			for _, v := range ignoredStatusCodes {
				if uint16(v) == code {
					return 2
				}
			}
		}
		a.setError(errors.New(fmt.Sprint("status code ", code)), "HTTP_ERROR", 1)
		return 3
	}
	return 0
}

//SetHTTPStatus : 内部使用: 添加 http状态, skip为跳过的调用栈
func (a *Action) SetHTTPStatus(code uint16, skip int) int {
	if a == nil || a.stateUsed != actionUsing {
		return 1
	}
	if a.statusCode == 0 {
		a.statusCode = code
	}
	if code >= 400 && code != 401 { //401认证失败，非错误代码
		ignoredStatusCodes := app.configs.serverArrays.Read(configServerConfigIArrayIgnoredStatusCodes, nil)
		if ignoredStatusCodes != nil {
			for _, v := range ignoredStatusCodes {
				if uint16(v) == code {
					return 2
				}
			}
		}
		a.setError(errors.New(fmt.Sprint("status code ", code)), "HTTP_ERROR", skip+1)
		return 3
	}
	return 0
}
func agentEnabled() bool {
	if app == nil {
		return false
	} else if configDisabled {
		return false
	} else {
		return readServerConfigBool(configServerConfigBoolAgentEnabled, true)
	}
}

//CreateAction : 在方法method中调用并 创建一个名为 name的事务,
func CreateAction(name string, method string) (*Action, error) {
	if app == nil {
		if configDisabled {
			return nil, errors.New("Agent disabled by local config file")
		}
		return nil, errors.New("Agent not Inited, please call AppInit() first")
	} else if app.actionPool.Size() > int32(app.configs.local.CIntegers.Read(configLocalIntegerNbsActionCacheMax, 10000)) {
		return nil, errors.New("Server busy, Skip one action")
	}
	return app.createAction(name, method)
}
func (a *Action) destroy() {
	if a == nil || a.stateUsed == actionUnused {
		return
	}
	a.category = ""
	a.url = ""
	a.path = ""
	a.method = ""
	a.httpMethod = ""
	a.trackID = ""
	a.actionID = ""
	for component := a.cache.Get(); component != nil; component = a.cache.Get() {
		component.(*Component).destroy()
	}
	for a.errors.Get() != nil {
	}
	a.callbacks = nil
	a.root = nil
	a.requestParams = nil
	a.customParams = nil
	a.stateUsed = actionUnused
	a.statusCode = 0
}
func getAction() *Action {
	if l := routineLocalGet(); l != nil {
		localStorage := l.(*RoutineLocal)
		return localStorage.action
	}
	return nil
}
func setAction(action *Action) {
	if action == nil {
		return
	}
	if l := routineLocalGet(); l == nil {
		localStorage := &RoutineLocal{action, nil, nil}
		routineLocalSet(localStorage)
		localStorage.action = action
	} else {
		localStorage := l.(*RoutineLocal)
		localStorage.action = action
	}
}

func getComponent() *Component {
	l := routineLocalGet()
	if l == nil {
		return nil
	}
	localStorage := l.(*RoutineLocal)
	if localStorage == nil {
		return nil
	}
	return localStorage.component
}
func setComponent(component *Component) {
	l := routineLocalGet()
	if l == nil {
		return
	}
	localStorage := l.(*RoutineLocal)
	if localStorage != nil {
		localStorage.component = component
	}
}

// LocalGet : 从协程局部存储器中取出key为 id的对象,没有则返回nil
func LocalGet(id int) interface{} {
	l := routineLocalGet()
	if l == nil {
		return nil
	}
	localStorage := l.(*RoutineLocal)
	if localStorage == nil {
		return nil
	}
	if localStorage.pointers == nil {
		return nil
	}
	if r, found := localStorage.pointers[id]; found {
		return r
	}
	return nil
}

// LocalSet : 以id为key, 将对象object写入协程局部存储器
func LocalSet(id int, object interface{}) {

	if l := routineLocalGet(); l == nil {
		localStorage := &RoutineLocal{nil, nil, map[int]interface{}{
			id: object,
		}}
		routineLocalSet(localStorage)
	} else {
		localStorage := l.(*RoutineLocal)
		if localStorage.pointers == nil {
			localStorage.pointers = map[int]interface{}{
				id: object,
			}
		} else {
			localStorage.pointers[id] = object
		}
	}
}

// LocalDelete : 从协程局部存储器中删除key为 id的对象,并返回这个对象
func LocalDelete(id int) interface{} {

	if l := routineLocalGet(); l != nil {
		localStorage := l.(*RoutineLocal)
		if localStorage.pointers != nil {
			if r, found := localStorage.pointers[id]; found {
				delete(localStorage.pointers, id)
				return r
			}
		}
	}
	return nil
}

// GetComponent : 辅助功能函数: 将存储到协程局部存储器的组件对象取出
func GetComponent() *Component {
	return getComponent()
}

// SetComponent : 辅助功能函数: 将组件存到协程局部存储器
func SetComponent(c *Component) {
	setComponent(c)
}

// GetAction : 辅助功能函数: 将存储到协程局部存储器的Action对象取出
func GetAction() *Action {
	return getAction()
}

// SetAction : 辅助功能函数: 将事务存到协程局部存储器
func SetAction(action *Action) {
	setAction(action)
}
