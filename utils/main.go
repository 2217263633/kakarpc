package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type CusUtils struct{}

func UtilsInit() *CusUtils {
	return &CusUtils{}
}

func (utils CusUtils) ListToString(arrs []string) string {
	tempStr := ""
	for _, arr := range arrs {
		tempStr += arr + ","
	}
	if len(tempStr) > 0 {
		return tempStr[:len(tempStr)-1]
	}
	return tempStr
}

func (utils CusUtils) ListToAny(arrs []any) string {
	tempStr := ""
	for _, arr := range arrs {
		var x = fmt.Sprintf("%T", arr)
		// logger.Info(x)
		if x == "float64" {
			floatStr := fmt.Sprintf("%f", arr)
			tempStr += floatStr + ","
		} else if x == "string" {
			tempStr += "'" + arr.(string) + "'" + ","
		} else if x == "int64" {
			floatStr := fmt.Sprintf("%d", arr)
			tempStr += floatStr + ","
		} else if x == "time.Time" {
			var value = arr.(time.Time).Format("2006-01-02 15:04:05")
			tempStr += "'" + value + "',"
		} else if x == "int" {
			floatStr := fmt.Sprintf("%d", arr)
			tempStr += floatStr + ","
		}
	}
	if len(tempStr) > 0 {
		return tempStr[:len(tempStr)-1]
	}
	return tempStr
}

func (utils CusUtils) ListToListStr(arrs []any) []string {
	var arrStr []string = make([]string, 0)

	for _, arr := range arrs {
		var x = fmt.Sprintf("%T", arr)
		if x == "float64" {
			floatStr := fmt.Sprintf("%f", arr)
			arrStr = append(arrStr, floatStr)
		} else if x == "string" {
			arrStr = append(arrStr, "'"+arr.(string)+"'")
		} else if x == "int64" {
			floatStr := fmt.Sprintf("%d", arr)
			arrStr = append(arrStr, floatStr)
		} else if x == "time.Time" {
			floatStr := arr.(time.Time).Format("2006-01-02 15:04:05")
			arrStr = append(arrStr, "'"+floatStr+"'")
		} else if x == "int" {
			floatStr := fmt.Sprintf("%d", arr)
			arrStr = append(arrStr, floatStr)
		}
	}
	return arrStr
}

// UPDATE
func (utils CusUtils) ListLog(fields []string, arrs []any) string {
	var tring = ""
	if len(fields) != len(arrs) {
		return tring
	}
	arrStr := utils.ListToListStr(arrs)
	// logger.Info(fields, arrStr)
	for i := 0; i < len(fields); i++ {
		// logger.Info(fields[i], arrStr[i])
		tring += fields[i] + "=" + arrStr[i] + ","
	}

	if len(tring) > 0 {
		tring = tring[:len(tring)-1]
	}
	return tring
}

func (utils CusUtils) MergeMap(x map[string]interface{}, y map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{})
	for i, v := range x {
		if _, ok := y[i]; !ok {
			y[i] = v
		}
	}
	n = y
	return n
}

// match 一般指 id   goal指的是parent_id
func (utils CusUtils) RecursionList(arr []map[string]interface{}, match string, goal string) []map[string]interface{} {
	// var _tempArrs []map[string]interface{} = make([]map[string]interface{}, 0)
	for _, v := range arr {
		for _, j := range arr {
			if j[goal] == v[match] {
				j["isdelete"] = true
				if v["children"] == nil {
					var children []map[string]interface{} = make([]map[string]interface{}, 0)
					children = append(children, j)
					v["children"] = children
				} else {
					v["children"] = append(v["children"].([]map[string]interface{}), j)
				}
				// break
			}
		}
	}

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i]["isdelete"] == true {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	return arr

}

// 作为补充 会以第一个name为开始寻找
// 一种递归的方法，不是特别好
func ReciursionName(arr []map[string]interface{}, id string, parent_id string, name string) map[string]interface{} {
	_map := make(map[string]interface{})
	for _, v := range arr {
		if v["name"] == name {
			_map = v
			break
		}
	}

	if len(_map) == 0 {
		return _map
	}

	for _, v := range arr {
		if v[parent_id] == _map[id] {
			_sd := make([]map[string]interface{}, 0)
			if _map["children"] != nil {
				_sd = _map["children"].([]map[string]interface{})
			}
			_sd = append(_sd, ReciursionName(arr, id, parent_id, v["name"].(string)))
			_map["children"] = _sd
		}
	}

	return _map
}

