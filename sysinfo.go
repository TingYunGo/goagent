// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type cpuInfo struct {
	User   uint64
	Kernel uint64
	Idle   uint64
	IoWait uint64
}

func (i *cpuInfo) ProcessUse() uint64 {
	return i.User + i.Kernel
}
func (i *cpuInfo) FullUse() uint64 {
	return i.User + i.Kernel + i.Idle + i.IoWait
}

type sysInfo struct {
	cpuProcess cpuInfo
	cpuSystem  cpuInfo
	vmSize     uint64
	vmRss      uint64
	Threads    int
	FdSize     int
	err        error
}

func (s *sysInfo) Init() *sysInfo {
	s.err = nil
	if runtime.GOOS != "linux" {
		s.err = errors.New("not linux")
		return nil
	}

	s.FdSize, s.err = elementCount("/proc/self/fd/")
	s.Threads, s.err = elementCount("/proc/self/task/")
	bytes, err := ioutil.ReadFile("/proc/self/stat")
	if err != nil {
		s.err = err
		return nil
	}
	array := strings.Split(string(bytes), " ")
	if len(array) < 15 {
		s.err = errors.New("Unknown or UnSupport linux version")
		return nil
	}
	s.cpuProcess.User, _ = strconv.ParseUint(array[13], 10, 64)
	s.cpuProcess.Kernel, _ = strconv.ParseUint(array[14], 10, 64)

	bytes, err = ioutil.ReadFile("/proc/stat")
	if err != nil {
		s.err = err
		return nil
	}
	array = strings.Split(strings.Split(string(bytes), "\n")[0], " ")
	s.cpuSystem.User, _ = strconv.ParseUint(array[1], 10, 64)
	s.cpuSystem.Kernel, _ = strconv.ParseUint(array[3], 10, 64)
	s.cpuSystem.Idle, _ = strconv.ParseUint(array[4], 10, 64)
	s.cpuSystem.IoWait, _ = strconv.ParseUint(array[5], 10, 64)

	bytes, err = ioutil.ReadFile("/proc/self/status")
	if err != nil {
		s.err = err
		return nil
	}
	array = strings.Split(string(bytes), "\n")
	for _, it := range array {
		a := strings.Split(it, ":")
		if a[0] == "VmSize" {
			valarray := strings.Split(strings.Trim(strings.Trim(a[1], "\t"), " "), " ")
			s.vmSize, _ = strconv.ParseUint(valarray[0], 10, 64)
		} else if a[0] == "VmRSS" {
			valarray := strings.Split(strings.Trim(strings.Trim(a[1], "\t"), " "), " ")
			s.vmRss, _ = strconv.ParseUint(valarray[0], 10, 64)
		}
	}
	return s
}
func elementCount(dir string) (int, error) {
	file, err := os.Open(dir)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	names, err := file.Readdirnames(0)
	if err != nil {
		return -1, err
	}
	return len(names), nil
}
