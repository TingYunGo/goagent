// Copyright 2022 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64
// +build cgo

package tingyun3

/*
#cgo LDFLAGS: -L${SRCDIR} -ltingyungosdk

extern int tingyun_go_init(void *);

*/
import "C"
