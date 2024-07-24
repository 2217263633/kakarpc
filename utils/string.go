package utils

import (
	"strings"
)

func (utils CusUtils) Include(strs string, findStr string, split string) bool {
	var strArr = strings.Split(strs, split)
	var newStrs = strings.ReplaceAll(strs, " ", "")
	if newStrs == "" {
		return false
	}
	for _, str := range strArr {
		if str == findStr {
			return true
		}
	}
	return false
}

func (utils CusUtils) IncludeDelete(strs string, findStr string, split string) string {
	var strArr = strings.Split(strs, split)
	var findStrArr = strings.Split(findStr, split)
	var retuStr = ""
	var index = 0
	for _, str := range strArr {
		isFind := false
		for _, findStr := range findStrArr {
			if str == findStr {
				isFind = true
				break
			}
		}

		if !isFind {
			if str != "" {
				if index == 0 {
					retuStr += str
				} else {
					retuStr += split + str
				}
				index++
			}
		}
	}
	return retuStr
}

func (utils CusUtils) IndexOf(arr []string, val string) int {
	index := -1
	for i, ar := range arr {
		if ar == val {
			index = i
			break
		}
	}
	return index
}

func (utils CusUtils) IndexOfInt(arr []int, val int) int {
	index := -1
	for i, ar := range arr {
		if ar == val {
			index = i
			break
		}
	}
	return index
}

func (utils CusUtils) IncludeAdd(sqlStr string, reviceStrs string) string {
	var vals = ""
	if !utils.Include(sqlStr, reviceStrs, ",") {
		vals = sqlStr + "," + reviceStrs
	} else {
		vals = utils.IncludeDelete(sqlStr, reviceStrs, ",")
		if vals == "" {
			vals = " "
		}
	}
	return vals
}

func ContainsStr(arr []string, val string) int {
	index := -1
	for i, ar := range arr {
		if ar == val {
			index = i
		}
	}
	return index
}

func ContainsNum(arr []interface{}, val int) int {
	index := -1
	for i, ar := range arr {
		if ar.(float64) == float64(val) {
			index = i
			break
		}
	}
	return index
}
