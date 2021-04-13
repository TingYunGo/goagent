// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"fmt"
	"strings"
	"time"

	"git.codemonky.net/TingYunGo/goagent/protoc"
	proto "github.com/golang/protobuf/proto"
)

type structAppData struct {
	sys         *sysInfo
	runtimeData runtimeBlock
	traces      *protoc.ActionTraces
}

func (r *structAppData) init() *structAppData {
	r.traces = nil
	r.sys = nil
	return r
}

func (r *structAppData) end(perf *runtimePerf) {
	perf.Snap()
	r.runtimeData.Read(perf)
	perf.Reset()
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
	if r.traces == nil {
		r.traces = &protoc.ActionTraces{}
	}

	trace := &protoc.ActionTrace{}
	trace.Action = action.GetName()
	trace.Time = action.time.begin.Unix()
	trace.Duration = (int64)(action.time.duration / time.Millisecond)
	trace.Rid = unicID(action.time.begin, action)
	if array := strings.Split(action.trackID, ";"); len(array) > 1 {
		for i := 0; i < len(array); i++ {
			item := array[i]
			switch string(item[0:2]) {
			case "e=":
				trace.Refid = item[2:]
			case "c=":
				trace.Cross = item[2:]
			case "n=":
				trace.Tmd5 = item[2:]
			case "x=":
				trace.Tid = item[2:]
			default:
			}
		}
		// entryTrace := parseTrackId(action.trackId)
		// if entryTrace != nil {
		// 	entryTrace["time"] = mkTime(action.time.duration, 0, 0, 0, 0, 0, 0)
		// 	action.customParams["entryTrace"] = entryTrace
		// }

	} else {
		trace.Tmd5 = md5sum(trace.Action)
		trace.Refid = action.trackID
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
	//trace.NoSample: 当前是否启用采样
	//
	detail := trace.Detail
	if detail == nil {
		trace.Detail = &protoc.TraceDetail{}
		detail = trace.Detail
	}
	detail.Custom = action.customParams
	detail.QueryStringParameters = action.url
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
		if len(component.errors) > 0 {
			for i := 0; i < len(component.errors); i++ {
				exception := &protoc.TracerException{}
				exception.Error = true
				exception.Msg = fmt.Sprint(component.errors[i].e)
				exception.Name = component.errors[i].eType
				exception.Stack = component.errors[i].stack
				traceItem.Exception = append(traceItem.Exception, exception)
			}
		}
		var params *protoc.TracerParams = nil
		if component._type != ComponentDefault {
			params = &protoc.TracerParams{}
			params.Operation = component.sql
			params.Protocol = component.protocol
			params.Key = component.table
			params.Instance = component.instance
			params.Vendor = component.getVender()
			params.TxData = component.txdata
			traceItem.Params = params
		}
		detail.Tracers = append(detail.Tracers, traceItem)
	}
	action.root = nil
	// componentMap.Clear()
	//添加trace
	r.traces.Traces = append(r.traces.Traces, trace)
}

const (
	//采样值,采样时刻进程内的go程数
	metricNumGoroutine = "GoRuntime/NULL/Goroutine"
	//单位时间内事件次数, =>一个累加值 在两次采样之间的差值。
	metricNumCgoCall = "GoRuntime/NULL/CgoCall"
	//单位时间内GC耗时的累加和,单位毫秒
	metricPauseTotalMs = "GoRuntime/NULL/PauseTotalMs"
	//单位时间内,每次GC耗时的5值统计性能数据
	metricGCTime = "GC/NULL/Time"
	//单位时间内 Free的次数
	metricFrees = "GoRuntime/NULL/Frees"
	//单位时间内 Malloc的次数
	metricMallocs = "GoRuntime/NULL/Mallocs"
	//单位时间内 Lookup的次数
	metricLookups = "GoRuntime/NULL/Lookups"
	//采样值,系统总的申请内存数 MB
	metricMemTotalSys = "Memory/NULL/MemSys"
	//采样值,系统栈内存数 MB
	metricMemStackSys = "Memory/Stack/StackSys"
	//采样值,系统堆内存数 MB
	metricMemHeapSys = "Memory/Heap/HeapSys"
	//采样值,系统内存区间结构数
	metricMSpanSys = "Memory/MSpan/MSpanSys"
	//采样值,系统内存Cache结构数
	metricMCacheSys = "Memory/MCache/MCacheSys"
	//采样值,系统内存BuckHash数
	metricBuckHashSys = "Memory/NULL/BuckHashSys"
	//采样值,使用中的堆内存数 MB
	metricHeapInuse = "Memory/Heap/HeapInuse"
	//采样值,使用中的栈内存数 MB
	metricStackInuse = "Memory/Stack/StackInuse"
	//采样值,使用中的内存区间结构数
	metricMSpanInuse = "Memory/MSpan/MSpanInuse"
	//采样值,使用中的内存Cache结构数
	metricMCacheInuse     = "Memory/MCache/MCacheInuse"
	metricUserTime        = "CPU/NULL/UserTime"
	metricUserUtilization = "CPU/NULL/UserUtilization"
	metricmem             = "Memory/NULL/PhysicalUsed"
	//采样值,进程打开文件句柄数(linux)
	metricFDSize = "Process/NULL/FD"
	//采样值,进程内的系统线程数(linux)
	metricThreads = "Process/NULL/Threads"
)

//数据序列化
func (r *structAppData) Serialize() ([]byte, error) {
	return proto.Marshal(r.traces)
}

//释放内存
func (r *structAppData) destroy() {
	r.traces = nil
	r.sys = nil
}
func (a *application) GetReportBlock() *structAppData {
	if a.reportQueue.Size() == 0 {
		a.reportQueue.PushBack(new(structAppData).init())
	}
	data, _ := a.reportQueue.Back().Value()
	if datablock := data.(*structAppData); datablock.traces != nil && len(datablock.traces.Traces) >= 5000 {
		a.reportQueue.PushBack(new(structAppData).init())
	}
	data, _ = a.reportQueue.Back().Value()
	return data.(*structAppData)
}
