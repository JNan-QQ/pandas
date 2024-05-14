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
	"cmp"
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"github.com/shopspring/decimal"
	"math"
	"reflect"
	"slices"
	"strings"
)

type Series struct {
	Name     string
	elements []Element
	t        Type
	indexes  []int
}

type Type string

const (
	String Type = "string"
	Int    Type = "int"
	Float  Type = "float64"
	Bool   Type = "bool"
)

// NewSeries 创建数据列
//
//	values: 数据切片
//	dType: 数据类型，可选String、Int、float64、Bool
//	name: 数据列名称
func NewSeries[S interface{ ~[]E }, E int | float64 | string | bool](values S, dType Type, name string) (*Series, error) {
	s := &Series{
		Name:     name,
		t:        dType,
		elements: NewElements(dType, 0),
	}

	if values == nil {
		return s, nil
	} else if reflect.TypeOf(values).String() != fmt.Sprintf("[]%s", dType) {
		return nil, fmt.Errorf("输入切片与指定数据类型不匹配")
	} else {
		if err := s.Append(values); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// LoadRecords 用字符串切片创建指定类型数据列
//
//	values: 数据切片
//	dType: 数据类型，可选String、Int、float64、Bool
//	name: 数据列名称
func LoadRecords(values []string, t Type, name string) *Series {
	ns := &Series{
		Name:     name,
		elements: NewElements(t, len(values)),
		t:        t,
	}
	for i, value := range values {
		ns.Element(i).Set(value)
	}
	ns.InitIndex()
	return ns
}

// InitIndex 重置索引
func (s *Series) InitIndex() {
	s.indexes = make([]int, 0)
	for i := 0; i < s.Len(); i++ {
		s.indexes = append(s.indexes, i)
	}
}

// Append 向数据集后添加元素,可以为单个元素或元素切片,最好保证数据类型一致
//
//	values：Series、int []int、string []string ...
func (s *Series) Append(values interface{}) error {
	if reflect.TypeOf(values).Kind() == reflect.Slice {
		value := reflect.ValueOf(values)
		for i := 0; i < value.Len(); i++ {
			if err := s.Append(value.Index(i).Interface()); err != nil {
				return err
			}
		}
	} else {
		var x Element
		switch v := values.(type) {
		case int:
			x = new(intElement)
			x.Set(v)
		case float64:
			x = new(floatElement)
			x.Set(v)
		case string:
			x = new(stringElement)
			x.Set(v)
		case bool:
			x = new(boolElement)
			x.Set(v)
		case Element:
			x = v
		case Series:
			s.elements = append(s.elements, v.elements...)
			s.InitIndex()
			return nil
		default:
			return fmt.Errorf("不支持的数据类型 %v, %s", v, reflect.TypeOf(v))
		}
		s.elements = append(s.elements, x.copy())
	}
	s.InitIndex()
	return nil
}

// Drop 删除指定索引的元素
func (s *Series) Drop(indexes ...int) *Series {
	slices.SortFunc(slices.Compact(indexes), func(a, b int) int {
		return cmp.Compare(b, a)
	})
	ns := s.Copy()
	for _, index := range indexes {
		ns.elements = slices.Delete(ns.elements, index, index+1)
		ns.indexes = slices.Delete(ns.indexes, index, index+1)
	}
	return ns
}

// Concat 将新数据集加原数据集后,如果类型不同，以原数据集为准
func (s *Series) Concat(x Series) error {
	if x.t != s.t {
		fmt.Printf("两个数据类型不同，正在尝试转换...")
		if err := x.SetType(s.t); err != nil {
			return err
		}
		fmt.Println(" done")
	}

	s.elements = slices.Concat(s.elements, x.elements)

	s.InitIndex()

	return nil
}

// SubSet 保留对应索引的元素
func (s *Series) SubSet(indexes ...int) (*Series, error) {
	if slices.Max(indexes) >= len(s.elements) {
		return nil, fmt.Errorf("index out of range")
	}
	var elements []Element
	for _, index := range indexes {
		elements = append(elements, s.elements[index].copy())
	}

	newSeries := Series{
		Name:     s.Name,
		elements: elements,
		t:        s.t,
		indexes:  indexes,
	}

	return &newSeries, nil
}

// Elements 返回数据集元素对象切片
func (s *Series) Elements() []Element {
	return s.elements
}

// Element 指定索引元素
func (s *Series) Element(i int) Element {
	return s.elements[i]
}

// Format 批量处理数据集
//
//	index: 元素索引
//	elem: 元素对象
func (s *Series) Format(f func(index int, elem Element) Element) {
	for i, element := range s.elements {
		s.elements[i].update(f(i, element))
	}
}

// SetType 改变数据集类型
func (s *Series) SetType(t Type) error {
	if s.t == t {
		return nil
	}

	newSeries := Series{
		Name:     s.Name,
		elements: NewElements(t, 0),
		t:        t,
		indexes:  s.indexes,
	}

	var values any
	switch t {
	case String:
		values = s.Records()
	case Int:
		values = s.Int()
	case Float:
		values = s.Float()
	case Bool:
		values = s.Bool()
	default:
		return fmt.Errorf("未知数据类型！")
	}

	if err := newSeries.Append(values); err != nil {
		return err
	}

	*s = newSeries

	return nil
}

// SortIndex 生成升序、降序索引变化
func (s *Series) SortIndex(reverse bool) []int {
	var elements = s.Copy().Elements()
	var indexes = slices.Clone(s.indexes)
	for i := 0; i < s.Len(); i++ {
		ele := slices.MinFunc(elements[i:], func(a, b Element) int {
			if reverse {
				a, b = b, a
			}
			switch s.t {
			case String:
				return cmp.Compare(a.Records(), b.Records())
			case Int:
				return cmp.Compare(a.Int(), b.Int())
			case Float:
				return cmp.Compare(a.Float(), b.Float())
			default:
				panic("bool 不支持排序")
			}
		})

		j := slices.Index(elements[i:], ele) + i
		elements[i], elements[j] = elements[j], elements[i]
		indexes[i], indexes[j] = indexes[j], indexes[i]
	}

	return indexes
}

// 自定义输出
func (s *Series) String() string {
	return fmt.Sprintf("字段名：%s\n数 据：%v\n索 引：%v\n类 型：%s\n", s.Name, s.Records(), s.indexes, s.t)
}

// Len 返回数据集大小
func (s *Series) Len() int {
	return len(s.elements)
}

// Type 返回类型
func (s *Series) Type() string {
	return string(s.t)
}

// HasNaN 判断是否存在空值
func (s *Series) HasNaN() bool {
	for _, element := range s.elements {
		if element.isNaN() {
			return true
		}
	}
	return false
}

// Copy 复制
func (s *Series) Copy() *Series {
	ns := Series{
		Name:     s.Name,
		elements: NewElements(s.t, s.Len()),
		t:        s.t,
		indexes:  s.indexes,
	}
	for i, element := range s.elements {
		ns.elements[i].update(element)
	}

	return &ns
}

// Records 将数据集中的元素作为字符串返回
func (s *Series) Records() []string {
	var x []string
	for _, element := range s.elements {
		x = append(x, element.Records())
	}
	return x
}

// Float 将数据集中的元素作为浮点数返回，如果数据转换失败自动设为 math.NaN()
func (s *Series) Float() []float64 {
	var x []float64
	for _, element := range s.elements {
		x = append(x, element.Float())
	}
	return x
}

// Int 将数据集中的元素作为整数返回，如果数据转换失败自动设为 math.MinInt()
func (s *Series) Int() []int {
	var x []int
	for _, element := range s.elements {
		x = append(x, element.Int())
	}
	return x
}

// Bool 将数据集中的元素作为浮布尔值返回
func (s *Series) Bool() []bool {
	var x []bool
	for _, element := range s.elements {
		x = append(x, element.Bool())
	}
	return x
}

func (s *Series) Any() []any {
	var x []any
	for _, element := range s.elements {
		x = append(x, element.Value())
	}
	return x
}

// Indexes 返回索引切片
func (s *Series) Indexes() []int {
	return s.indexes
}

// RelationalOperator 关系运算
type RelationalOperator int

// ArithmeticOperator 算术运算
type ArithmeticOperator int

const (
	// Equal 相等
	Equal RelationalOperator = iota
	// NotEqual 不相等
	NotEqual
	// LessThan 小于
	LessThan
	// LessOrEqual 小于等于
	LessOrEqual
	// GreaterThan 大于
	GreaterThan
	// GreaterOrEqual 大于等于
	GreaterOrEqual
	// Contains 包含
	Contains
	// StartsWith 以...开始
	StartsWith
	// EndsWith 以...结束
	EndsWith
	// In 在...列表里
	In
	// NotIn 不在...列表里
	NotIn
)

const (
	// Addition 加
	Addition ArithmeticOperator = iota
	// Subtraction 减
	Subtraction
	// Multiplication 乘
	Multiplication
	// Division 除
	Division
	// Remainder 求于
	Remainder
)

// Filter 过滤数据集
func (s *Series) Filter(operator RelationalOperator, values any) (*Series, error) {
	// 判断待比对数据类型，是否与原数据集相同
	vT := reflect.TypeOf(values).String()
	if !(vT == fmt.Sprintf("[]%s", s.t) || vT == string(s.t)) {
		return nil, fmt.Errorf("输入数据类型与数据集类型不同")
	}

	// 判断对应数据类型的方法是否合法
	switch s.t {
	case String:
		if !slices.Contains([]RelationalOperator{0, 1, 6, 7, 8, 9, 10}, operator) {
			return nil, fmt.Errorf("string 类型数据无法执行操作%d", operator)
		}

	case Int:
		if !slices.Contains([]RelationalOperator{0, 1, 2, 3, 4, 5, 9, 10}, operator) {
			return nil, fmt.Errorf("int 类型数据无法执行操作%d", operator)
		}
	case Float:
		if !slices.Contains([]RelationalOperator{0, 1, 2, 3, 4, 5, 9, 10}, operator) {
			return nil, fmt.Errorf("float64 类型数据无法执行操作%d", operator)
		}
	case Bool:
		if !slices.Contains([]RelationalOperator{0, 1}, operator) {
			return nil, fmt.Errorf("bool 类型数据无法执行操作%d", operator)
		}
	}

	if (operator == 9 || operator == 10) && reflect.TypeOf(values).Kind() != reflect.Slice {
		return nil, fmt.Errorf("in / NotIn 需要输入切片作为参数")
	}

	_, indexes := data.Filter(s.elements, func(element Element) bool {
		switch operator {
		case 0:
			return element.Value() == values
		case 1:
			return element.Value() != values
		case 2, 3, 4, 5:
			if math.IsNaN(element.Float()) {
				return false
			}
			if reflect.TypeOf(values).Kind() == reflect.Int {
				switch operator {
				case 2:
					return element.Int() < values.(int)
				case 3:
					return element.Int() <= values.(int)
				case 4:
					return element.Int() > values.(int)
				case 5:
					return element.Int() >= values.(int)
				}
			} else {
				switch operator {
				case 2:
					return element.Float() < values.(float64)
				case 3:
					return element.Float() <= values.(float64)
				case 4:
					return element.Float() > values.(float64)
				case 5:
					return element.Float() >= values.(float64)
				}
			}
		case 6:
			return strings.Contains(element.Records(), values.(string))
		case 7:
			return strings.HasPrefix(element.Records(), values.(string))
		case 8:
			return strings.HasSuffix(element.Records(), values.(string))
		case 9, 10:
			var newValues []any
			v := reflect.ValueOf(values)
			for i := 0; i < v.Len(); i++ {
				newValues = append(newValues, v.Index(i).Interface())
			}
			if operator == 9 {
				return slices.Contains(newValues, element.Value())
			} else {
				return !slices.Contains(newValues, element.Value())
			}
		}
		return false
	})

	if subSet, err := s.SubSet(indexes...); err != nil {
		return nil, err
	} else {
		return subSet, nil
	}
}

// NewElements 生成 Elements 接口的对象
func NewElements(t Type, l int) []Element {
	var ne []Element
	for i := 0; i < l; i++ {
		var x Element
		switch t {
		case String:
			x = new(stringElement)
		case Int:
			x = new(intElement)
		case Float:
			x = new(floatElement)
		case Bool:
			x = new(boolElement)
		}
		ne = append(ne, x)
	}
	return ne
}

// Arithmetic 算术运算
func (s *Series) Arithmetic(operator ArithmeticOperator, x Series) (*Series, error) {
	if s.t == Bool || x.t == Bool {
		return nil, fmt.Errorf("布尔值不支持算术运算")
	}
	if s.Len() != x.Len() {
		return nil, fmt.Errorf("长度不相等！")
	}
	if operator != 0 && (s.t == String || x.t == String) {
		return nil, fmt.Errorf("字符串不支持该操作")
	}

	ns := &Series{Name: s.Name}
	// 保证结果准确
	if s.t == String || x.t == String {
		ns.t = String
	} else if s.t == Float || x.t == Float {
		ns.t = Float
	} else {
		if operator == 3 {
			ns.t = Float
		} else {
			ns.t = Int
		}
	}
	ns.elements = NewElements(ns.t, s.Len())

	switch operator {
	case 0:
		for i := 0; i < s.Len(); i++ {
			if ns.t == String {
				ns.Element(i).Set(s.Element(i).Records() + x.Element(i).Records())
			} else {
				f, _ := decimal.NewFromFloat(s.Element(i).Float()).
					Add(decimal.NewFromFloat(x.Element(i).Float())).Float64()
				ns.Element(i).Set(f)
			}
		}
	case 1:
		for i := 0; i < s.Len(); i++ {
			f, _ := decimal.NewFromFloat(s.Element(i).Float()).
				Sub(decimal.NewFromFloat(x.Element(i).Float())).Float64()
			ns.Element(i).Set(f)
		}
	case 2:
		for i := 0; i < s.Len(); i++ {
			f, _ := decimal.NewFromFloat(s.Element(i).Float()).
				Mul(decimal.NewFromFloat(x.Element(i).Float())).Float64()
			ns.Element(i).Set(f)
		}
	case 3:
		for i := 0; i < s.Len(); i++ {
			if x.Element(i).Float() == 0 {
				ns.Element(i).Set(math.NaN())
			} else {
				f, _ := decimal.NewFromFloat(s.Element(i).Float()).
					Div(decimal.NewFromFloat(x.Element(i).Float())).Float64()
				ns.Element(i).Set(f)
			}
		}
	case 4:
		for i := 0; i < s.Len(); i++ {
			if x.Element(i).Float() == 0 {
				ns.Element(i).Set(math.NaN())
			} else {
				f, _ := decimal.NewFromFloat(s.Element(i).Float()).
					Mod(decimal.NewFromFloat(x.Element(i).Float())).Float64()
				ns.Element(i).Set(f)
			}
		}
	}
	return ns, nil
}
