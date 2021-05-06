// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package configfile

//go:nosplit
func isSpace(t byte) bool {
	return (t <= 0x20) && (t > 0)
}

//go:nosplit
func isStringToken(ch byte) bool {
	return ch == 0x22 || ch == 0x27
}

//go:nosplit
func skipLineSpace(value []byte) []byte {
	i := 0
	for ; i < len(value); i++ {
		if value[i] == 0x0A || (!isSpace(value[i])) {
			break
		}
	}
	return value[i:]
}

//go:nosplit
func skipLine(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		if data[i] == 0x0A {
			return data[i+1:]
		}
	}
	return data[len(data):]
}

//go:nosplit
func getStringLen(data []byte, isEnd func(byte) bool) int {
	var stringToken byte = 0
	validTokens := 0
	preIsTrans := false
	preIsSpace := true
	i := 0
	for ; i < len(data); i++ {
		if preIsTrans {
			preIsTrans = false
			preIsSpace = false
			continue
		}
		if data[i] == '\\' {
			preIsSpace = false
			preIsTrans = true
			continue
		}
		if (validTokens & 1) == 1 {
			if data[i] == stringToken {
				validTokens++
			}
			continue
		}
		if isSpace(data[i]) {
			//行末
			if data[i] == 0x0A {
				break
			}
			preIsSpace = true
			continue
		}
		if isStringToken(data[i]) {
			preIsSpace = false
			validTokens++
			stringToken = data[i]
			continue
		}
		if data[i] == '#' && preIsSpace {
			break
		}
		if isEnd(data[i]) {
			break
		}
		preIsSpace = false
	}
	return i
}

//go:nosplit
func getChar(ch byte) byte {
	switch ch {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 'v':
		return '\v'
	case '\\':
		return '\\'
	default:
		return ch
	}
}
func parseString(data, output []byte) int {
	size := 0
	var strToken byte = 0
	insideStr := 0
	preIsTrans := false
	i := 0
	for ; i < len(data); i++ {
		if preIsTrans {
			preIsTrans = false
			output[size] = getChar(data[i])
			size++
			continue
		}
		if data[i] == '\\' {
			preIsTrans = true
			continue
		}
		if (insideStr & 1) == 1 {
			if data[i] == strToken {
				insideStr++
			} else {
				output[size] = data[i]
				size++
			}
			continue
		}
		if isStringToken(data[i]) {
			insideStr++
			strToken = data[i]
			continue
		}
		output[size] = data[i]
		size++
	}
	return size
}
func trimSpace(value string) string {
	begin := 0
	for ; begin < len(value); begin++ {
		if !isSpace(value[begin]) {
			break
		}
	}
	value = value[begin:]
	end := len(value)
	if end == 0 {
		return value
	}
	for ; end > 0; end-- {
		if !isSpace(value[end-1]) {
			break
		}
	}
	return value[:end]
}

// Value 配置项值
type Value struct {
	value   interface{}
	isArray bool
}

// IsArray 配置项值是否是一个Array
func (v *Value) IsArray() bool {
	return v.isArray
}

// Get 读取配置项值 : string
func (v *Value) Get() string {
	if v.isArray {
		return ""
	}
	return v.value.(string)
}

// ArrayCount 读取配置项Array的count (如果是Array)
func (v *Value) ArrayCount() int {
	if !v.isArray {
		return -1
	}
	a := v.value.([]*Value)
	return len(a)
}

// GetArrayItem 取Array子项
func (v *Value) GetArrayItem(index int) *Value {
	if !v.isArray {
		return nil
	}
	if a := v.value.([]*Value); len(a) > index {
		return a[index]
	}
	return nil
}
func (v *Value) valid() bool {
	if v.isArray {
		return true
	}
	return len(v.value.(string)) > 0
}
func (v *Value) scanf(data []byte, isEnd func(byte) bool) int {
	sizeRaw := len(data)
	data = skipLineSpace(data)
	if len(data) == 0 {
		v.isArray = false
		v.value = ""
		return sizeRaw
	}
	v.isArray = (data[0] == '[')
	if !v.isArray {
		slen := getStringLen(data, isEnd)
		if slen == 0 {
			v.value = ""
			return sizeRaw - len(data)
		}
		buffer := make([]byte, slen)
		strLen := parseString(data[:slen], buffer)
		v.value = trimSpace(string(buffer[:strLen]))
		data = data[slen:]
		return sizeRaw - len(data)
	}
	data = data[1:]
	itemArray := []*Value{}
	for len(data) > 0 {
		item := &Value{}
		parsed := item.scanf(data, func(ch byte) bool {
			return ch == ']' || isEnd(ch)
		})
		if item.valid() {
			itemArray = append(itemArray, item)
		}
		data = data[parsed:]
		if len(data) == 0 {
			break
		}
		if data[0] == ']' {
			data = skipLineSpace(data[1:])
			break
		}
		data = skipLine(data)
	}
	v.value = itemArray
	return sizeRaw - len(data)
}

// Config 配置项集合
type Config struct {
	values map[string]map[string]*Value
}

// Session 取配置项中指定名字的Session
func (c *Config) Session(sessionName string) map[string]*Value {
	if c.values == nil {
		return nil
	}
	if v, found := c.values[sessionName]; found {
		return v
	}
	return nil
}

//Parse : 解析conf 文件, 变量使用key=value方式
func Parse(data []byte) *Config {
	r := &Config{values: map[string]map[string]*Value{}}
	session := ""
	for len(data) > 0 {
		value := &Value{}
		parsed := value.scanf(data, func(ch byte) bool {
			return ch == '=' || ch == 0x0A || ch == '#'
		})
		data = data[parsed:]
		if value.IsArray() {
			name := value.GetArrayItem(0)
			if name != nil && !name.IsArray() {
				session = trimSpace(name.Get())
			}
			data = skipLine(data)
			continue
		}
		if len(data) == 0 {
			break
		}
		if data[0] == '#' || data[0] == 0x0A {
			data = skipLine(data)
			continue
		}
		//data[0] == '='
		data = data[1:]
		propertyName := value.Get()
		parsed = value.scanf(data, func(ch byte) bool {
			return ch == 0x0A || ch == '#'
		})
		data = skipLine(data[parsed:])
		dataSession := r.Session(session)
		if dataSession == nil {
			r.values[session] = map[string]*Value{}
			dataSession = r.Session(session)
		}
		dataSession[propertyName] = value
	}
	return r
}
