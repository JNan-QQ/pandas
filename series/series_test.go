package series

import (
	"fmt"
)

func ExampleNewSeries() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	fmt.Println(s1)
	//	output:
	//字段名：number1
	//数 据：[1 2 test 2 4 6]
	//索 引：[0 1 2 3 4 5]
	//类 型：string
}

func ExampleNewSeries_a() {
	s1, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number2")
	fmt.Println(s1)
	//	output:
	//字段名：number2
	//数 据：[1 2 3 4 6]
	//索 引：[0 1 2 3 4]
	//类 型：int
}

func ExampleNewSeries_err() {
	s1, err := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, Float, "number3")
	fmt.Println(s1, err)
	//	output:
	//<nil> 输入切片与指定数据类型不匹配
}

func ExampleSeries_Append() {
	s1, _ := NewSeries(
		[]string{"1", "2", "test", "2", "4", "6"},
		String,
		"number1",
	)
	_ = s1.Append([]string{"1", "2", "test", "2", "4", "6"})
	fmt.Println(s1)
	//	output:
	//字段名：number1
	//数 据：[1 2 test 2 4 6 1 2 test 2 4 6]
	//索 引：[0 1 2 3 4 5 6 7 8 9 10 11]
	//类 型：string
}

func ExampleSeries_Append_a() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	_ = s1.Append(1)
	fmt.Println(s1)
	//	output:
	//字段名：number1
	//数 据：[1 2 test 2 4 6 1]
	//索 引：[0 1 2 3 4 5 6]
	//类 型：string
}

func ExampleSeries_SetType() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	fmt.Println(s1)
	_ = s1.SetType(Int)
	fmt.Println(s1)
	//	output:
	//	字段名：number1
	//数 据：[1 2 test 2 4 6]
	//索引：[0 1 2 3 4 5]
	//类 型：string
	//
	//test 不能转换为 Int，已置为无限小
	//字段名：number1
	//数 据：[1 2 NaN 2 4 6]
	//索 引：[0 1 2 3 4 5]
	//类 型：int
}

func ExampleSeries_Concat() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	s2, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number2")
	_ = s1.Concat(*s2)
	fmt.Println(s1)
	// Output:
	//两个数据类型不同，正在尝试转换... done
	//字段名：number1
	//数 据：[1 2 test 2 4 6 1 2 3 4 6]
	//索 引：[0 1 2 3 4 5 6 7 8 9 10]
	//类 型：string

}

func ExampleSeries_Format() {
	s1, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number1")
	// value + index
	s1.Format(func(index int, elem Element) Element {
		e := elem.Int()
		elem.Set(e + index)
		return elem
	})

	fmt.Println(s1)

	// Output: 字段名：number1
	//数 据：[1 3 5 7 10]
	//索 引：[0 1 2 3 4]
	//类 型：int
}

func ExampleSeries_Format_a() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number2")
	s1.Format(func(index int, elem Element) Element {
		e := elem.Records()
		elem.Set(e + "ioc")
		return elem
	})

	fmt.Println(s1)

	// Output:
	//字段名：number2
	//数 据：[1ioc 2ioc testioc 2ioc 4ioc 6ioc]
	//索 引：[0 1 2 3 4 5]
	//类 型：string
}

func ExampleSeries_Filter() {
	s1, _ := NewSeries([]string{"aee", "2", "test", "2", "any", "6"}, String, "number1")
	ns, _ := s1.Filter(Equal, "2")
	fmt.Println(ns)
	fmt.Println("-----------------------------")
	ns1, _ := s1.Filter(NotEqual, "2")
	fmt.Println(ns1)
	fmt.Println("-----------------------------")
	ns2, _ := s1.Filter(Contains, "e")
	fmt.Println(ns2)
	fmt.Println("-----------------------------")
	ns3, _ := s1.Filter(StartsWith, "a")
	fmt.Println(ns3)
	fmt.Println("-----------------------------")
	ns4, _ := s1.Filter(EndsWith, "e")
	fmt.Println(ns4)

	// output:
	//字段名：number1
	//数 据：[2 2]
	//索 引：[1 3]
	//类 型：string
	//
	//-----------------------------
	//字段名：number1
	//数 据：[aee test any 6]
	//索 引：[0 2 4 5]
	//类 型：string
	//
	//-----------------------------
	//字段名：number1
	//数 据：[aee test]
	//索 引：[0 2]
	//类 型：string
	//
	//-----------------------------
	//字段名：number1
	//数 据：[aee any]
	//索 引：[0 4]
	//类 型：string
	//
	//-----------------------------
	//字段名：number1
	//数 据：[aee]
	//索 引：[0]
	//类 型：string
}

