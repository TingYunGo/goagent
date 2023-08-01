// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
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
	ComponentOracle     = 37
	ComponentMongo      = 48
	ComponentMemCache   = 49
	ComponentRedis      = 50
	ComponentMQC        = 56
	ComponentMQP        = 57
	ComponentExternal   = 64
	componentUnused     = 255
)

var dbNameMap = [32]string{0: "Database", 1: "MySQL", 2: "PostgreSql", 3: "MSSQL", 4: "SQLite", 5: "Oracle", 16: "MongoDB", 17: "MemCache", 18: "Redis"}

const (
	actionUsing    = 1
	actionFinished = 2
	actionUnused   = 0
)

//Action : 事务对象
type Action struct {
	name           string
	category       string
	url            string
	path           string
	method         string
	httpMethod     string
	clientIP       string
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
	callbacks      pool.SerialReadPool
	tracerIDMaker  int32
	statusCode     uint16
	stateUsed      uint8
	trackEnable    bool
	isTask         bool
	enabledBack    bool
}

func (a *Action) IsTask() bool {
	if a == nil {
		return false
	}
	return a.isTask
}
func (a *Action) checkComponent() bool {
	if a == nil {
		return false
	}
	maxCacheSize := int32(app.configs.local.CIntegers.Read(configLocalIntegerComponentMax, 3000))
	return a.cache.Size() < maxCacheSize
}
func fixSQL(sql string) string {
	maxSize := int(app.configs.local.CIntegers.Read(configLocalIntegerMaxSQLSize, 5000))
	if len(sql) > maxSize {
		sql = sql[0:maxSize]
	}
	return sql
}

func parseURL(url string) (string, string) {
	protocol := ""
	if tystring.SubString(url, 0, 8) == "https://" {
		protocol = "https"
	} else if tystring.SubString(url, 0, 7) == "http://" {
		protocol = "http"
	} else if tystring.SubString(url, 0, 7) == "grpc://" {
		protocol = "grpc"
	}
	urlRequest := parseUriRequest(parseURI(url))
	if len(urlRequest) == 0 {
		urlRequest = "/"
	} else if len(urlRequest) > 1 && urlRequest[0] == '/' {
		urlRequest = urlRequest[1:]
	}
	return protocol, urlRequest
}

//CreateExternalComponent : 创建Web Service性能分解组件
//参数:
//    url    : 调用Web Service的url,格式: http(s)://host/uri, 例如 http://www.baidu.com/
//    method : 发起这个Web Service调用的类名.方法名, 例如 http.Get
func (a *Action) CreateExternalComponent(url string, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
		return nil
	}
	c := &Component{
		action:         a,
		instance:       url,
		method:         method,
		protocol:       "",
		op:             "",
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		exID:           false,
		statusCode:     0,
		callStack:      nil,
		time:           timeRange{time.Now(), -1},
		_type:          ComponentExternal,
	}
	c.SetURL(url)
	a.cache.Put(c)
	return c
}

// OnEnd : 注册一个在事务结束时执行的回调函数
func (a *Action) OnEnd(cb func()) {
	if a != nil && a.stateUsed == actionUsing {
		a.callbacks.Put(cb)
	}
}

// CreateMQComponent : 创建一个消息队列组件
//   vender : mq类型: kafka/rabbit MQ/ActiveMQ
func (a *Action) CreateMQComponent(vender string, isConsumer bool, host, queue string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
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
		statusCode:     0,
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

func getRedisInstanceName(host, key string) string {
	if len(key) == 0 {
		return host
	}
	redisUseKey := int(readLocalConfigInteger(configLocalIntegerRedisInstanceUseKey, 0))
	if key[0] == '[' {
		if (redisUseKey & 2) == 0 {
			return host
		}
	} else {
		if (redisUseKey & 1) == 0 {
			return host
		}
	}
	return host + "/" + key
}

// CreateRedisComponent : 创建一个Redis数据库访问组件
func (a *Action) CreateRedisComponent(host, cmd, key, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
		return nil
	}
	key = getValidString(key, "NULL")

	c := &Component{
		action:         a,
		method:         method,
		instance:       getRedisInstanceName(getValidString(host, "NULL"), key),
		table:          key,
		op:             cmd,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		exID:           false,
		callStack:      nil,
		statusCode:     0,
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
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
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
		op:             fixSQL(op),
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		statusCode:     0,
		_type:          dbType,
	}
	a.cache.Put(c)
	return c
}

// CreateMongoComponent 创建 Mongo 组件
func (a *Action) CreateMongoComponent(host, database, collection, op, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
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
		statusCode:     0,
		_type:          ComponentMongo,
	}
	a.cache.Put(c)
	return c
}

