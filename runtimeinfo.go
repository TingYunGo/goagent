// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

type RuntimeSnap struct {
	CountOfGoRoutine    uint64
	CountOfCgoCall      uint64
	CountOfGC           uint64
	CountOfFrees        uint64
	CountOfMallocs      uint64
	CountOfLookups      uint64
	CountOfFD           uint32
	CountOfThreads      uint32
	TimeOfGC_ns         uint64
	SizeOfMemTotalSys   uint64
	SizeOfMemStackSys   uint64
	SizeOfMSpanSys      uint64
	SizeOfMemHeapSys    uint64
	SizeOfMCacheSys     uint64
	SizeOfBuckHashSys   uint64
	SizeOfMemHeapInuse  uint64
	SizeOfMemStackInuse uint64
	SizeOfMSpanInuse    uint64
	SizeOfMCacheInuse   uint64
	SizeOfMemVM         uint64
	SizeOfMemRss        uint64
	CpuProcess          cpuInfo
	CpuSystem           cpuInfo
}

func (snap *RuntimeSnap) Snap() {
	snap.CountOfCgoCall = uint64(runtime.NumCgoCall())
	memState := &runtime.MemStats{}
	runtime.ReadMemStats(memState)
	snap.CountOfGoRoutine = uint64(runtime.NumGoroutine())
	snap.SizeOfMemTotalSys = memState.Sys
	snap.SizeOfMemHeapSys = memState.HeapSys
	snap.SizeOfMemStackSys = memState.StackSys
	snap.SizeOfMSpanSys = memState.MSpanSys
	snap.SizeOfMCacheSys = memState.MCacheSys
	snap.SizeOfBuckHashSys = memState.BuckHashSys
	snap.SizeOfMemHeapInuse = memState.HeapInuse
	snap.SizeOfMSpanInuse = memState.MSpanInuse
	snap.SizeOfMCacheInuse = memState.MCacheInuse
	snap.SizeOfMemStackInuse = memState.StackInuse
	snap.CountOfMallocs = memState.Mallocs
	snap.CountOfFrees = memState.Frees
	snap.CountOfLookups = memState.Lookups
	snap.TimeOfGC_ns = memState.PauseTotalNs
	snap.CountOfGC = uint64(memState.NumGC)

	if runtime.GOOS != "linux" {
		return
	}
	fdSize, _ := elementCount("/proc/self/fd/")
	threads, _ := elementCount("/proc/self/task/")
	snap.CountOfFD = uint32(fdSize)
	snap.CountOfThreads = uint32(threads)

	bytes, _ := ioutil.ReadFile("/proc/self/stat")
	array := strings.Split(string(bytes), " ")
	if len(array) < 15 {
		return
	}
	snap.CpuProcess.User, _ = strconv.ParseUint(array[13], 10, 64)
	snap.CpuProcess.Kernel, _ = strconv.ParseUint(array[14], 10, 64)

	bytes, _ = ioutil.ReadFile("/proc/stat")

	array = strings.Split(strings.Split(string(bytes), "\n")[0], " ")
	snap.CpuSystem.User, _ = strconv.ParseUint(array[1], 10, 64)
	snap.CpuSystem.Kernel, _ = strconv.ParseUint(array[3], 10, 64)
	snap.CpuSystem.Idle, _ = strconv.ParseUint(array[4], 10, 64)
	snap.CpuSystem.IoWait, _ = strconv.ParseUint(array[5], 10, 64)

	bytes, _ = ioutil.ReadFile("/proc/self/status")
	array = strings.Split(string(bytes), "\n")
	for _, it := range array {
		a := strings.Split(it, ":")
		if a[0] == "VmSize" {
			valarray := strings.Split(strings.Trim(strings.Trim(a[1], "\t"), " "), " ")
			snap.SizeOfMemVM, _ = strconv.ParseUint(valarray[0], 10, 64)
		} else if a[0] == "VmRSS" {
			valarray := strings.Split(strings.Trim(strings.Trim(a[1], "\t"), " "), " ")
			snap.SizeOfMemRss, _ = strconv.ParseUint(valarray[0], 10, 64)
		}
	}
}

