// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

type structPerformance struct {
	sum         float64 //累加和
	valueMax    float64 //最大值
	valueMin    float64 //最小值
	accessCount int32   //累加次数
}

func (p *structPerformance) IntSlice() []int64 {
	r := make([]int64, 4)
	r[0] = int64(p.accessCount)
	r[1] = int64(p.sum)
	r[2] = int64(p.valueMax)
	r[3] = int64(p.valueMin)
	return r
}
func (p *structPerformance) FloatSlice() []interface{} {
	r := make([]interface{}, 4)
	r[0] = p.accessCount
	r[1] = p.sum
	r[2] = p.valueMax
	r[3] = p.valueMin
	return r
}

func newStructPerformance() *structPerformance {
	r := &structPerformance{}
	r.Reset()
	return r
}
func (p *structPerformance) Reset() *structPerformance {
	p.sum = 0
	p.valueMax = 0
	p.valueMin = 0
	p.accessCount = 0
	return p
}

func (p *structPerformance) Append(q *structPerformance) {
	if q.accessCount > 0 {
		p.sum += q.sum
		if p.accessCount == 0 || q.valueMax > p.valueMax {
			p.valueMax = q.valueMax
		}
		if p.accessCount == 0 || q.valueMin < p.valueMin {
			p.valueMin = q.valueMin
		}
		p.accessCount += q.accessCount
	}
}
func (p *structPerformance) AddComponent(value float64, excl float64) {
	p.sum += value
	if p.accessCount == 0 {
		p.valueMax = excl
		p.valueMin = excl
	} else {
		if p.valueMax < excl {
			p.valueMax = excl
		}
		if p.valueMin > excl {
			p.valueMin = excl
		}
	}
	p.accessCount++
}
func (p *structPerformance) AddValue(value float64, excl float64) {
	p.sum += value
	if p.accessCount == 0 {
		p.valueMax = value
		p.valueMin = value
	} else {
		if p.valueMax < value {
			p.valueMax = value
		}
		if p.valueMin > value {
			p.valueMin = value
		}
	}
	p.accessCount++
}
