// Copyright 2022 冯立强 fenglq@tingyun.com.  All rights reserved.

package database

import (
	"database/sql"
	"database/sql/driver"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

type StructSQLDB110 struct {
	driver driver.Connector
}
type StructSQLDB111 struct {
	waitDuration int64
	driver       driver.Connector
}

func matchObject110(first, second *StructSQLDB110) bool {
	return first.driver == second.driver
}
func matchObject111(first, second *StructSQLDB111) bool {
	return first.driver == second.driver
}

func matchObject(first, second *sql.DB) bool {
	vervect := strings.Split(runtime.Version(), ".")
	subver, _ := strconv.ParseInt(vervect[1], 10, 32)
	if subver < 11 {
		return matchObject110((*StructSQLDB110)(unsafe.Pointer(first)), (*StructSQLDB110)(unsafe.Pointer(second)))
	} else {
		return matchObject111((*StructSQLDB111)(unsafe.Pointer(first)), (*StructSQLDB111)(unsafe.Pointer(second)))
	}
}
