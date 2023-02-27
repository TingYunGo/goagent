// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tystring

//go:nosplit
func IsAlpha(ch uint8) bool {
	t := (ch | ('A' ^ 'a')) - 'a'
	return t <= 'z'-'a'
}

//go:nosplit
func ToLower(ch uint8) uint8 {

	if t := ch | ('A' ^ 'a'); t-'a' <= 'z'-'a' {
		return t
	}
	return ch
}

//go:nosplit
func min(a, b int) int {

	if a < b {
		return a
	}
	return b
}

//CaseCMP: a < b => -1; a > b => 1; a == b => 0
func CaseCMP(a, b string) int {
	size := min(len(a), len(b))
	for i := 0; i < size; i++ {
		lower_char_a := ToLower(a[i])
		lower_char_b := ToLower(b[i])
		if lower_char_a == lower_char_b {
			continue
		}
		if lower_char_b < lower_char_a {
			return 1
		}
		return -1
	}
	if len(a) == len(b) {
		return 0
	}
	if len(a) < len(b) {
		return -1
	}
	return 1
}

//go:nosplit
func SubString(str string, begin, size int) string {
	if len(str) <= begin {
		return ""
	}
	return str[begin:min(len(str), begin+size)]
}
func FindString(array []string, target string) int {
	Range := len(array)
	begin := 0
	for Range > 0 {
		middle := begin + Range/2
		if comp := CaseCMP(target, array[middle]); comp == 0 {
			return middle

		} else if comp < 0 {
			Range = middle - begin

		} else {
			Range -= middle + 1 - begin
			begin = middle + 1
		}
	}
	return -1
}
func IsSpace(t byte) bool {
	return (t <= 0x20) && (t > 0)
}
func TrimString(value string, isSep func(byte) bool) string {
	begin := 0
	for ; begin < len(value); begin++ {
		if !isSep(value[begin]) {
			break
		}
	}
	value = value[begin:]
	end := len(value)
	if end == 0 {
		return value
	}
	for ; end > 0; end-- {
		if !isSep(value[end-1]) {
			break
		}
	}
	return value[:end]
}

func SplitMapString(source string, isSep func(byte) bool, handler func(string, string)) {
	if handler == nil || isSep == nil {
		return
	}
	source = TrimString(source, IsSpace)
	keyLen := len(source)
	for i := 0; i < keyLen; i++ {
		if isSep(source[i]) {
			keyLen = i
			break
		}
	}
	value := ""
	found := false
	for i := keyLen; i < len(source); i++ {
		if isSep(source[i]) {
			found = true
		} else {
			value = source[i:]
			break
		}
	}
	if keyLen > 0 || found || len(value) > 0 {
		handler(source[:keyLen], value)
	}
	return
}

func SplitStrings(source string, isSep func(byte) bool, handler func(string) bool) {
	if handler == nil || isSep == nil {
		return
	}
	begin := -1
	for i := 0; i < len(source); i++ {
		if isSep(source[i]) {
			if begin > -1 {
				if handler(source[begin:i]) {
					return
				}
				begin = -1
			}
		} else {
			if begin == -1 {
				begin = i
			}
		}
	}
	if begin > -1 && begin < len(source) {
		handler(source[begin:len(source)])
	}
	return
}
