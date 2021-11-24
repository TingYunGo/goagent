// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"time"
)

/*
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Stack                []string `protobuf:"bytes,3,rep,name=stack,proto3" json:"stack,omitempty"`
	Error                bool     `protobuf:"varint,4,opt,name=error,proto3" json:"error,omitempty"`

*/
type errInfo struct {
	happenTime time.Time
	e          string
	stack      []string
	eType      string
	isError    bool
}

func (i *errInfo) Destroy() {
	i.e = ""
	i.stack = []string{}
}
