/**
  Copyright (c) [2024] [JiangNan]
  [pandas] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package dataframe

import (
	"encoding/csv"
	"fmt"
	"gitee.com/jn-qq/go-tools/pandas/series"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
)

const (
	XLSX int = iota
	CSV
)

// Sheets 数据对象
type Sheets struct {
	SCol      int           // 开始列, 默认 1
	SRow      int           // 开始行, 默认 1
	ECol      int           // 结束列, 默认最后一列
	ERow      int           // 结束行, 默认最后一行
	Header    []string      // 表头，默认第一行
	SheetName string        // 工作部名 XLSX 特有
	ColsType  []series.Type // 列类型
}

// ReadXLSX 从XLSX中读取表格
func ReadXLSX(filePath string, sheets ...Sheets) (map[string]*DataFrame, error) {
	// 读取文档
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	frames := make(map[string]*DataFrame)
	// 遍历sheet
	for _, sheet := range sheets {
		// 初始化开始行列
		if sheet.SCol == 0 {
			sheet.SCol = 1
		}
		if sheet.SRow == 0 {
			sheet.SRow = 1
		}
		rows, err := f.Rows(sheet.SheetName)
		if err != nil {
			return nil, err
		}
		// 表头，数据
		header := sheet.Header
		var sheetData [][]string
		// 遍历行
		row := 1
		for rows.Next() {
			if row < sheet.SRow || (sheet.ERow != 0 && row > sheet.ERow) {
				row += 1
				continue
			} else {
				row += 1
			}
			// 行数据
			columns, err := rows.Columns()
			if err != nil {
				return nil, err
			}

			if sheet.ECol == 0 {
				sheet.ECol = len(columns)
			}

			if header == nil {
				// 表头
				header = columns[sheet.SCol-1 : sheet.ECol]
			} else {
				// 表数据
				sheetData = append(sheetData, columns[sheet.SCol-1:sheet.ECol])
			}
		}

		record, err := LoadRecord(sheetData, header, sheet.ColsType)
		if err != nil {
			return nil, err
		}
		frames[sheet.SheetName] = record
	}
	return frames, nil
}

// ReadCSV 从CSV中读取表格
func ReadCSV(filePath string, sheet Sheets) (*DataFrame, error) {
	// 初始化开始行列
	if sheet.SCol == 0 {
		sheet.SCol = 1
	}
	if sheet.SRow == 0 {
		sheet.SRow = 1
	}
	//打开文件(只读模式)，创建io.read接口实例
	opencast, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := opencast.Close(); err != nil {
			panic(err)
		}
	}()
	// 创建csv对象
	reader := csv.NewReader(opencast)
	header := sheet.Header
	var sheetData [][]string
	row := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if row < sheet.SRow || (sheet.ERow != 0 && row > sheet.ERow) {
			row += 1
			continue
		} else {
			row += 1
		}

		if sheet.ECol == 0 {
			sheet.ECol = len(record)
		}

		if header == nil {
			header = record[sheet.SCol-1 : sheet.ECol]
		} else {
			sheetData = append(sheetData, record[sheet.SCol-1:sheet.ECol])
		}
	}
	record, err := LoadRecord(sheetData, header, sheet.ColsType)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (df *DataFrame) WriteToCSV(p string) error {
	newFile, err := os.Create(p)
	if err != nil {
		return err
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			return
		}
	}()
	// 写入UTF-8 BOM，防止中文乱码
	if _, err := newFile.WriteString("\xEF\xBB\xBF"); err != nil {
		return err
	}
	// 写数据到csv文件
	w := csv.NewWriter(newFile)
	// WriteAll方法使用Write方法向w写入多条记录，并在最后调用Flush方法清空缓存。
	if err := w.WriteAll(df.Records(true, true)); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func (df *DataFrame) WriteToXLSX(p, sheetName string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	if err := f.SetDefaultFont("楷体"); err != nil {
		return err
	}
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	// 写入表头
	hd := df.Names()
	if err := f.SetSheetRow(sheetName, "A1", &hd); err != nil {
		return err
	}
	// 设置表头格式
	h1, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	_ = f.SetRowStyle(sheetName, 1, 1, h1)

	// 按列写入
	for i, column := range df.columns {
		d1 := column.Any()
		if err = f.SetSheetCol(sheetName, fmt.Sprintf("%s2", string(rune(65+i))), &d1); err != nil {
			return err
		}
	}

	// 设置内容格式
	d1, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 12,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
	})
	_ = f.SetRowStyle(sheetName, 2, df.rows+1, d1)

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err = f.SaveAs(p); err != nil {
		return err
	}

	return nil
}
