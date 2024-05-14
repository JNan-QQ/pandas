/**
  Copyright (c) [2024] [JiangNan]
  [go-tools] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package series

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type Element interface {
	// Set 设置值
	Set(any)
	// Records 返回字符串
	Records() string
	// Int 返回整数
	Int() int
	// Float 返回浮点数
	Float() float64
	// Bool 返回布尔值
	Bool() bool
	// Value 任一类型
	Value() any
	// Type 返回数据类型
	dType() string
	// Copy 复制值
	copy() Element
	// 判断是否为NaN
	isNaN() bool
	// 更新
	update(Element)
}

// 字符串数据格式，实现接口 Element
type stringElement string

// 整数数据格式，实现接口 Element
type intElement int

// 浮点数数据格式，实现接口 Element
type floatElement float64

// 布尔值数据格式，实现接口 Element
type boolElement bool

//////////////////////////////////////
//            元素                  //
/////////////////////////////////////

func (s *stringElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if slices.Contains([]string{"", "NaN", "nan", "null", "Null"}, val) {
			*s = "NaN"
		} else {
			*s = stringElement(val)
		}
	case int:
		*s = stringElement(strconv.Itoa(val))
	case float64:
		*s = stringElement(strconv.FormatFloat(value.(float64), 'f', 6, 64))
	case bool:
		if val {
			*s = "true"
		} else {
			*s = "false"
		}
	default:
		*s = "NaN"
	}
}
func (i *intElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if newInt, err := strconv.Atoi(val); err == nil {
			*i = intElement(newInt)
		} else {
			*i = math.MinInt
		}
	case int:
		*i = intElement(val)
	case float64:
		*i = intElement(val)
	case bool:
		if val {
			*i = 1
		} else {
			*i = 0
		}
	default:
		*i = math.MinInt
	}
}
func (f *floatElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if newFloat, err := strconv.ParseFloat(val, 10); err != nil {
			*f = floatElement(math.NaN())
		} else {
			*f = floatElement(newFloat)
		}
	case int:
		*f = floatElement(val)
	case float64:
		*f = floatElement(val)
	case bool:
		if val {
			*f = 1
		} else {
			*f = 0
		}
	default:
		*f = floatElement(math.NaN())
	}
}
func (b *boolElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if slices.Contains([]string{"false", "0", "F", "f"}, val) {
			*b = false
		} else {
			*b = true
		}
	case int, float64:
		if val == 0 {
			*b = false
		} else {
			*b = true
		}
	case bool:
		*b = boolElement(val)
	default:
		panic("错误的数据类型")
	}
}

//--------------------------------//

func (s *stringElement) Records() string {
	return string(*s)
}
func (i *intElement) Records() string {
	if i.isNaN() {
		return "NaN"
	}
	return strconv.Itoa(int(*i))
}
func (f *floatElement) Records() string {
	if f.isNaN() {
		return "NaN"
	}
	return strconv.FormatFloat(float64(*f), 'f', -1, 64)
}
func (b *boolElement) Records() string {
	if *b {
		return "true"
	} else {
		return "false"
	}
}

//--------------------------------//

func (s *stringElement) Int() int {
	if s.isNaN() {
		return math.MinInt
	}
	if i, err := strconv.Atoi(string(*s)); err != nil {
		fmt.Printf("%s 不能转换为 Int，已置为无限小\n", string(*s))
		return math.MinInt
	} else {
		return i
	}
}
func (i *intElement) Int() int {
	return int(*i)
}
func (f *floatElement) Int() int {
	return int(*f)
}
func (b *boolElement) Int() int {
	if *b {
		return 1
	} else {
		return 0
	}
}

//--------------------------------//

func (s *stringElement) Float() float64 {
	if s.isNaN() {
		return math.NaN()
	}
	if f, err := strconv.ParseFloat(string(*s), 64); err != nil {
		return math.NaN()
	} else {
		return f
	}
}
func (i *intElement) Float() float64 {
	if *i == math.MinInt {
		return math.NaN()
	}
	return float64(*i)
}
func (f *floatElement) Float() float64 {
	return float64(*f)
}
func (b *boolElement) Float() float64 {
	if *b {
		return 1
	} else {
		return 0
	}
}

//--------------------------------//

func (s *stringElement) Bool() bool {
	if slices.Contains([]string{"false", "0", "F", "f"}, strings.ToLower(string(*s))) {
		return false
	} else {
		return true
	}
}
func (i *intElement) Bool() bool {
	if *i > 0 {
		return true
	} else {
		return false
	}
}
func (f *floatElement) Bool() bool {
	if *f > 0 && !math.IsNaN(float64(*f)) {
		return true
	} else {
		return false
	}
}
func (b *boolElement) Bool() bool {
	return bool(*b)
}

//--------------------------------//

func (s *stringElement) dType() string {
	return string(String)
}
func (i *intElement) dType() string {
	return string(Int)
}
func (f *floatElement) dType() string {
	return string(Float)
}
func (b *boolElement) dType() string {
	return string(Bool)
}

//--------------------------------//

func (s *stringElement) copy() Element {
	s2 := new(stringElement)
	*s2 = *s
	return s2
}
func (i *intElement) copy() Element {
	i2 := new(intElement)
	*i2 = *i
	return i2
}
func (f *floatElement) copy() Element {
	f2 := new(floatElement)
	*f2 = *f
	return f2
}
func (b *boolElement) copy() Element {
	b2 := new(boolElement)
	*b2 = *b
	return b2
}

//--------------------------------//

func (s *stringElement) isNaN() bool {
	if *s == "NaN" {
		return true
	}
	return false
}
func (i *intElement) isNaN() bool {
	if *i == math.MinInt {
		return true
	}
	return false
}
func (f *floatElement) isNaN() bool {
	if math.IsNaN(float64(*f)) {
		return true
	}
	return false
}
func (b *boolElement) isNaN() bool {
	return false
}

//--------------------------------//

func (s *stringElement) String() string {
	if s.isNaN() {
		return "NaN"
	} else {
		return s.Records()
	}
}
func (i *intElement) String() string {
	if *i == math.MinInt {
		return "NaN"
	} else {
		return i.Records()
	}
}

//--------------------------------//

func (s *stringElement) update(elem Element) {
	s.Set(elem.Records())
}

func (i *intElement) update(elem Element) {
	i.Set(elem.Int())
}

func (f *floatElement) update(elem Element) {
	f.Set(elem.Float())
}

func (b *boolElement) update(elem Element) {
	b.Set(elem.Bool())
}

//--------------------------------//

func (s *stringElement) Value() any {
	return s.Records()
}

func (f *floatElement) Value() any {
	return f.Float()
}

func (b *boolElement) Value() any {
	return b.Bool()
}

func (i *intElement) Value() any {
	return i.Int()
}
