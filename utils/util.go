// Package util tool class function
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Debug TODO: 在终端打印日志,输出具体位置的链接，点击链接可跳转到打印日志的地方
func Debug(logs ...interface{}) {
	// 0表示跳转到下面这行代码的位置，1表示跳转到使用Debug()方法的位置，2表示再上一层位置，以此类推
	_, file, line, _ := runtime.Caller(1)

	var data string
	for _, v := range logs {
		msg, _ := json.Marshal(v)
		data = data + " " + fmt.Sprintf("%v", string(msg))
	}

	fmt.Println(file+":"+strconv.Itoa(line), data)
}

func PageSize(page, pageSize int) (int, int, int) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return page, pageSize, offset
}

func StringsSplit(str string) []string {
	if str == "" {
		return nil
	}
	return strings.Split(str, ",")
}

// StrToTime 字符串转时间 2023-11-17 20:05:44
func StrToTime(str string) time.Time {
	t, _ := time.ParseInLocation(time.DateTime, str, time.Local)
	return t
}

// StrToTimeRFC 字符串转时间 2023-11-17T20:05:44Z
func StrToTimeRFC(str string) time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, str, time.Local)
	return t
}

// ReadFile 读取文件
func ReadFile(url string) ([]byte, error) {
	file, err := ioutil.ReadFile(url)
	return file, err
}

// GetStartAndEndOfMonth 获取某年某月的开始时间和结束时间
func GetStartAndEndOfMonth(year, month int) (time.Time, time.Time) {
	// 获取指定年月的第一天
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// 获取下个月的第一天，然后减去一秒得到当前月的最后一天
	nextMonth := startOfMonth.AddDate(0, 1, 0)
	endOfMonth := nextMonth.Add(-time.Second)

	return startOfMonth, endOfMonth
}

// GetDatesBetween 获取两个时间的日期集
func GetDatesBetween(start string, end string) ([]string, error) {
	startDate, err := time.ParseInLocation(time.DateTime, start, time.Local)
	if err != nil {
		return nil, err
	}

	endDate, err := time.ParseInLocation(time.DateTime, end, time.Local)
	if err != nil {
		return nil, err
	}

	var dates []string

	for !startDate.After(endDate) {
		day := startDate.Format(time.DateTime)
		dates = append(dates, day)
		startDate = startDate.AddDate(0, 0, 1) // 增加一天
	}

	return dates, nil
}

// FindIntersection 两个切片的交集
func FindIntersection(slice1, slice2 []string) []string {
	set := make(map[string]bool)
	var intersection []string

	for _, num := range slice1 {
		set[num] = true
	}

	for _, num := range slice2 {
		if set[num] {
			intersection = append(intersection, num)
		}
	}

	return intersection
}

// GetRandString 获取长度为len*2的16进制的随机字符串
func GetRandString(prefix string) string {
	bt := make([]byte, 3)
	_, _ = rand.Read(bt)
	return prefix + time.Now().Local().Format("20060102150405") + hex.EncodeToString(bt)
}

func GetDateOfDay(months string, day int) (string, error) {
	// 检查月份和天数是否有效
	pt, err := time.Parse("2006-01", months)
	if err != nil {
		return "", err
	}
	year, month := pt.Year(), pt.Month()
	if month < 1 || month > 12 {
		return "", fmt.Errorf("月份无效: %d", month)
	}
	if day < 1 || day > 31 {
		return "", fmt.Errorf("天数无效: %d", day)
	}

	// 创建时间对象
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// 检查日期是否有效（例如 2 月 30 日是无效日期）
	if t.Month() != month || t.Day() != day {
		return "", fmt.Errorf("日期无效: %d-%02d-%02d", year, month, day)
	}
	// 返回日期字符串
	return t.Format("2006-01-02"), nil
}

// SliceOperator SliceOperate 切片操作
// 获取两切片的 合集 交集 差值(A—B) 差值(B—A)
func SliceOperator[T any](A, B []T, equalFunc func(a, b T) bool) (union, intersection, differenceAB, differenceBA []T) {
	// 使用切片来存储元素
	var setA []T
	var setB []T

	// 将数组 A 的元素存入 setA（去重）
	for _, item := range A {
		if !contains(setA, item, equalFunc) {
			setA = append(setA, item)
		}
	}

	// 将数组 B 的元素存入 setB（去重）
	for _, item := range B {
		if !contains(setB, item, equalFunc) {
			setB = append(setB, item)
		}
	}

	// 计算交集
	for _, itemA := range setA {
		for _, itemB := range setB {
			if equalFunc(itemA, itemB) {
				intersection = append(intersection, itemA)
				break
			}
		}
	}

	// 计算并集
	union = append(union, setA...)
	for _, itemB := range setB {
		if !contains(union, itemB, equalFunc) {
			union = append(union, itemB)
		}
	}

	// 计算差集 (A - B)
	for _, itemA := range setA {
		if !contains(setB, itemA, equalFunc) {
			differenceAB = append(differenceAB, itemA)
		}
	}

	// 计算差集 (B - A)
	for _, itemB := range setB {
		if !contains(setA, itemB, equalFunc) {
			differenceBA = append(differenceBA, itemB)
		}
	}

	return intersection, union, differenceAB, differenceBA
}

// SliceRemoveDuplicates 泛型去重函数，支持自定义比较函数
func SliceRemoveDuplicates[T any](slice []T, equalFunc func(a, b T) bool) []T {
	seen := make(map[int]struct{}) // 用于存储已出现的元素索引
	result := []T{}                // 去重后的结果
	for i, item := range slice {
		isDuplicate := false
		for j := range seen {
			if equalFunc(slice[j], item) { // 使用自定义比较函数
				isDuplicate = true
				break
			}
		}
		if !isDuplicate { // 如果元素未出现过
			seen[i] = struct{}{}          // 标记为已出现
			result = append(result, item) // 添加到结果中
		}
	}
	return result
}

// DurationDuplicateCheck 检查时间段是否有重叠
func DurationDuplicateCheck(timeRanges [][2]string) bool {
	var parsedRanges [][2]time.Time
	baseDate := time.Now().Format("2006-01-02") // 设定为当天日期
	layout := "2006-01-02 15:04"
	for _, tr := range timeRanges {
		start, _ := time.Parse(layout, baseDate+" "+tr[0])
		end, _ := time.Parse(layout, baseDate+" "+tr[1])
		// 如果开始时间大于结束时间，说明跨天，结束时间加一天
		if start.After(end) {
			end = end.Add(24 * time.Hour)
		}
		parsedRanges = append(parsedRanges, [2]time.Time{start, end})
	}
	// 按照开始时间排序
	sort.Slice(parsedRanges, func(i, j int) bool {
		return parsedRanges[i][0].Before(parsedRanges[j][0])
	})
	// 检查时间段是否重叠
	for i := 1; i < len(parsedRanges); i++ {
		prevEnd := parsedRanges[i-1][1]
		currStart := parsedRanges[i][0]
		// 当前时间段的开始时间如果小于前一个时间段的结束时间，则存在重叠
		if currStart.Before(prevEnd) {
			return true
		}
	}
	return false
}

// contains 检查切片中是否包含某个元素
func contains[T any](slice []T, item T, equalFunc func(a, b T) bool) bool {
	for _, s := range slice {
		if equalFunc(s, item) {
			return true
		}
	}
	return false
}