//CreateSQLComponent : 以 SQL语句创建一个数据库组件
func (a *Action) CreateSQLComponent(dbType uint8, host string, dbname string, sql string, method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
		return nil
	}
	nameID := dbType - ComponentDefaultDB
	if nameID < 0 || nameID >= 32 {
		return nil
	}
	c := &Component{
		action:         a,
		method:         method,
		op:             fixSQL(sql),
		vender:         getValidString(dbNameMap[nameID], "UnDefDatabase"),
		instance:       getValidString(host, "NULL") + "/" + getValidString(dbname, "NULL"),
		callStack:      nil,
		tracerParentID: a.current.tracerID,
		tracerID:       a.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		statusCode:     0,
		_type:          dbType,
	}
	a.cache.Put(c)
	return c
}

//CreateComponent : 创建函数/方法类型的组件
//参数
//    method : 类名.方法名, 例如 main.user.login
func (a *Action) CreateComponent(method string) *Component {
	if app == nil || a == nil || a.stateUsed != actionUsing || !app.inited {
		return nil
	}
	if !a.checkComponent() {
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
		statusCode:     0,
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
	if !readLocalConfigBool(configLocalBoolTransactionEnabled, true) {
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
	if !readLocalConfigBool(configLocalBoolTransactionEnabled, true) {
		return
	}
	if readServerConfigBool(configServerConfigBoolActionTracerEnabled, false) {
		a.trackID = id
	}
}

func (a *Action) SetBackEnabled(enabled bool) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.enabledBack = enabled
}

func formatActionName(instance string, method string, prefix string) string {
	if len(instance) == 0 {
		classEnd := strings.Index(method, ".")
		if classEnd != -1 {
			instance = method[0:classEnd]
			method = method[classEnd+1:]
		}
	}
	// mlen := len(method)
	// if mlen > 1 && method[0:1] == "/" {
	// 	method = method[1:mlen]
	// }

	if len(method) > 0 && method[0:1] == "/" {
		method = method[1:]
	}
	return prefix + instance + "/" + method
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
	if name == "CLIENTIP" {
		a.clientIP = method
		return
	} else if name == "URI" {
		a.path = method
	} else if name == "CLASS" {
		a.root.classname = method
	} else {
		a.method = method
		a.root.method = method
	}
	if name == "URI" {
		if len(a.category) == 0 || !readServerConfigBool(ServerConfigBoolAutoActionNaming, true) {
			a.category = name
		}
	} else {
		a.category = name
	}
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
	if len(a.name) > 0 {
		return a.prefixName() + a.name
	}
	path := a.path
	category := a.category
	if a.category == "CLASS" {
		if len(a.method) > 0 {
			category = a.root.classname
			path = a.method
		} else {
			path = a.root.classname
		}
	} else if a.category != "URI" {
		path = a.method
	}
	if path == a.path {
		path = formatRestfulURI(path, int(readLocalConfigInteger(ConfigLocalIntegerRestfulUUIDMinSize, 8)))
	}
	return formatActionName(category, path, a.prefixName())
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

func (a *Action) Duration() time.Duration {
	if a == nil {
		return 0
	}
	if a.stateUsed == actionUnused {
		return 0
	}
	if a.stateUsed == actionFinished {
		return a.time.duration
	}
	if a.stateUsed == actionUsing {
		return time.Now().Sub(a.time.begin)
	}
	return 0
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
	a.setError(e, "ActionError", 1, true)
}

//SetException : 事务发生错误或异常时调用,记录事务的运行时错误信息
func (a *Action) SetException(e interface{}) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.setError(e, "ActionException", 1, false)
}

//Finish : 事务结束时调用
//HTTP请求时长 = Finish时刻 - CreateAction时刻
func (a *Action) Finish() {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	for callback := a.callbacks.Get(); callback != nil; callback = a.callbacks.Get() {
		cb := callback.(func())
		cb()
	}

	a.stateUsed = actionFinished
	if a.statusCode == 0 {
		a.statusCode = 200
	}
	if traceDisabled {
		a.destroy()
		return
	}
	appendAction(a)
}

var traceDisabled bool = false

func init() {
	//丢弃trace
	//export TINGYUN_GO_DEBUG=TRACE_DISABLED:OTHER
	splitStrings(os.Getenv("TINGYUN_GO_DEBUG"), func(item string) bool {
		if item != "TRACE_DISABLED" {
			return false
		}
		traceDisabled = true
		return true
	}, func(t byte) bool { return t == ':' })
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
		a.setError(errors.New(fmt.Sprint("status code ", code)), "HTTP_ERROR", 1, true)
		return 3
	}
	return 0
}

