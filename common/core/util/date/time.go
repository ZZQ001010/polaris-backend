package date

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/strs"
	"strconv"
	"time"
)

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

/**
在时间基础上增减时间
ParseDuration解析一个时间段字符串。一个时间段字符串是一个序列，每个片段包含可选的正负号、十进制数、
可选的小数部分和单位后缀，如"300ms"、"-1.5h"、"2h45m"。合法的单位有"ns"、"us" /"µs"、"ms"、"s"、"m"、"h"。
*/
func CTime(t time.Time, time_str string) time.Time {
	time_part, err := time.ParseDuration(time_str)
	if err != nil {
		return t
	}
	return t.Add(time_part)
}

//组装时间区间
// secondInterval        区间,
// orgTime    			 原始时kl
// symbol 				 增减或者减少符号,
// intervalUnit    		 增加或者减少单位,
func AssemblyDateTime(secondInterval int, orgTime time.Time, symbol string, intervalUnit string) string {
	//当前时间加上时间区间的时间
	plusTimeStr := fmt.Sprintf("%s"+strconv.Itoa(secondInterval)+"%s", symbol, intervalUnit)
	currentTimeAddInterval := CTime(orgTime, plusTimeStr)
	//转换时间区间的字符串时间 =  当前时间+ 时间区间
	dateTime := repairTimeZero(strconv.Itoa(currentTimeAddInterval.Hour())) + ":" + repairTimeZero(strconv.Itoa(currentTimeAddInterval.Minute()))
	return dateTime
}

func repairTimeZero(str string) string{
	if strs.Len(str) >= 2{
		return str
	}
	return "0" + str
}