// 作为补充 会以第一个name为开始寻找
func RecursionListName(arr []map[string]interface{}, match string, goal string, name string) []map[string]interface{} {
	// var _tempArrs []map[string]interface{} = make([]map[string]interface{}, 0)
	_app := make([]map[string]interface{}, 0)
	for _, v := range arr {
		if v["name"] != name {
			continue
		}
		// logger.Info(v["name"])
		_app = append(_app, v)
		// v["children"] = make([]map[string]interface{}, 0)
		for _, j := range arr {
			if j[goal] == v[match] {

				if v["children"] == nil {
					var children []map[string]interface{} = make([]map[string]interface{}, 0)
					children = append(children, j)
					v["children"] = children
				} else {
					v["children"] = append(v["children"].([]map[string]interface{}), j)
				}
				// // break
			} else {
				j["isdelete"] = true
			}
		}
	}

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i]["isdelete"] == true {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}

	if len(_app) > 0 {
		_app[0]["children"] = arr
	}

	return _app

}

func (utils CusUtils) Remove(arr []map[string]interface{}, index int) []map[string]interface{} {
	return append(arr[:index], arr[index+1:]...)
}

func (utils CusUtils) FindParent(mapValue map[string]interface{}, type_ string, parent_id string, id string,
	re map[string]interface{}) {
	// 判断自己是否有父节点
	if mapValue[type_+parent_id] != nil {
		mapValue[type_+id].(map[string]interface{})["isdelete"] = true
		_count := mapValue[type_+parent_id].(map[string]interface{})["count"].(int)
		_count -= 1
		mapValue[type_+parent_id].(map[string]interface{})["count"] = _count
		mapValue[type_+parent_id].(map[string]interface{})["children"] =
			append(mapValue[type_+parent_id].(map[string]interface{})["children"].([]interface{}), re)
	}
}

// parameter 存在的时候 不创建未绑定单元
func (utils CusUtils) FindLocationId(mapValue map[string]interface{}, type_ string, location_id int, id string,
	re map[string]interface{}, parameter string) {
	if type_ == "device_file" {

		if mapValue["location"+strconv.Itoa(location_id)] != nil {
			mapValue["location"+strconv.Itoa(location_id)].(map[string]interface{})["children"] =
				append(mapValue["location"+strconv.Itoa(location_id)].(map[string]interface{})["children"].([]interface{}), re)
			mapValue["device_file"+id].(map[string]interface{})["isdelete"] = true
		} else {
			if parameter == "" {
				mapValue["location"+strconv.Itoa(location_id)] = map[string]interface{}{
					"id":          location_id,
					"name":        "未绑定地点" + strconv.Itoa(location_id),
					"children":    []interface{}{re},
					"type":        "location",
					"business_id": "",
					"isdelete":    false,
				}
				mapValue["device_file"+id].(map[string]interface{})["isdelete"] = true
			}
		}

	} else if type_ == "device" {
		if mapValue["device_file"+strconv.Itoa(location_id)] != nil {
			mapValue["device_file"+strconv.Itoa(location_id)].(map[string]interface{})["children"] =
				append(mapValue["device_file"+strconv.Itoa(location_id)].(map[string]interface{})["children"].([]interface{}), re)
		} else {
			if parameter == "" {
				mapValue["device_file"+strconv.Itoa(location_id)] = map[string]interface{}{
					"id":          location_id,
					"name":        "未绑定单元" + strconv.Itoa(location_id) + re["name"].(string),
					"children":    []interface{}{re},
					"type":        "device_file",
					"isdelete":    false,
					"business_id": "",
				}
			}
		}
		mapValue["device"+id].(map[string]interface{})["isdelete"] = true
	}
}

func (utils CusUtils) ToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func (utils CusUtils) ToDegrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func (utils CusUtils) Calc_center(lngs []float64, lats []float64) (float64, float64) {
	_lng := 0.0
	_lat := 0.0
	count := len(lngs)

	for i := 0; i < count; i++ {
		_lng += utils.ToRadians(lngs[i])
		_lat += utils.ToRadians(lats[i])
	}

	lng_avg := utils.ToDegrees(_lng / float64(count))
	lat_avg := utils.ToDegrees(_lat / float64(count))
	return lng_avg, lat_avg
}

func (utils CusUtils) Calcu_distanc(lat1, lng1, lat2, lng2 float64) int {
	lat1Rad := utils.ToRadians(lat1)
	lon1Rad := utils.ToRadians(lng1)
	lat2Rad := utils.ToRadians(lat2)
	lon2Rad := utils.ToRadians(lng2)
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := 6371 * c * 1000.0
	return int(distance)
}
