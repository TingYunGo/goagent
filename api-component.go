// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"net/url"
	"time"
)

// Component : 构成事务的组件过程
type Component struct {
	action         *Action
	parent         *Component
	method         string
	classname      string
	txdata         string
	protocol       string
	vender         string
	instance       string
	table          string
	op             string
	callStack      []string
	errors         []*errInfo
	time           timeRange
	tracerID       int32
	tracerParentID int32
	exID           bool
	_type          uint8
}

func (c *Component) getVender() string {
	if c == nil {
		return ""
	}
	if len(c.vender) > 0 {
		return c.vender
	}
	switch c._type {
	case ComponentDefaultDB:
		return "DB"
	case ComponentMysql:
		return "MySQL"
	case ComponentPostgreSQL:
		return "PostgreSQL"
	case ComponentMSSQL:
		return "MSSQL"
	case ComponentSQLite:
		return "SQLite"
	case ComponentOracle:
		return "Oracle"
	case ComponentMongo:
		return "MongoDB"
	case ComponentRedis:
		return "Redis"
	case ComponentMQC:
		return "MQC"
	case ComponentMQP:
		return "MQP"
	case ComponentMemCache:
		return "Memcache"
	case ComponentExternal:
	case ComponentDefault:
	default:
	}
	return ""
}

// GetAction : 取对应的事务对象
func (c *Component) GetAction() *Action {
	if c == nil {
		return nil
	}
	return c.action
}
func (c *Component) setError(e interface{}, errType string, isError bool) {
	if c == nil {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	c.errors = append(c.errors, &errInfo{errTime, patchSize(toString(e), 4000), callStack(1), errType, isError})
}

// SetError : 组件错误捕获
func (c *Component) SetError(e interface{}, errType string, skipStack int) {
	if skipStack < 0 {
		skipStack = 0
	}
	if c == nil {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	c.errors = append(c.errors, &errInfo{errTime, patchSize(toString(e), 4000), callStack(skipStack + 1), errType, true})
}

// SetException : 组件异常捕获
func (c *Component) SetException(e interface{}, errType string, skipStack int) {
	if skipStack < 0 {
		skipStack = 0
	}
	if c == nil {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	c.errors = append(c.errors, &errInfo{errTime, patchSize(toString(e), 4000), callStack(skipStack + 1), errType, false})
}

// Finish : 停止组件计时
// 性能分解组件时长 = Finish时刻 - CreateComponent时刻
// 当时长超出堆栈阈值时，记录当前组件的代码堆栈
//go:noinline
func (c *Component) Finish() {
	c.End(1)
}
func (c *Component) SetMethod(method string) {
	if c != nil {
		c.method = method
	}
}

//End : 内部使用, skip为跳过的调用栈数
//go:noinline
func (c *Component) End(skip int) {
	c.FixStackEnd(skip+1, func(a string) bool {
		return false
	})
}

//End : 内部使用, skip为跳过的调用栈数
func (c *Component) FixStackEnd(skip int, checkRemovedFunction func(string) bool) {
	if c != nil {
		if c.action == nil {
			return
		}
		if c.parent == nil && c != c.action.root {
			c.parent = c.action.root
		}
		c.time.End()
		if c._type == ComponentDefault && c.action.current == c && c != c.action.root {
			c.action.current = c.parent
		}
		if len(c.callStack) > 0 {
			return
		}
		if readServerConfigBool(configServerConfigBoolActionTracerStackTraceEnabled, true) {
			//超阈值取callstack
			if c.time.duration >= time.Duration(readServerConfigInt(configServerConfigIntegerActionTracerStacktraceThreshold, 500))*time.Millisecond {
				c.callStack = validCallStack(skip+1, checkRemovedFunction)
			}
		}
	}
}

//CreateTrackID : 跨应用追踪接口,用于调用端,生成一个跨应用追踪id,通过http头或者私有协议发送到被调用端
//
//返回值: 字符串,一个包含授权id,应用id,实例id,事务id等信息的追踪id
func (c *Component) CreateTrackID() string {
	if app == nil || c == nil || c.action == nil || c._type != ComponentExternal {
		return ""
	}
	if !readLocalConfigBool(configLocalBoolTransactionEnabled, true) {
		return ""
	}
	if !readServerConfigBool(configServerConfigBoolActionTracerEnabled, false) {
		return ""
	}
	//c=CALL_LIST;x=TRANSACTION_TRACE_GUID;e=EXTERNAL_TRACE_GUID;n=TRANSACTION_NAME_MD5;
	if secID := app.configs.server.CStrings.Read(configServerStringTingyunIDSecret, ""); len(secID) != 0 {
		c.exID = true
		callList, transactionID := c.action.parseTrackID()
		if len(callList) > 0 {
			callList = "," + callList
		}
		if len(transactionID) == 0 {
			transactionID = c.action.unicID()
		}
		return "c=S|" + secID + callList + ";x=" + transactionID + ";e=" + c.unicID() + ";n=" + md5sum(c.action.GetName()) + ";"
	}
	return ""
}

//SetTxData : 跨应用追踪接口,用于调用端,将被调用端返回的事务性能数据保存到外部调用组件
//
//参数: 被调用端返回的事务的性能数据
func (c *Component) SetTxData(txData string) {
	if app == nil || c == nil || c.action == nil || c._type != ComponentExternal {
		return
	}
	if !readLocalConfigBool(configLocalBoolTransactionEnabled, true) {
		return
	}
	if readServerConfigBool(configServerConfigBoolActionTracerEnabled, false) {
		c.txdata = txData
	}
}

//AppendSQL : 用于数据库组件,通过此接口将sql查询语句保存到数据库组件,在报表慢事务追踪列表展示
//
//参数: sql语句
func (c *Component) AppendSQL(sql string) {
	if app == nil || c == nil || c.action == nil ||
		(c._type != ComponentExternal && c._type != ComponentDefaultDB && c._type != ComponentMysql && c._type != ComponentPostgreSQL && c._type != ComponentMSSQL && c._type != ComponentSQLite) {
		return
	}
	c.op = fixSQL(sql)
}

// CreateComponent : 在函数/方法中调用其他函数/方法时,如果认为有必要,调用此方法测量子过程性能
func (c *Component) CreateComponent(method string) *Component {
	if c == nil || c.action == nil || c._type != ComponentDefault {
		return nil
	}
	if !c.action.checkComponent() {
		return nil
	}
	r := &Component{
		action:         c.action,
		parent:         c.action.current,
		method:         url.QueryEscape(method),
		callStack:      nil,
		tracerParentID: c.tracerID,
		tracerID:       c.action.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		exID:           false,
		_type:          ComponentDefault,
	}
	c.action.current = c
	c.action.cache.Put(r)
	return r
}

func (c *Component) destroy() {
	if c._type == componentUnused {
		return
	}
	c.method = ""
	c.parent = nil
	c.classname = ""
	c.txdata = ""
	c.protocol = ""
	c.vender = ""
	c.instance = ""
	c.table = ""
	c.op = ""
	c.callStack = nil
	c.errors = nil
	c.action = nil
	c._type = componentUnused
}

//FixBegin : 校正事务开始时间
func (c *Component) FixBegin(begin time.Time) {
	c.time.begin = begin
}

func (c *Component) unicID() string {
	if c.exID {
		return unicID(c.time.begin, c)
	}
	return ""
}
