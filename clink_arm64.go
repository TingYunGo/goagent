// Copyright 2022 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build arm64
// +build cgo

package tingyun3

/*
#cgo LDFLAGS: -L${SRCDIR} -ltingyungoarm64

extern int tingyun_go_init(void *);

*/
import "C"