//FixBegin : 内部使用, 重置事务开始时间
func (a *Action) FixBegin(begin time.Time) {
	if a == nil {
		return
	}
	a.time.begin = begin
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
		a.setError(errors.New(fmt.Sprint("status code ", code)), "HTTP_ERROR", skip+1, true)
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

func Enabled() bool {
	if app == nil {
		return false
	}
	if configDisabled {
		return false
	}
	if !readLocalConfigBool(configLocalBoolAgentEnable, true) {
		return false
	}
	return readServerConfigBool(configServerConfigBoolAgentEnabled, true)
}

//CreateAction : 在方法method中调用并 创建一个名为 name的事务,
func CreateAction(name string, method string) (*Action, error) {
	if app == nil || !app.inited || app.serverCtrl.login_time == 0 {
		if configDisabled {
			return nil, errors.New("Agent disabled by local config file")
		}
		return nil, errors.New("Agent not Inited")
	} else if app.actionPool.Size() > int32(app.configs.local.CIntegers.Read(configLocalIntegerNbsActionCacheMax, 10000)) {
		return nil, errors.New("Server busy, Skip one action")
	}
	return app.createAction(name, method, false)
}
func CreateTask(method string) (*Action, error) {
	if app == nil || !app.inited || app.serverCtrl.login_time == 0 {
		if configDisabled {
			return nil, errors.New("Agent disabled by local config file")
		}
		return nil, errors.New("Agent not Inited")
	} else if app.actionPool.Size() > int32(app.configs.local.CIntegers.Read(configLocalIntegerNbsActionCacheMax, 10000)) {
		return nil, errors.New("Server busy, Skip one action")
	}
	a, e := app.createAction("Job", method, true)
	return a, e
}

func (a *Action) destroy() {
	if a == nil || a.stateUsed == actionUnused {
		return
	}
	a.name = ""
	a.category = ""
	a.url = ""
	a.path = ""
	a.method = ""
	a.httpMethod = ""
	a.clientIP = ""
	a.trackID = ""
	a.actionID = ""
	for component := a.cache.Get(); component != nil; component = a.cache.Get() {
		component.(*Component).destroy()
	}
	for a.errors.Get() != nil {
	}
	a.root = nil
	a.requestParams = nil
	a.customParams = nil
	a.stateUsed = actionUnused
	a.statusCode = 0
}
func getAction() *Action {
	var res *Action = nil
	routineLocalExec(func(local interface{}) interface{} {
		if local != nil {
			localStorage := local.(*RoutineLocal)
			res = localStorage.action
		}
		return local
	})
	return res
}
func setAction(action *Action) {
	if action == nil {
		return
	}
	routineLocalExec(func(local interface{}) interface{} {
		if local == nil {
			return &RoutineLocal{action, nil, nil}
		} else {
			localStorage := local.(*RoutineLocal)
			localStorage.action = action
			return local
		}
	})
}

func getComponent() *Component {
	var res *Component = nil
	routineLocalExec(func(local interface{}) interface{} {
		if local != nil {
			localStorage := local.(*RoutineLocal)
			res = localStorage.component
		}
		return local
	})
	return res
}
func setComponent(component *Component) {
	routineLocalExec(func(local interface{}) interface{} {
		if local != nil {
			localStorage := local.(*RoutineLocal)
			localStorage.component = component
		}
		return local
	})
}

// LocalGet : 从协程局部存储器中取出key为 id的对象,没有则返回nil
func LocalGet(id int) interface{} {
	var res interface{} = nil
	routineLocalExec(func(local interface{}) interface{} {
		if local == nil {
			return nil
		}
		localStorage := local.(*RoutineLocal)
		if localStorage == nil {
			return nil
		}
		if localStorage.pointers != nil {
			if r, found := localStorage.pointers[id]; found {
				res = r
			}
		}
		return local
	})
	return res
}

// LocalSet : 以id为key, 将对象object写入协程局部存储器
func LocalSet(id int, object interface{}) {
	routineLocalExec(func(local interface{}) interface{} {
		if local == nil {
			return &RoutineLocal{nil, nil, map[int]interface{}{
				id: object,
			}}
		} else {
			localStorage := local.(*RoutineLocal)
			if localStorage.pointers == nil {
				localStorage.pointers = map[int]interface{}{
					id: object,
				}
			} else {
				localStorage.pointers[id] = object
			}
			return local
		}
	})
}

// LocalDelete : 从协程局部存储器中删除key为 id的对象,并返回这个对象
func LocalDelete(id int) interface{} {
	var res interface{} = nil
	routineLocalExec(func(local interface{}) interface{} {
		if local == nil {
			return nil
		}
		localStorage := local.(*RoutineLocal)
		if localStorage.pointers != nil {
			if r, found := localStorage.pointers[id]; found {
				res = r
				delete(localStorage.pointers, id)
				if localStorage.action == nil && localStorage.component == nil && len(localStorage.pointers) == 0 {
					return nil
				}
				return local
			}
		}
		return local
	})
	return res
}
func (a *Action) GetTransactionID() string {
	if a == nil {
		return ""
	}
	return a.getTransactionID()
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

func FindAction(ctx context.Context) (action *Action, sync bool) {
	if action := getAction(); action != nil {
		return action, true
	} else if ctx == nil {
		return nil, true
	} else if value := ctx.Value("TingYunWebAction"); value != nil {
		if a, ok := value.(*Action); ok {
			return a, false
		}
	}
	return nil, true
}

// SetAction : 辅助功能函数: 将事务存到协程局部存储器
func SetAction(action *Action) {
	setAction(action)
}
