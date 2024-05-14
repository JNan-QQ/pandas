package dataframe

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"gitee.com/jn-qq/go-tools/pandas/series"
	"strconv"
)

func ExampleLoadMap() {
	df, _ := LoadMap(map[string]any{
		"name":  data.CreateSlice("Join", 100),
		"phone": data.CreateSlice("15963578965", 100),
	})
	fmt.Println(df)
	// Output:	+------------------------------+
	//|   DataFrame Size：2 x 100    |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join   | 15963578965 |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//| 6     | Join   | 15963578965 |
	//| 7     | Join   | 15963578965 |
	//| 8     | Join   | 15963578965 |
	//| 9     | Join   | 15963578965 |
	//| 10    | Join   | 15963578965 |
	//| 11    | Join   | 15963578965 |
	//| 12    | Join   | 15963578965 |
	//| 13    | Join   | 15963578965 |
	//| 14    | Join   | 15963578965 |
	//| 15    | Join   | 15963578965 |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| 95    | Join   | 15963578965 |
	//| 96    | Join   | 15963578965 |
	//| 97    | Join   | 15963578965 |
	//| 98    | Join   | 15963578965 |
	//| 99    | Join   | 15963578965 |
	//| 100   | Join   | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | string      |
	//+-------+--------+-------------+

}

func ExampleLoadRecord() {
	df, _ := LoadRecord(
		data.CreateSlice([]string{"Join", "15963578965"}, 100),
		[]string{"name", "phone"},
		[]series.Type{series.String, series.Int},
	)
	fmt.Println(df)
	// Output:	+------------------------------+
	//|   DataFrame Size：2 x 100    |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join   | 15963578965 |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//| 6     | Join   | 15963578965 |
	//| 7     | Join   | 15963578965 |
	//| 8     | Join   | 15963578965 |
	//| 9     | Join   | 15963578965 |
	//| 10    | Join   | 15963578965 |
	//| 11    | Join   | 15963578965 |
	//| 12    | Join   | 15963578965 |
	//| 13    | Join   | 15963578965 |
	//| 14    | Join   | 15963578965 |
	//| 15    | Join   | 15963578965 |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| 95    | Join   | 15963578965 |
	//| 96    | Join   | 15963578965 |
	//| 97    | Join   | 15963578965 |
	//| 98    | Join   | 15963578965 |
	//| 99    | Join   | 15963578965 |
	//| 100   | Join   | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}

func ExampleNew() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 100), data.CreateSlice(15963578965, 100)},
		[]string{"name", "phone"},
	)
	fmt.Println(df)
	//	Output:	+------------------------------+
	//|   DataFrame Size：2 x 100    |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join   | 15963578965 |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//| 6     | Join   | 15963578965 |
	//| 7     | Join   | 15963578965 |
	//| 8     | Join   | 15963578965 |
	//| 9     | Join   | 15963578965 |
	//| 10    | Join   | 15963578965 |
	//| 11    | Join   | 15963578965 |
	//| 12    | Join   | 15963578965 |
	//| 13    | Join   | 15963578965 |
	//| 14    | Join   | 15963578965 |
	//| 15    | Join   | 15963578965 |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| .     | .      | .           |
	//| 95    | Join   | 15963578965 |
	//| 96    | Join   | 15963578965 |
	//| 97    | Join   | 15963578965 |
	//| 98    | Join   | 15963578965 |
	//| 99    | Join   | 15963578965 |
	//| 100   | Join   | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}

func ExampleDataFrame_Records() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	fmt.Println(df.Records(false, false))
	fmt.Println(df.Records(false, true))
	fmt.Println(df.Records(true, true))
	// Output:[[Join ... Join] [15963578965 ... 15963578965]]
	//[[name Join ... Join] [phone 15963578965 ... 15963578965]]
	//[[name phone] [Join 15963578965] ... [Join 15963578965]]
}

