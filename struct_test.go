// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"fmt"
	"strings"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(md5sum("1234567890"))
}

func TestFileEnum(t *testing.T) {
	dir := "c:\\"
	count, err := elementCount(dir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("count of dir \"%s\" = %d\n", dir, count)
	}
}
