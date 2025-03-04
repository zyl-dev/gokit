package utils

import (
	"fmt"
	excelize "github.com/xuri/excelize/v2"
	"io"
	"strconv"
)

func ExportRotaGroupTpl(sheetName string, year, month int, members, lineEnumValues, isMonitorEnumValues []string) (contents []byte, err error) {
	days := getDaysInMonth(year, month)
	// 获取结束列编号
	_, lastCol := getColumnPair(days, "B")

	f := excelize.NewFile()
	_ = f.SetSheetName(f.GetSheetName(0), sheetName)

	//设置A列宽为20
	_ = f.SetColWidth(sheetName, "A", "A", 20)
	//设置 第一行 等高为160
	_ = f.SetRowHeight(sheetName, 1, 160)
	// 设置 第一行 A1单元格内容
	_ = f.SetCellValue(sheetName, "A1", "填写说明：\n1、值班成员需输入人员账号。\n2、模板默认为31日，导入时请根据月实际天数调整。\n3、正式导入前，请删除示例数据。")
	// 设置 第一行 A1单元格样式
	A1Style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",   // 靠左对齐
			Vertical:   "center", // 上下居中
		},
		Font: &excelize.Font{
			Size: 10,
		},
	})
	_ = f.SetCellStyle(sheetName, "A1", "A1", A1Style)
	//合并 第一行 单元格
	_ = f.MergeCell(sheetName, "A1", lastCol+"1")

	//设置第二行 A2单元格内容
	_ = f.SetCellValue(sheetName, "A2", "值班成员")
	//合并 A2:A3 两个单元格式
	_ = f.MergeCell(sheetName, "A2", "A3")
	A2A3Style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 靠左对齐
			Vertical:   "center", // 上下居中
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#A4C93F"}, Pattern: 1},
	})
	_ = f.SetCellStyle(sheetName, "A2", "A3", A2A3Style)

	//设置 第2，3行 样式
	Row3Style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 靠左对齐
			Vertical:   "center", // 上下居中
		},
	})
	_ = f.SetCellStyle(sheetName, "B3", lastCol+"3", Row3Style)

	for i := 1; i <= days; i++ {
		fc, sc := getColumnPair(i, "B")
		//设置 是否值班长 列宽
		_ = f.SetColWidth(sheetName, sc, sc, 11)
		// 设置第二行 单元格样式
		_ = f.SetCellValue(sheetName, fc+"2", strconv.Itoa(i))
		_ = f.MergeCell(sheetName, fc+"2", sc+"2")
		//设置第三行单元格值
		_ = f.SetCellValue(sheetName, fc+"3", "班线")
		_ = f.SetCellValue(sheetName, sc+"3", "是否值班长")

		// 设置 班线 数据验证规则
		_ = f.AddDataValidation(sheetName, &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       true,
			ShowErrorMessage: true,
			ShowInputMessage: true,
			ShowDropDown:     false,
			Sqref:            fmt.Sprintf("%s%d:%s%d", fc, 4, fc, 4+len(members)-1),
			Formula1:         fmt.Sprintf(`"%s"`, joinEnumValues(lineEnumValues)),
			Formula2:         fmt.Sprintf(`"%s"`, joinEnumValues(lineEnumValues)),
			Error:            stringPtr("请从下拉列表中选择一个值"),
			Prompt:           stringPtr("请从下拉列表中选择班次"),
		})

		// 设置 是否值班长 数据验证规则
		_ = f.AddDataValidation(sheetName, &excelize.DataValidation{
			Type:             "list",
			AllowBlank:       true,
			ShowErrorMessage: true,
			ShowInputMessage: true,
			ShowDropDown:     false,
			Sqref:            fmt.Sprintf("%s%d:%s%d", sc, 4, sc, 4+len(members)-1),
			Formula1:         fmt.Sprintf(`"%s"`, joinEnumValues(isMonitorEnumValues)),
			Formula2:         fmt.Sprintf(`"%s"`, joinEnumValues(isMonitorEnumValues)),
			Error:            stringPtr("请从下拉列表中选择一个值"),
			Prompt:           stringPtr("请从下拉列表中选择"),
		})
	}

	// 设置 第二行 B2到结束列表样式
	A2Style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 左右居中
			Vertical:   "center", // 上下居中
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#9ab3d4"}, Pattern: 1},
	})
	_ = f.SetCellStyle(sheetName, "B2", lastCol+"2", A2Style)

	// 写入值组员数据
	A4ToNStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 靠左对齐
			Vertical:   "center", // 上下居中
		},
	})
	for i, u := range members {
		_ = f.SetCellValue(sheetName, "A"+strconv.Itoa(4+i), u)
		_ = f.SetCellStyle(sheetName, "A"+strconv.Itoa(4+i), "A"+strconv.Itoa(4+i), A4ToNStyle)
	}
	//字节流
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(buf)
}

// joinEnumValues 将枚举值列表转换为逗号分隔的字符串
func joinEnumValues(values []string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += ","
		}
		result += value
	}
	return result
}

// stringPtr 返回字符串的指针
func stringPtr(s string) *string {
	return &s
}

// getDaysInMonth 获取某个月的天数
func getDaysInMonth(year, month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

func getColumnPair(day int, startCol string) (string, string) {
	// 将起始列编号转换为列号
	start := columnToIndex(startCol)

	// 计算目标列的偏移量
	offset := (day - 1) * 2

	// 计算第一列和第二列的列号
	firstCol := start + offset
	secondCol := firstCol + 1

	// 将列号转换为列编号
	return indexToColumn(firstCol), indexToColumn(secondCol)
}

// columnToIndex 将列编号转换为列号 (A=0, B=1, ..., Z=25, AA=26, AB=27, ...)
func columnToIndex(col string) int {
	index := 0
	for _, ch := range col {
		index = index*26 + int(ch-'A') + 1
	}
	return index - 1
}

// indexToColumn 将列号转换为列编号 (0=A, 1=B, ..., 25=Z, 26=AA, 27=AB, ...)
func indexToColumn(index int) string {
	col := ""
	for index >= 0 {
		col = string(rune('A'+(index%26))) + col
		index = index/26 - 1
	}
	return col
}

// isLeapYear 判断是否为闰年
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