func ExampleDataFrame_Set() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	_ = df.Set(0, []any{"Andy", 1111111111})
	fmt.Println(df)
	_ = df.Set(2, map[string]any{"phone": 222222222})
	fmt.Println(df)
	_ = df.Set(df.rows, []any{"Andy", 1111111111})
	fmt.Println(df)
	// output:+------------------------------+
	//|    DataFrame Size：2 x 5     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Andy   | 1111111111  |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+
	//
	//+------------------------------+
	//|    DataFrame Size：2 x 5     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Andy   | 1111111111  |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 222222222   |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+
	//
	//+------------------------------+
	//|    DataFrame Size：2 x 6     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Andy   | 1111111111  |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 222222222   |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//| 6     | Andy   | 1111111111  |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}

func ExampleDataFrame_AddRows() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	_ = df.AddRows([][]any{{"name1", 12345678}, {"name1", 12345678}, {"name1", 12345678}})
	fmt.Println(df)
	// output:+------------------------------+
	//|    DataFrame Size：2 x 8     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join   | 15963578965 |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Join   | 15963578965 |
	//| 5     | Join   | 15963578965 |
	//| 6     | name1  | 12345678    |
	//| 7     | name1  | 12345678    |
	//| 8     | name1  | 12345678    |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}

func ExampleDataFrame_Arrange() {
	df, _ := New(
		[]any{
			[]string{"伏旭歆", "管原炳", "仰芝凤", "万茵瑾", "左芊筱", "俞淑允", "宗茹淳", "卓虹", "司丽瑾", "岑泳继"},
			[]int{13935531105, 15665203778, 14583084372, 14779318181, 17606363473, 18950385204, 18659058185, 16628908658, 17590257481, 17254554855},
			[]int{35, 36, 42, 13, 20, 20, 14, 20, 30, 36},
		},
		[]string{"name", "phone", "age"},
	)
	_ = df.Arrange(SortByForward("name"), Order{ColumnName: "age", Reverse: true})
	fmt.Println(df)
	// output:+------------------------------------+
	//|       DataFrame Size：3 x 10       |
	//+-------+--------+-------------+-----+
	//| Index | name   | phone       | age |
	//+-------+--------+-------------+-----+
	//| 1     | 仰芝凤 | 14583084372 | 42  |
	//| 2     | 岑泳继 | 17254554855 | 36  |
	//| 3     | 管原炳 | 15665203778 | 36  |
	//| 4     | 伏旭歆 | 13935531105 | 35  |
	//| 5     | 司丽瑾 | 17590257481 | 30  |
	//| 6     | 卓虹   | 16628908658 | 20  |
	//| 7     | 左芊筱 | 17606363473 | 20  |
	//| 8     | 俞淑允 | 18950385204 | 20  |
	//| 9     | 宗茹淳 | 18659058185 | 14  |
	//| 10    | 万茵瑾 | 14779318181 | 13  |
	//+-------+--------+-------------+-----+
	//| Types | string | int         | int |
	//+-------+--------+-------------+-----+

}

func ExampleDataFrame_AddCol() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	_ = df.AddCol("addr", data.CreateSlice("xxxxx", 5), nil)
	fmt.Println(df)
	_ = df.AddCol("addr", data.CreateSlice("yyyy", 1), "yyy")
	fmt.Println(df)
	// output:+---------------------------------------+
	//|         DataFrame Size：3 x 5         |
	//+-------+--------+-------------+--------+
	//| Index | name   | phone       | addr   |
	//+-------+--------+-------------+--------+
	//| 1     | Join   | 15963578965 | xxxxx  |
	//| 2     | Join   | 15963578965 | xxxxx  |
	//| 3     | Join   | 15963578965 | xxxxx  |
	//| 4     | Join   | 15963578965 | xxxxx  |
	//| 5     | Join   | 15963578965 | xxxxx  |
	//+-------+--------+-------------+--------+
	//| Types | string | int         | string |
	//+-------+--------+-------------+--------+
	//
	//+---------------------------------------+
	//|         DataFrame Size：3 x 5         |
	//+-------+--------+-------------+--------+
	//| Index | name   | phone       | addr   |
	//+-------+--------+-------------+--------+
	//| 1     | Join   | 15963578965 | yyyy   |
	//| 2     | Join   | 15963578965 | yyy    |
	//| 3     | Join   | 15963578965 | yyy    |
	//| 4     | Join   | 15963578965 | yyy    |
	//| 5     | Join   | 15963578965 | yyy    |
	//+-------+--------+-------------+--------+
	//| Types | string | int         | string |
	//+-------+--------+-------------+--------+

}