func (snap *RuntimeSnap) Sub(lastSnap *RuntimeSnap) *runtimeBlock {
	res := &runtimeBlock{}
	res.NumGoroutine.AddValue(float64(snap.CountOfGoRoutine), 0)
	res.NumCgoCall.AddValue(float64(snap.CountOfCgoCall-lastSnap.CountOfCgoCall), 0)
	res.NumGC.AddValue(float64(snap.CountOfGC-lastSnap.CountOfGC), 0)
	res.GCTime.AddValue(float64(snap.TimeOfGC_ns-lastSnap.TimeOfGC_ns)/1000000, 0)
	res.MemTotalSys.AddValue(float64(snap.SizeOfMemTotalSys)/1048576, 0)
	res.MemStackSys.AddValue(float64(snap.SizeOfMemStackSys)/1048576, 0)
	res.MSpanSys.AddValue(float64(snap.SizeOfMSpanSys)/1048576, 0)
	res.MCacheSys.AddValue(float64(snap.SizeOfMCacheSys)/1048576, 0)
	res.MemHeapSys.AddValue(float64(snap.SizeOfMemHeapSys)/1048576, 0)
	res.BuckHashSys.AddValue(float64(snap.SizeOfBuckHashSys)/1048576, 0)
	res.Frees.AddValue(float64(snap.CountOfFrees-lastSnap.CountOfFrees), 0)
	res.Mallocs.AddValue(float64(snap.CountOfMallocs-lastSnap.CountOfMallocs), 0)
	res.Lookups.AddValue(float64(snap.CountOfLookups), 0)
	res.HeapInuse.AddValue(float64(snap.SizeOfMemHeapInuse)/1048576, 0)
	res.StackInuse.AddValue(float64(snap.SizeOfMemStackInuse)/1048576, 0)
	res.MSpanInuse.AddValue(float64(snap.SizeOfMSpanInuse)/1048576, 0)
	res.MCacheInuse.AddValue(float64(snap.SizeOfMCacheInuse)/1048576, 0)
	processUsed := snap.CpuProcess.ProcessUse() - lastSnap.CpuProcess.ProcessUse()
	res.UserTime.AddValue(float64(processUsed)/100, 0)
	hostCpuUsed := snap.CpuSystem.FullUse() - lastSnap.CpuSystem.FullUse()
	if hostCpuUsed == 0 {
		hostCpuUsed = 1
	}
	res.UserUtilization.AddValue(float64(processUsed)*100/float64(hostCpuUsed), 0)
	res.FDSize.AddValue(float64(snap.CountOfFD), 0)
	res.Threads.AddValue(float64(snap.CountOfThreads), 0)
	res.VMPeakSize.AddValue(float64(snap.SizeOfMemVM)/1024, 0)
	res.RssPeak.AddValue(float64(snap.SizeOfMemRss)/1024, 0)
	res.mem.AddValue(float64(snap.SizeOfMemRss)/1024, 0)
	return res
}

type runtimeBlock struct {
	NumGoroutine structPerformance
	NumCgoCall   structPerformance
	NumGC        structPerformance
	GCTime       structPerformance //1分钟
	MemTotalSys  structPerformance
	MemStackSys  structPerformance
	MSpanSys     structPerformance
	MCacheSys    structPerformance
	MemHeapSys   structPerformance
	BuckHashSys  structPerformance
	Frees        structPerformance
	Mallocs      structPerformance
	Lookups      structPerformance
	HeapInuse    structPerformance
	StackInuse   structPerformance
	MSpanInuse   structPerformance
	MCacheInuse  structPerformance

	UserTime        structPerformance
	UserUtilization structPerformance
	mem             structPerformance
	FDSize          structPerformance
	Threads         structPerformance
	VMPeakSize      structPerformance
	RssPeak         structPerformance
}

