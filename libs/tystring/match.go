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
