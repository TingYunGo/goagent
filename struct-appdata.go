// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"fmt"
	"strings"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent/protoc"
	proto "github.com/golang/protobuf/proto"
)

type structAppData struct {
	traces protoc.ActionTraces
}

func mapTraceType(t uint8) protoc.TracerType {
	switch t {
	case ComponentExternal:
		return protoc.TracerType_External
	case ComponentDefaultDB:
		return protoc.TracerType_Database
	case ComponentMysql:
		return protoc.TracerType_Database
	case ComponentPostgreSQL:
		return protoc.TracerType_Database
	case ComponentMSSQL:
		return protoc.TracerType_Database
	case ComponentSQLite:
		return protoc.TracerType_Database
	case ComponentMemCache:
		return protoc.TracerType_Memcached
	case ComponentMongo:
		return protoc.TracerType_Mongo
	case ComponentRedis:
		return protoc.TracerType_Redis
	case ComponentMQC:
		return protoc.TracerType_MQC
	case ComponentMQP:
		return protoc.TracerType_MQP
	case ComponentDefault:
		return protoc.TracerType_Java
	default:
		return protoc.TracerType_Java
	}
}

//tingyun3: 2021
//追加action数据
func (r *structAppData) Append(action *Action) {

	trace := &protoc.ActionTrace{}
	trace.Action = action.GetName()
	trace.Time = action.time.begin.Unix()
	trace.Duration = (int64)(action.time.duration / time.Millisecond)
	trace.Rid = action.unicID()
	array := strings.Split(action.trackID, ";")
	for _, item := range array {
		switch tystring.SubString(item, 0, 2) {
		case "e=":
			trace.Refid = tystring.SubString(item, 2, len(item)-2)
		case "c=":
			trace.Cross = tystring.SubString(item, 2, len(item)-2)
		case "n=":
			trace.Tmd5 = tystring.SubString(item, 2, len(item)-2)
		case "x=":
			trace.Tid = tystring.SubString(item, 2, len(item)-2)
		default:
		}
	}
	if trace.Tmd5 == "" {
		trace.Tmd5 = md5sum(trace.Action)
	}
	if len(trace.Tid) == 0 {
		trace.Tid = trace.Rid
	}
	trace.Status = int32(action.statusCode)
	if methodType, found := protoc.HttpMethod_value[action.httpMethod]; found {
		trace.Method = (protoc.HttpMethod)(methodType)
	} else {
		trace.Method = protoc.HttpMethod_UNKNOWN
	}
	trace.Url = action.url
	trace.Ip = action.clientIP
	//trace.NoSample: 当前是否启用采样
	//
	detail := trace.Detail
	if detail == nil {
		trace.Detail = &protoc.TraceDetail{}
		detail = trace.Detail
	}
	detail.Custom = action.customParams
	detail.QueryStringParameters = parseQueryString(parseURI(action.url))
	detail.RequestHeader = action.requestParams
	detail.ResponseHeader = action.responseParams
	for action.errors.Size() > 0 {
		exception := action.errors.Get()
		action.root.errors = append(action.root.errors, exception.(*errInfo))
	}
	action.root.method = action.method
	action.root.time = action.time
	for action.cache.Size() > 0 {
		c := action.cache.Get()
		if c == nil {
			continue
		}
		component := c.(*Component)
		traceItem := &protoc.TracerItem{}
		traceItem.TracerId = component.tracerID
		traceItem.ParentTracerId = component.tracerParentID
		traceItem.Type = mapTraceType(component._type)
		beginOffset := component.time.begin.Sub(action.time.begin)
		traceItem.Start = int64(beginOffset / time.Millisecond)
		traceItem.End = int64((beginOffset + component.time.duration) / time.Millisecond)
		traceItem.Metric = ""
		traceItem.Clazz = component.classname
		traceItem.Method = component.method
		traceItem.Backtrace = component.callStack
		for i := 0; i < len(component.errors); i++ {
			exception := &protoc.TracerException{}
			exception.Error = component.errors[i].isError
			exception.Msg = fmt.Sprint(component.errors[i].e)
			exception.Name = component.errors[i].eType
			exception.Stack = component.errors[i].stack
			traceItem.Exception = append(traceItem.Exception, exception)
		}
		var params *protoc.TracerParams = nil
		if component._type != ComponentDefault {
			params = &protoc.TracerParams{}
			params.Operation = component.op
			params.Protocol = component.protocol
			params.Key = component.table
			params.Instance = component.instance
			params.Vendor = component.getVender()
			params.TxData = component.txdata
			if component._type == ComponentExternal {
				params.ExternalId = component.unicID()
				params.Instance = parseHost(component.instance)
			} else {
				params.Instance = component.instance
			}
			traceItem.Params = params
		}
		detail.Tracers = append(detail.Tracers, traceItem)
	}
	action.root = nil
	action.current = nil
	r.traces.Traces = append(r.traces.Traces, trace)
}

//数据序列化
func (r *structAppData) Serialize() ([]byte, error) {
	return proto.Marshal(&r.traces)
}

func (r *structAppData) destroy() {
	if len(r.traces.Traces) > 0 {
		for _, trace := range r.traces.Traces {
			trace.Reset()
		}
	}
	r.traces.Reset()
}
func (a *application) GetReportBlock(reportMax, saveCount int) *structAppData {
	if a.reportQueue.Size() == 0 {
		a.reportQueue.PushBack(&structAppData{})
	}
	data, _ := a.reportQueue.Back().Value()
	if datablock := data.(*structAppData); len(datablock.traces.Traces) >= reportMax {
		if a.reportQueue.Size() > saveCount {
			return nil
		}
		a.reportQueue.PushBack(&structAppData{})
		data, _ = a.reportQueue.Back().Value()
	}
	return data.(*structAppData)
}
