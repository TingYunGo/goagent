// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"
)

// Component : 构成事务的组件过程
type Component struct {
	action         *Action
	name           string
	method         string
	classname      string
	txdata         string
	extSecretID    string
	protocol       string
	vender         string
	host           string
	instance       string
	table          string
	op             string
	sql            string
	callStack      []string
	errors         []*errInfo
	time           timeRange
	aloneTime      time.Duration
	remoteDuration float64
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
func (c *Component) setError(e interface{}, errType string) {
	if c == nil {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	c.errors = append(c.errors, &errInfo{errTime, e, callStack(1), errType})
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
	c.errors = append(c.errors, &errInfo{errTime, e, callStack(skipStack + 1), errType})
}

// Finish : 停止组件计时
// 性能分解组件时长 = Finish时刻 - CreateComponent时刻
// 当时长超出堆栈阈值时，记录当前组件的代码堆栈
//go:noinline
func (c *Component) Finish() {
	c.End(1)
}

//End : 内部使用, skip为跳过的调用栈数
func (c *Component) End(skip int) {
	if c != nil {
		c.time.End()
		if c._type != ComponentDefault {
			c.aloneTime = c.time.duration
		}
		if len(c.callStack) > 0 {
			return
		}
		if readServerConfigBool(configServerConfigBoolActionTracerStackTraceEnabled, true) {
			//超阈值取callstack
			if c.time.duration >= time.Duration(readServerConfigInt(configServerConfigIntegerActionTracerStacktraceThreshold, 500))*time.Millisecond {
				c.callStack = callStack(skip + 1)
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
	if enabled := readServerConfigBool(configServerConfigBoolActionTracerEnabled, false); !enabled {
		return ""
	}
	//TINGYUN_ID_SECRET;c=CALLER_TYPE;r=REQ_ID;x=TX_ID;e=EXTERNAL_ID;p=PROTOCOL
	//时间+对象地址=>生成exId
	if secID := app.configs.server.CStrings.Read(configServerStringTingyunIDSecret, ""); len(secID) != 0 {
		c.exID = true
		// protocol := "http"
		// if arr := strings.Split(c.name, "://"); len(arr) > 1 {
		// 	protocol = arr[0]
		// }
		if len(c.action.trackID) > 0 && strings.Contains(c.action.trackID, ";n=") {
			preTrackIds := strings.Split(c.action.trackID, ";")
			return "c=S|" + secID + ";" + tystring.SubString(preTrackIds[0], 2, len(preTrackIds[0])-2) + ";" + preTrackIds[1] + ";e=" + c.unicID() + ";n=" + md5sum(c.action.GetName())
		}
		return "c=S|" + secID + ";x=" + c.action.unicID() + ";e=" + c.unicID() + ";n=" + md5sum(c.action.GetName())
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
	jsonData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(txData), &jsonData); err != nil {
		return
	}
	if tr, err := jsonReadInt(jsonData, "tr"); err == nil {
		c.action.trackEnable = (tr != 0)
	}
	c.txdata = txData
}

//AppendSQL : 用于数据库组件,通过此接口将sql查询语句保存到数据库组件,在报表慢事务追踪列表展示
//
//参数: sql语句
func (c *Component) AppendSQL(sql string) {
	if app == nil || c == nil || c.action == nil ||
		(c._type != ComponentExternal && c._type != ComponentDefaultDB && c._type != ComponentMysql && c._type != ComponentPostgreSQL && c._type != ComponentMSSQL && c._type != ComponentSQLite) {
		return
	}
	c.sql = sql
}

// CreateComponent : 在函数/方法中调用其他函数/方法时,如果认为有必要,调用此方法测量子过程性能
func (c *Component) CreateComponent(method string) *Component {
	if c == nil || c.action == nil || c._type != ComponentDefault {
		return nil
	}
	r := &Component{
		action:         c.action,
		name:           "",
		method:         url.QueryEscape(method),
		callStack:      nil,
		tracerParentID: c.tracerID,
		tracerID:       c.action.makeTracerID(),
		time:           timeRange{time.Now(), -1},
		aloneTime:      0,
		remoteDuration: 0,
		exID:           false,
		_type:          ComponentDefault,
	}
	c.action.cache.Put(r)
	return r
}