func metricValue(name string, perf *structPerformance) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"value": map[string]interface{}{
			"count": perf.accessCount,
			"total": perf.sum,
			"min":   perf.valueMin,
			"max":   perf.valueMax,
		},
	}
}

func (p *runtimeBlock) Merge(q *runtimeBlock) {
	p.NumGoroutine.Append(&q.NumGoroutine)
	p.NumCgoCall.Append(&q.NumCgoCall)
	p.NumGC.Append(&q.NumGC)
	p.GCTime.Append(&q.GCTime)
	p.MemTotalSys.Append(&q.MemTotalSys)
	p.MemStackSys.Append(&q.MemStackSys)
	p.MSpanSys.Append(&q.MSpanSys)
	p.MCacheSys.Append(&q.MCacheSys)
	p.MemHeapSys.Append(&q.MemHeapSys)
	p.BuckHashSys.Append(&q.BuckHashSys)
	p.Frees.Append(&q.Frees)
	p.Mallocs.Append(&q.Mallocs)
	p.Lookups.Append(&q.Lookups)
	p.HeapInuse.Append(&q.HeapInuse)
	p.StackInuse.Append(&q.StackInuse)
	p.MSpanInuse.Append(&q.MSpanInuse)
	p.MCacheInuse.Append(&q.MCacheInuse)

	p.UserTime.Append(&q.UserTime)
	p.UserUtilization.Append(&q.UserUtilization)
	p.mem.Append(&q.mem)
	p.FDSize.Append(&q.FDSize)
	p.Threads.Append(&q.Threads)
	p.VMPeakSize.Append(&q.VMPeakSize)
	p.RssPeak.Append(&q.RssPeak)
}
func (p *runtimeBlock) Serialize() ([]byte, error) {
	ret := make([]interface{}, 24)
	ret[0] = metricValue("GoRuntime/NULL/Goroutine", &p.NumGoroutine)
	ret[1] = metricValue("GoRuntime/NULL/CgoCall", &p.NumCgoCall)
	ret[2] = metricValue("GoRuntime/NULL/Frees", &p.Frees)
	ret[3] = metricValue("GoRuntime/NULL/Mallocs", &p.Mallocs)
	ret[4] = metricValue("GoRuntime/NULL/Lookups", &p.Lookups)
	ret[5] = metricValue("GC/NULL/Count", &p.NumGC)
	ret[6] = metricValue("GC/NULL/Time", &p.GCTime)
	ret[7] = metricValue("CPU/NULL/UserTime", &p.UserTime)
	ret[8] = metricValue("CPU/NULL/UserUtilization", &p.UserUtilization)
	ret[9] = metricValue("Memory/NULL/PhysicalUsed", &p.mem)
	ret[10] = metricValue("Memory/NULL/VmSize", &p.VMPeakSize)
	ret[11] = metricValue("Memory/NULL/VmRSS", &p.RssPeak)
	ret[12] = metricValue("Memory/NULL/MemSys", &p.MemTotalSys)
	ret[13] = metricValue("Memory/Stack/StackSys", &p.MemStackSys)
	ret[14] = metricValue("Memory/MSpan/MSpanSys", &p.MSpanSys)
	ret[15] = metricValue("Memory/MCache/MCacheSys", &p.MCacheSys)
	ret[16] = metricValue("Memory/Heap/HeapSys", &p.MemHeapSys)
	ret[17] = metricValue("Memory/NULL/BuckHashSys", &p.BuckHashSys)
	ret[18] = metricValue("Memory/Heap/HeapInuse", &p.HeapInuse)
	ret[19] = metricValue("Memory/Stack/StackInuse", &p.StackInuse)
	ret[20] = metricValue("Memory/MSpan/MSpanInuse", &p.MSpanInuse)
	ret[21] = metricValue("Memory/MCache/MCacheInuse", &p.MCacheInuse)
	ret[22] = metricValue("FD/NULL/Count", &p.FDSize)
	ret[23] = metricValue("Thread/NULL/Count", &p.Threads)
	return json.Marshal(ret)
}