func ExampleDataFrame_Concat() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 3), data.CreateSlice(15963578965, 3)},
		[]string{"name", "phone"},
	)
	df1, _ := New(
		[]any{data.CreateSlice("Mary", 3), data.CreateSlice(19645698705, 3)},
		[]string{"name", "phone"},
	)
	_ = df.Concat(*df1, false)
	fmt.Println(df)
	// output:+------------------------------+
	//|    DataFrame Size：2 x 6     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join   | 15963578965 |
	//| 2     | Join   | 15963578965 |
	//| 3     | Join   | 15963578965 |
	//| 4     | Mary   | 19645698705 |
	//| 5     | Mary   | 19645698705 |
	//| 6     | Mary   | 19645698705 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}

func ExampleDataFrame_Groups() {
	df, _ := New(
		[]any{
			[]string{"喻靖元", "尤淇方", "方文栋", "郝晨轩", "养海露", "弘展鹏", "滕安平", "谷灵雁", "陶海露", "乔瀚天"},
			[]string{"男", "男", "男", "男", "女", "男", "男", "女", "女", "男"},
			[]int{51, 29, 44, 21, 26, 29, 68, 21, 29, 52},
		},
		[]string{"name", "sex", "age"},
	)
	group, _ := df.Groups("sex")
	for s, frame := range group {
		fmt.Println("group_name is " + s)
		fmt.Println(frame)
	}
	// output:
	//group_name is 男
	//+-------------------------------+
	//|     DataFrame Size：3 x 7     |
	//+-------+--------+--------+-----+
	//| Index | name   | sex    | age |
	//+-------+--------+--------+-----+
	//| 1     | 喻靖元 | 男     | 51  |
	//| 2     | 尤淇方 | 男     | 29  |
	//| 3     | 方文栋 | 男     | 44  |
	//| 4     | 郝晨轩 | 男     | 21  |
	//| 5     | 弘展鹏 | 男     | 29  |
	//| 6     | 滕安平 | 男     | 68  |
	//| 7     | 乔瀚天 | 男     | 52  |
	//+-------+--------+--------+-----+
	//| Types | string | string | int |
	//+-------+--------+--------+-----+
	//
	//group_name is 女
	//+-------------------------------+
	//|     DataFrame Size：3 x 3     |
	//+-------+--------+--------+-----+
	//| Index | name   | sex    | age |
	//+-------+--------+--------+-----+
	//| 1     | 养海露 | 女     | 26  |
	//| 2     | 谷灵雁 | 女     | 21  |
	//| 3     | 陶海露 | 女     | 29  |
	//+-------+--------+--------+-----+
	//| Types | string | string | int |
	//+-------+--------+--------+-----+

}

