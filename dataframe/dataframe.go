/**
  Copyright (c) [2024] [JiangNan]
  [go-tools] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package dataframe

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"gitee.com/jn-qq/go-tools/pandas/series"
	"github.com/apcera/termtables"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type DataFrame struct {
	columns []series.Series
	cols    int
	rows    int
}

// New 创建 DataFrame 数据对象
//
//	columns: 待输入数据，可以为 []any 。series.Series / []int / []float64 / []string / []bool
//	colsName: 列名，当 columns 为 series.Series 可为nil
func New(columns []any, colsName []string) (*DataFrame, error) {
	df := &DataFrame{columns: make([]series.Series, 0), cols: 0, rows: 0}
	if columns == nil {
		return df, nil
	} else if !equalLength(columns, 0) {
		return df, fmt.Errorf("columns must equal length")
	}

	for i, column := range columns {
		var ns *series.Series
		var err error
		switch value := column.(type) {
		case *series.Series:
			ns = value.Copy()
		case []bool:
			ns, err = series.NewSeries(value, series.Bool, colsName[i])
		case []float64:
			ns, err = series.NewSeries(value, series.Float, colsName[i])
		case []string:
			ns, err = series.NewSeries(value, series.String, colsName[i])
		case []int:
			ns, err = series.NewSeries(value, series.Int, colsName[i])
		}
		if err != nil {
			return nil, err
		}
		df.columns = append(df.columns, *ns)
	}

	df.Size()

	return df, nil
}

// LoadRecord 用二维字符串切片创建 DataFrame 数据对象
//
//	rows: 待输入数据
//	colsName：每列名称，当为 nil 时，rows[0]作为每列名称
//	colsType：每列数据类型
func LoadRecord(rows [][]string, colsName []string, colsType []series.Type) (*DataFrame, error) {
	// 1.检查输入数据
	if rows == nil {
		return nil, fmt.Errorf("输入数据不能为空")
	} else if !equalLength(rows, 0) {
		return nil, fmt.Errorf("各子切片数据个数不同")
	} else if len(rows[0]) != len(colsType) {
		return nil, fmt.Errorf("数据列数与类型列数不相等")
	}
	// 判断 数据列数是否与类型、名称对应
	maxCol := len(colsType)
	if colsName != nil && maxCol != len(colsName) {
		return nil, fmt.Errorf("数据列数与名称列数不相等")
	}
	if colsName == nil {
		colsName = rows[0]
		rows = rows[1:]
	}
	maxRow := len(rows)

	// 2.创建 二维表
	df := &DataFrame{
		columns: make([]series.Series, 0),
		cols:    maxCol,
		rows:    maxRow,
	}
	// 遍历每一行统计列数据
	values := make([][]string, maxCol)
	for i := 0; i < maxRow; i++ {
		for j := 0; j < maxCol; j++ {
			if values[j] == nil {
				values[j] = []string{rows[i][j]}
			} else {
				values[j] = append(values[j], rows[i][j])
			}
		}
	}
	// 生成 series.Series 对象
	for i, value := range values {
		df.columns = append(df.columns, *series.LoadRecords(value, colsType[i], colsName[i]))
	}
	return df, nil
}

// LoadMap 通过列名集合创建表
//
//	values map[string]T, T = []int、[]string、[]float64、[]bool
func LoadMap(values map[string]any) (*DataFrame, error) {
	if values == nil {
		return nil, fmt.Errorf("请输入数据")
	} else if !equalLength(values, 0) {
		return nil, fmt.Errorf("各列数据不相等")
	}
	var colsName []string
	var columns []any
	for key, value := range values {
		colsName = append(colsName, key)
		columns = append(columns, value)
	}
	df, err := New(columns, colsName)
	if err != nil {
		return nil, err
	}
	return df, nil

}

// Records 返回字符串切片
//
//	isRow = false 返回列切片 isRow = true  返回行切片
//	hasColName 是否返回列名，第一个元素，或切片
func (df *DataFrame) Records(isRow bool, hasColName bool) [][]string {
	var res [][]string
	if isRow {
		if hasColName {
			res = append(res, df.Names())
		}
		for i := 0; i < df.rows; i++ {
			var rows []string
			for _, column := range df.columns {
				rows = append(rows, column.Element(i).Records())
			}
			res = append(res, rows)
		}
	} else {
		for _, column := range df.columns {
			values := column.Records()

			if hasColName {
				values = slices.Insert(values, 0, column.Name)
			}
			res = append(res, values)
		}
	}

	return res
}

// Names 返回列名
func (df *DataFrame) Names() []string {
	var names []string
	for _, column := range df.columns {
		names = append(names, column.Name)
	}
	return names
}

// Types 返回列类型
func (df *DataFrame) Types() []string {
	var _types []string
	for _, column := range df.columns {
		_types = append(_types, column.Type())
	}
	return _types
}

// 自定义输出
func (df *DataFrame) String() string {
	return df.Print(false)
}

func (df *DataFrame) Print(isComplete bool) (str string) {

	// 创建表格对象
	table := termtables.CreateTable()
	// 添加标题
	if df.cols == 0 || df.rows == 0 {
		table.AddTitle("DataFrame Is Empty")
	} else {
		table.AddTitle(fmt.Sprintf("DataFrame Size：%d x %d", df.cols, df.rows))
	}
	// 添加表头
	colsName := slices.Insert(df.Names(), 0, "Index")
	colsType := slices.Insert(df.Types(), 0, "Types")
	var headers, dTypes []any
	for i := 0; i < df.cols+1; i++ {
		headers = append(headers, colsName[i])
		dTypes = append(dTypes, colsType[i])
	}
	table.AddHeaders(headers...)

	// 添加表内容
	for i, rows := range df.Records(true, false) {
		if !isComplete && df.rows > 50 {
			if i == 15 {
				table.AddRow(data.SliceToAny(strings.Split(strings.Repeat(".", df.cols+1), ""))...)
				table.AddRow(data.SliceToAny(strings.Split(strings.Repeat(".", df.cols+1), ""))...)
				table.AddRow(data.SliceToAny(strings.Split(strings.Repeat(".", df.cols+1), ""))...)
				continue
			} else if i > 15 && i < df.rows-6 {
				continue
			}
		}
		table.AddRow(data.SliceToAny(append([]string{strconv.Itoa(i + 1)}, rows...))...)
	}

	// 添加表脚
	table.AddSeparator()
	table.AddRow(dTypes...)

	return table.Render()
}

// Size 更新并返回二维数组大小
func (df *DataFrame) Size() (cols, rows int) {
	df.cols = len(df.columns)
	if df.cols > 0 {
		df.rows = df.columns[0].Len()
	} else {
		df.rows = 0
	}
	return df.cols, df.rows
}

func (df *DataFrame) NCols() int {
	return df.cols
}

func (df *DataFrame) NRows() int {
	return df.rows
}

// Columns 返回列
func (df *DataFrame) Columns(name string) (series.Series, error) {
	if indexCol := slices.IndexFunc(df.columns, func(s series.Series) bool { return s.Name == name }); indexCol == -1 {
		return series.Series{}, fmt.Errorf("name: %s is not found", name)
	} else {
		return df.columns[indexCol], nil
	}
}

func (df *DataFrame) SelectCols(names ...string) *DataFrame {
	var elements []any
	for _, name := range names {
		if indexCol := slices.IndexFunc(df.columns, func(s series.Series) bool { return s.Name == name }); indexCol == -1 {
			return nil
		} else {
			elements = append(elements, *df.columns[indexCol].Copy())
		}
	}

	if frame, err := New(elements, nil); err != nil {
		return nil
	} else {
		return frame
	}
}

// Rows 返回行
func (df *DataFrame) Rows(r int) map[string]series.Element {
	if r >= df.rows {
		return nil
	}
	var rows = make(map[string]series.Element)
	for _, column := range df.columns {
		rows[column.Name] = column.Element(r)
	}
	return rows
}

// Cell 返回指定单元格元素
func (df *DataFrame) Cell(r int, name string) series.Element {
	s, _ := df.Columns(name)
	return s.Element(r)
}

// Copy 复制
func (df *DataFrame) Copy() *DataFrame {
	frame := DataFrame{
		columns: slices.Clone(df.columns),
	}
	frame.Size()
	return &frame
}

// Set 设置  index 行的值
//
//	values：可选[]any、map[string]any ,[]any要更改行的所有元素
//	index < DataFrame.rows 更新行，index >= DataFrame.rows 添加行
func (df *DataFrame) Set(index int, values any) error {
	switch value := values.(type) {
	case []any:
		if len(value) != df.cols {
			return fmt.Errorf("length of columns must equal %d", df.cols)
		}
		for i, v := range value {
			if index >= df.rows {
				if err := df.columns[i].Append(v); err != nil {
					return err
				}
			} else {
				df.columns[i].Element(index).Set(v)
			}
		}
	case map[string]any:
		for k, v := range value {
			i := slices.IndexFunc(df.columns, func(s series.Series) bool { return s.Name == k })
			if i == -1 {
				return fmt.Errorf("column %s not found", k)
			}
			if index >= df.rows {
				if err := df.columns[i].Append(v); err != nil {
					return err
				}
			} else {
				df.columns[i].Element(index).Set(v)
			}
		}
	default:
		return fmt.Errorf("type %T not supported", value)
	}
	df.Size()
	return nil
}

// AddRows 向列表末尾添加行
func (df *DataFrame) AddRows(values [][]any) error {
	if !equalLength(values, df.cols) {
		return fmt.Errorf("length of columns must equal %d", df.cols)
	}
	for _, value := range values {
		if err := df.Set(df.rows, value); err != nil {
			return err
		}
	}
	return nil
}

// AddCol 添加列
//
//	name：列名。如果已存在更新，否则添加
//	values：可选 series.Series []E {int | float64 | bool | string}
//	defaultValue：当 values 长度不足时，自动添加
func (df *DataFrame) AddCol(name string, values any, defaultValue any) error {
	var ns *series.Series
	switch value := values.(type) {
	case *series.Series:
		// 补长度
		if value.Len() < df.rows {
			if defaultValue == nil {
				return fmt.Errorf("length of default series must equal %d", df.rows)
			}
			if err := value.Append(data.CreateSlice(defaultValue, df.rows-value.Len())); err != nil {
				return err
			}
		} else if value.Len() > df.rows {
			value, _ = value.SubSet(data.Range(0, df.rows, 1)...)
		}
		ns = value.Copy()
	case []int:
		if len(value) < df.rows {
			value = append(value, data.CreateSlice(defaultValue.(int), df.rows-len(value))...)
		}
		ns, _ = series.NewSeries(value, series.Int, name)
	case []string:
		if len(value) < df.rows {
			value = append(value, data.CreateSlice(defaultValue.(string), df.rows-len(value))...)
		}
		ns, _ = series.NewSeries(value, series.String, name)
	case []float64:
		if len(value) < df.rows {
			value = append(value, data.CreateSlice(defaultValue.(float64), df.rows-len(value))...)
		}
		ns, _ = series.NewSeries(value, series.Float, name)
	case []bool:
		if len(value) < df.rows {
			value = append(value, data.CreateSlice(defaultValue.(bool), df.rows-len(value))...)
		}
		ns, _ = series.NewSeries(value[:df.rows], series.Bool, name)
	default:
		return fmt.Errorf("type %T not supported", value)
	}

	// 更新或添加
	indexCol := slices.IndexFunc(df.columns, func(s series.Series) bool { return s.Name == ns.Name })
	if indexCol == -1 {
		df.columns = append(df.columns, *ns)
	} else {
		df.columns[indexCol] = *ns
	}
	df.Size()
	return nil
}

// Concat 合并两个表
//
//	isColumn：是否合并在右侧 ，如果两个表列名相同，则更新原表列
func (df *DataFrame) Concat(x DataFrame, isColumn bool) error {
	if isColumn && df.rows != x.rows {
		return fmt.Errorf("rows must equal %d", df.rows)
	} else if !isColumn && df.cols != x.cols {
		return fmt.Errorf("columns must equal %d", df.cols)
	}
	if isColumn {
		for _, column := range x.columns {
			if err := df.AddCol("", column, nil); err != nil {
				return err
			}
		}
	} else {
		n1, n2 := df.Names(), x.Names()
		slices.Sort(n1)
		slices.Sort(n2)
		if !slices.Equal(n1, n2) {
			return fmt.Errorf("两个表列名不相同")
		}
		for i, name := range df.Names() {
			ns, err := x.Columns(name)
			if err != nil {
				return err
			}
			if err := df.columns[i].Append(*ns.Copy()); err != nil {
				return err
			}
		}
		df.Size()
	}
	return nil
}

// DropCols 批量删除
func (df *DataFrame) DropCols(names ...string) {
	df.columns = slices.DeleteFunc(df.columns, func(s series.Series) bool { return slices.Contains(names, s.Name) })
	df.Size()
}

// Rename 批量命名
func (df *DataFrame) Rename(cols map[string]string) {
	for _, column := range df.columns {
		if value, ok := cols[column.Name]; ok {
			column.Name = value
		} else {
			fmt.Printf("column %s not found", column.Name)
		}
	}
}

// Arrange 排序
func (df *DataFrame) Arrange(order ...Order) error {
	frame := df.Copy()
	for _, o := range order {
		i := slices.IndexFunc(frame.Names(), func(s string) bool { return s == o.ColumnName })
		if i == -1 {
			return fmt.Errorf("column %s not found", o.ColumnName)
		}
		indexes := frame.columns[i].SortIndex(o.Reverse)
		for i := 0; i < frame.cols; i++ {
			ns, _ := frame.columns[i].SubSet(indexes...)
			ns.InitIndex()
			frame.columns[i] = *ns
		}
	}
	df.columns = frame.columns
	return nil
}

func (df *DataFrame) SubSet(indexes ...int) (*DataFrame, error) {
	frame := df.Copy()
	for i := 0; i < frame.cols; i++ {
		set, err := frame.columns[i].SubSet(indexes...)
		if err != nil {
			return nil, err
		}
		frame.columns[i] = *set
	}
	frame.Size()
	return frame, nil
}

// Filter 过滤
func (df *DataFrame) Filter(filters ...F) (*DataFrame, error) {
	var indexes []int
	for _, filter := range filters {
		ns, err := df.Columns(filter.Column)
		if err != nil {
			return nil, err
		}
		if s, err := ns.Filter(filter.Operator, filter.values); err != nil {
			return nil, err
		} else {
			if indexes == nil {
				indexes = s.Indexes()
			} else if filter.OR {
				indexes = slices.Compact(append(indexes, s.Indexes()...))
			} else {
				indexes = data.Overlap(indexes, s.Indexes())
			}
		}
	}
	set, err := df.SubSet(indexes...)
	if err != nil {
		return nil, err
	}
	return set, nil
}

// Order 排序结构
type Order struct {
	// 列名
	ColumnName string
	// 倒叙
	Reverse bool
}

// SortByForward 正序
func SortByForward(name string) Order {
	return Order{ColumnName: name, Reverse: false}
}

// SortByReverse 倒叙
func SortByReverse(name string) Order {
	return Order{ColumnName: name, Reverse: true}
}

// F 过滤条件
//
//	OR与前一个过滤的关系
type F struct {
	Column   string
	Operator series.RelationalOperator
	values   any
	OR       bool
}

// 判断切片长度是否相等
func equalLength(values any, baseLen int) bool {
	value := reflect.ValueOf(values)
	if reflect.TypeOf(values).Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			var l int
			if value.Index(i).Type().String() == "series.Series" {
				l = int(value.Index(i).MethodByName("Len").Int())
			} else {
				l = reflect.ValueOf(value.Index(i).Interface()).Len()
			}

			if baseLen == 0 {
				baseLen = l
			}

			if baseLen != l {
				return false
			}
		}
	} else if reflect.TypeOf(values).Kind() == reflect.Map {
		it := value.MapRange()
		for it.Next() {
			if baseLen == 0 {
				baseLen = reflect.ValueOf(it.Value().Interface()).Len()
			}
			if reflect.ValueOf(it.Value().Interface()).Len() != baseLen {
				return false
			}
		}
	}

	return true
}

func (df *DataFrame) Groups(names ...string) (map[string]*DataFrame, error) {
	if names == nil {
		return nil, fmt.Errorf("names is nil")
	}
	var ns *series.Series
	for i, name := range names {
		ns1, err := df.Columns(name)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			ns = ns1.Copy()
			if err = ns.SetType(series.String); err != nil {
				return nil, err
			}
			ns.Name = "_" + ns.Name
		} else {
			ns, err = ns.Arithmetic(series.Addition, ns1)
			if err != nil {
				return nil, err
			}
			ns.Name += "_" + ns1.Name
		}
	}

	frame := df.Copy()
	err := frame.AddCol("", ns, "")
	if err != nil {
		return nil, err
	}

	group := make(map[string]*DataFrame, 0)
	for _, s := range slices.Compact(ns.Records()) {
		df1, err := frame.Filter(F{
			Column:   ns.Name,
			Operator: series.Equal,
			values:   s,
			OR:       false,
		})
		if err != nil {
			return nil, err
		}
		df1.DropCols(ns.Name)
		group[s] = df1
	}

	return group, nil
}

func (df *DataFrame) FormatCols(f func(index int, elem series.Element) series.Element, cols ...string) error {
	for _, col := range cols {
		column, err := df.Columns(col)
		if err != nil {
			return err
		}
		column.Format(f)
	}
	return nil
}