func ExampleSeries_Filter_a() {
	s1, _ := NewSeries([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 3, 100, 200, 300, 3, 8, 4}, Int, "number1")
	ns, _ := s1.Filter(Equal, 3)
	fmt.Println(ns)
	fmt.Println("-----------------------------")
	ns1, _ := s1.Filter(LessThan, 6)
	fmt.Println(ns1)
	fmt.Println("-----------------------------")
	ns2, _ := s1.Filter(LessOrEqual, 8)
	fmt.Println(ns2)
	fmt.Println("-----------------------------")
	ns3, _ := s1.Filter(GreaterOrEqual, 10)
	fmt.Println(ns3)
	fmt.Println("-----------------------------")
	ns4, _ := s1.Filter(In, []int{3, 8, 200})
	fmt.Println(ns4)

	// output:
	//字段名：number1
	//数 据：[3 3 3]
	//索 引：[2 10 14]
	//类 型：int
	//
	//-----------------------------
	//字段名：number1
	//数 据：[1 2 3 4 5 3 3 4]
	//索 引：[0 1 2 3 4 10 14 16]
	//类 型：int
	//
	//-----------------------------
	//字段名：number1
	//数 据：[1 2 3 4 5 6 7 8 3 3 8 4]
	//索 引：[0 1 2 3 4 5 6 7 10 14 15 16]
	//类 型：int
	//
	//-----------------------------
	//字段名：number1
	//数 据：[10 100 200 300]
	//索 引：[9 11 12 13]
	//类 型：int
	//
	//-----------------------------
	//字段名：number1
	//数 据：[3 8 3 200 3 8]
	//索 引：[2 7 10 12 14 15]
	//类 型：int
}

func ExampleSeries_Copy() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	s2 := s1.Copy()
	fmt.Println(&s1 == &s2)
	//	output:
	//false
}

func ExampleSeries_SubSet() {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	ns, _ := s1.SubSet(1, 2, 3)
	fmt.Println(ns)
	//	output:
	//	字段名：number1
	//数 据：[2 test 2]
	//索 引：[1 2 3]
	//类 型：string
}

func ExampleLoadRecords() {
	s1 := LoadRecords([]string{"1", "2", "test", "2", "4", "6"}, Int, "number1")
	fmt.Println(s1)
	//	output:
	//	字段名：number1
	//数 据：[1 2 NaN 2 4 6]
	//索 引：[0 1 2 3 4 5]
	//类 型：int
}

func ExampleSeries_SortIndex() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	fmt.Println(s1)
	fmt.Println(s1.SortIndex(false))
	fmt.Println(s1.SortIndex(true))
	//	output:
	//	字段名：number1
	//数 据：[10 1 5 3 7 5 1 36 5]
	//索 引：[0 1 2 3 4 5 6 7 8]
	//类 型：int
	//
	//[1 6 3 2 5 8 4 0 7]
	//[7 0 4 2 5 8 3 1 6]
}

func ExampleSeries_Drop() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	s1 = s1.Drop(3, 6, 4)
	fmt.Println(s1)
	//	output:
	//字段名：number1
	//数 据：[10 1 5 5 36 5]
	//索 引：[0 1 2 5 7 8]
	//类 型：int

}

func ExampleSeries_Arithmetic_add() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	s2, _ := NewSeries([]float64{10.1, 1.0, 5.5, 3.3, 7.0, 5.0, 0, 36.0, 5.0}, Float, "number1")
	s, _ := s1.Arithmetic(Addition, *s2)
	fmt.Println(s)
	// output:字段名：number1
	//数 据：[20.1 2 10.5 6.3 14 10 1 72 10]
	//索 引：[]
	//类 型：float64
}

func ExampleSeries_Arithmetic_sub() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	s2, _ := NewSeries([]float64{10.1, 1.0, 5.5, 3.3, 7.0, 5.0, 0, 36.0, 5.0}, Float, "number1")
	s, _ := s1.Arithmetic(Subtraction, *s2)
	fmt.Println(s)
	// output:
	//字段名：number1
	//数 据：[-0.1 0 -0.5 -0.3 0 0 1 0 0]
	//索 引：[]
	//类 型：float64

}
func ExampleSeries_Arithmetic_mul() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	s2, _ := NewSeries([]float64{10.1, 1.0, 5.5, 3.3, 7.0, 5.0, 0, 36.0, 5.0}, Float, "number1")
	s, _ := s1.Arithmetic(Multiplication, *s2)
	fmt.Println(s)
	// output:
	//字段名：number1
	//数 据：[101 1 27.5 9.9 49 25 0 1296 25]
	//索 引：[]
	//类 型：float64

}
func ExampleSeries_Arithmetic_div() {
	s1, _ := NewSeries([]int{20, 10, 5, 9, 25}, Int, "number1")
	s2, _ := NewSeries([]float64{2.5, 2, 5, 0, 12.5}, Float, "number1")
	s, _ := s1.Arithmetic(Division, *s2)
	fmt.Println(s)
	// output:字段名：number1
	//数 据：[8 5 1 NaN 2]
	//索 引：[]
	//类 型：float64

}
func ExampleSeries_Arithmetic_mod() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7}, Int, "number1")
	s2, _ := NewSeries([]float64{3, 1.0, 3, 0.5, 2}, Float, "number1")
	s, _ := s1.Arithmetic(Remainder, *s2)
	fmt.Println(s)
	// output:字段名：number1
	//数 据：[1 0 2 0 1]
	//索 引：[]
	//类 型：float64
}

func ExampleSeries_Arithmetic_a2() {
	s1, _ := NewSeries([]int{10, 1, 5, 3, 7, 5, 1, 36, 5}, Int, "number1")
	s2, _ := NewSeries([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, String, "name2")
	s, _ := s2.Arithmetic(Addition, *s1)
	fmt.Println(s)
	//output:字段名：name2
	//数 据：[a10 b1 c5 d3 e7 f5 g1 h36 i5]
	//索 引：[]
	//类 型：string

}