func ExampleDataFrame_Groups_two() {
	df, _ := New(
		[]any{
			[]string{"喻靖元", "尤淇方", "方文栋", "郝晨轩", "养海露", "弘展鹏", "滕安平", "谷灵雁", "陶海露", "乔瀚天"},
			[]string{"男", "男", "男", "男", "女", "男", "男", "女", "女", "男"},
			[]int{51, 29, 44, 21, 26, 29, 68, 21, 29, 52},
		},
		[]string{"name", "sex", "age"},
	)
	group, _ := df.Groups("sex", "age")
	for s, frame := range group {
		fmt.Println("group_name is " + s)
		fmt.Println(frame)
	}
	// output:
	//group_name is 男51
	//+-------------------------------+
	//|     DataFrame Size：3 x 1     |
	//+-------+--------+--------+-----+
	//| Index | name   | sex    | age |
	//+-------+--------+--------+-----+
	//| 1     | 喻靖元 | 男     | 51  |
	//+-------+--------+--------+-----+
	//| Types | string | string | int |
	//+-------+--------+--------+-----+
	//
	//	......
	//
	//group_name is 男29
	//+-------------------------------+
	//|     DataFrame Size：3 x 2     |
	//+-------+--------+--------+-----+
	//| Index | name   | sex    | age |
	//+-------+--------+--------+-----+
	//| 1     | 尤淇方 | 男     | 29  |
	//| 2     | 弘展鹏 | 男     | 29  |
	//+-------+--------+--------+-----+
	//| Types | string | string | int |
	//+-------+--------+--------+-----+
	//
	//group_name is 男21
	//+-------------------------------+
	//|     DataFrame Size：3 x 1     |
	//+-------+--------+--------+-----+
	//| Index | name   | sex    | age |
	//+-------+--------+--------+-----+
	//| 1     | 郝晨轩 | 男     | 21  |
	//+-------+--------+--------+-----+
	//| Types | string | string | int |
	//+-------+--------+--------+-----+

}

func ExampleReadXLSX() {
	dfs, err := ReadXLSX("test.xlsx",
		Sheets{
			SheetName: "Sheet3",
			ColsType:  []series.Type{series.String, series.Int},
		},
	)
	if err != nil {
		panic(err)
	}
	for s, frame := range dfs {
		fmt.Println("Sheet Name is", s)
		fmt.Println(frame)
	}
}

func ExampleReadXLSX_a() {
	dfs, err := ReadXLSX("test.xlsx",
		Sheets{
			SheetName: "Sheet3",
			SRow:      5,
			SCol:      2,
			ColsType:  []series.Type{series.String, series.Int},
		},
	)
	if err != nil {
		panic(err)
	}
	for s, frame := range dfs {
		fmt.Println("Sheet Name is", s)
		fmt.Println(frame)
	}
}

func ExampleReadCSV() {
	readCSV, err := ReadCSV("test.csv", Sheets{
		ColsType: []series.Type{series.String, series.Int},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(readCSV)
}

func ExampleDataFrame_WriteToXLSX() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 100), data.CreateSlice(15963578965, 100)},
		[]string{"name", "phone"},
	)
	err := df.WriteToXLSX("test.xlsx", "Sheet3")
	if err != nil {
		panic(err)
	}
}

func ExampleDataFrame_WriteToCSV() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 100), data.CreateSlice(15963578965, 100)},
		[]string{"name", "phone"},
	)
	err := df.WriteToCSV("test.csv")
	if err != nil {
		panic(err)
	}
}

func ExampleDataFrame_FormatCols() {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	err := df.FormatCols(func(index int, elem series.Element) series.Element {
		elem.Set(elem.Records() + strconv.Itoa(index))
		return elem
	}, "name")
	if err != nil {
		panic(err)
	}
	fmt.Println(df)
	//output:+------------------------------+
	//|    DataFrame Size：2 x 5     |
	//+-------+--------+-------------+
	//| Index | name   | phone       |
	//+-------+--------+-------------+
	//| 1     | Join0  | 15963578965 |
	//| 2     | Join1  | 15963578965 |
	//| 3     | Join2  | 15963578965 |
	//| 4     | Join3  | 15963578965 |
	//| 5     | Join4  | 15963578965 |
	//+-------+--------+-------------+
	//| Types | string | int         |
	//+-------+--------+-------------+

}
