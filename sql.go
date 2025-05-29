package myrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/2217263633/kakarpc/tool"
)

func QuerySql(sql string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.QueryData",
		Param:        sql}, &data)
	if data == nil || data["data"] == nil {
		return []map[string]interface{}{}, errors.New("数据库服务已离线，请联系管理员")
	}

	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if data["err"] == nil {
		return list_sql, nil
	}
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}

	return list_sql, data["err"].(error)
}

func JudgeTable(table string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.JudgeTable",
		Param:        table}, &data)
	if data == nil || data["data"] == nil {
		return []map[string]interface{}{}, errors.New("数据库服务已离线，请联系管理员")
	}
	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}

	return list_sql, data["err"].(error)
}

// resu ,total,size,error
func PageSql(sql SqlStruct) ([]map[string]interface{}, int, int, error) {
	map_sql := sql.ToMap()
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.PageSql",
		Param:        map_sql}, &data)
	if data == nil || data["data"] == nil {
		return []map[string]interface{}{}, 0, 0, errors.New("数据库服务已离线，请联系管理员")
	}
	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, 0, sql.Size, errors.New(data["err"].(string))
	}

	if fmt.Sprintf("%T", data["err"]) == "error" {
		return list_sql, 0, sql.Size, data["err"].(error)
	}

	return list_sql, data["total"].(int), data["size"].(int), nil
}

func CallToken(token string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "用户服务",
		Method:       "UserService.ParseToken",
		Param:        token}, &data)

	if data == nil || data["data"] == nil {
		return []map[string]interface{}{}, errors.New("数据异常")
	}

	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if data["err"] == nil {
		return list_sql, nil
	}
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}
	return list_sql, data["err"].(error)
}

func CallAny(method string, param any, chinese_name string) (interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: chinese_name,
		Method:       method,
		Param:        param}, &data)

	if data == nil || data["data"] == nil {
		return []map[string]interface{}{}, errors.New("数据异常")
	}

	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if data["err"] == nil {
		return list_sql, nil
	}
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}
	return list_sql, data["err"].(error)
}

func ErrDeal(err error, res *map[string]interface{}) {
	if err != nil {
		(*res)["state"] = 401
		// 因为rpc的原因 必须传string
		(*res)["err"] = err.Error()
	} else {
		(*res)["state"] = 200
	}
}

func StructToSql(fileStruct any) ([]string, []interface{}) {
	t := reflect.TypeOf(fileStruct)

	value := reflect.ValueOf(fileStruct)
	_map := tool.StructToMap(fileStruct)

	var typeStr []string = make([]string, 0)
	var valueStr []interface{} = make([]interface{}, 0)
	// 例外  因为现在map被排序了，导致有些字段不在第一位，所以这里需要判断一下
	// 我们把 Company_id 放到第一位，其他字段放到最后一位 2024年12月31日 kasia
	if _map["Company_id"] != nil {
		key, isbool := t.FieldByName("Company_id")
		if isbool {
			if value.FieldByName(key.Name).Int() != 0 {
				valueStr = append(valueStr, value.FieldByName(key.Name).Int())
				typeStr = append(typeStr, strings.ToLower(key.Name))
			}
			delete(_map, "Company_id")
		}
	}
	for k := range _map {
		key, isbool := t.FieldByName(k)
		if isbool {
			if key.Tag.Get("show") != "" {
				continue
			} else if key.Type.Name() == "int" {
				if value.FieldByName(key.Name).Int() != 0 {
					valueStr = append(valueStr, value.FieldByName(key.Name).Int())
					typeStr = append(typeStr, strings.ToLower(key.Name))
				}
			} else if key.Type.Name() == "float64" {
				if value.FieldByName(key.Name).Float() != 0.0 {
					valueStr = append(valueStr, value.FieldByName(key.Name).Float())
					typeStr = append(typeStr, strings.ToLower(key.Name))
				}
			} else if key.Type.Name() == "float32" {
				if value.FieldByName(key.Name).Float() != 0.0 {
					valueStr = append(valueStr, value.FieldByName(key.Name).Float())
					typeStr = append(typeStr, strings.ToLower(key.Name))
				}
			} else if key.Type.Name() == "string" {
				if len(value.FieldByName(key.Name).String()) > 0 {
					valueStr = append(valueStr, value.FieldByName(key.Name).String())
					typeStr = append(typeStr, strings.ToLower(key.Name))
				}
			} else if key.Type.Name() == "Time" {

				if value.FieldByName(key.Name).Interface().(time.Time).Year() < 2000 {
					continue
				}
				valueStr = append(valueStr, value.FieldByName(key.Name).Interface().(time.Time).Local().Format("2006-01-02 15:04:05"))
				typeStr = append(typeStr, strings.ToLower(key.Name))
			} else if key.Type.String() == "[]int" {
				var intArr = value.FieldByName(key.Name).Interface().([]int)
				if len(intArr) == 0 {
					continue
				}
				intStr := ""
				for index, intVal := range intArr {
					if index == len(intArr)-1 {
						intStr += strconv.Itoa(intVal)
					} else {
						intStr += strconv.Itoa(intVal) + ","
					}
				}
				valueStr = append(valueStr, intStr)
				typeStr = append(typeStr, strings.ToLower(key.Name))
			} else if key.Type.String() == "[]string" {
				var strArr = value.FieldByName(key.Name).Interface().([]string)
				intStr := ""
				for index, intVal := range strArr {
					if index == len(strArr)-1 {
						intStr += intVal
					} else {
						intStr += intVal + ","
					}
				}
				valueStr = append(valueStr, intStr)
				typeStr = append(typeStr, strings.ToLower(key.Name))
			} else if key.Type.String() == "bool" {
				var strArr = value.FieldByName(key.Name).Interface().(bool)
				intStr := 0
				if strArr {
					intStr = 1
				}
				valueStr = append(valueStr, intStr)
				typeStr = append(typeStr, strings.ToLower(key.Name))
			} else {
				_type := fmt.Sprintf("%T", value.FieldByName(key.Name).Interface())
				if _type == "types.WarningType" {
					if value.FieldByName(key.Name).Int() != 0 {
						valueStr = append(valueStr, value.FieldByName(key.Name).Int())
						typeStr = append(typeStr, strings.ToLower(key.Name))
					}
				}

				//  else if strings.Contains(t.Field(i).Type.String(), "From") {
				// 	// logger.Info(t.Field(i).Type.String(), t.Field(i).Name, value.Field(i), "---",
				// 	// 	t.Field(i).Type.Kind(),
				// 	// 	_map,
				// 	// )
				// 	// fields, values := StructToSql(_map[t.Field(i).Name])
				// 	// logger.Info(fields, values)
				// }
				// logger.Info(key.Type.Name(), value.FieldByName(key.Name))
			}
		}
	}
	return typeStr, valueStr
}

func StructToSql2(fileStruct any, obj map[string]interface{}) ([]string, []interface{}) {
	var typeStr []string = make([]string, 0)
	var valueStr []interface{} = make([]interface{}, 0)
	t := reflect.TypeOf(fileStruct)
	fieleNum := t.NumField()
	value := reflect.ValueOf(fileStruct)
	for i := 0; i < fieleNum; i++ {
		var _names = strings.ToLower(t.Field(i).Name)
		if t.Field(i).Tag.Get("show") != "" {
			continue
		}
		if obj[_names] != nil {
			if t.Field(i).Type.Name() == "Time" {
				if value.Field(i).Interface().(time.Time).Year() < 2000 {
				} else {
					typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
					valueStr = append(valueStr, value.Field(i).Interface().(time.Time).Local().Format("2006-01-02 15:04:05"))
				}
				continue
			}
			if t.Field(i).Type.String() == "[]int" {
				var intArr = obj[_names].([]interface{})
				if len(intArr) == 0 {
					continue
				}
				intStr := ""
				for index, intVal := range intArr {
					if index == len(intArr)-1 {
						intStr += strconv.Itoa(int(intVal.(float64)))
					} else {
						intStr += strconv.Itoa(int(intVal.(float64))) + ","
					}
				}
				valueStr = append(valueStr, intStr)
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
				continue
			}
			// logger.Info(value.Field(i).Interface(), fmt.Sprintf("%T", value.Field(i).Interface()), ":",
			// 	strings.ToLower(t.Field(i).Name), fmt.Sprintf(`%T`, strings.ToLower(t.Field(i).Name)))

			typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
			// valueStr = append(valueStr, obj[strings.ToLower(t.Field(i).Name)])
			valueStr = append(valueStr, value.Field(i).Interface())
		} else if _names == "companyid" {
			if obj["companyId"] != nil {
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
				valueStr = append(valueStr, obj["companyId"])
			}
		}
		// logger.Info(obj[_names], _names)
	}

	return typeStr, valueStr
}

func InsertTable(sql string) error {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.InsertData",
		Param:        sql}, &data)
	if data == nil || data["state"] == nil {
		return errors.New("数据库服务已离线，请联系管理员")
	}
	// var list_sql []map[string]interface{}
	// json.Unmarshal(data["data"].([]byte), &list_sql)
	if data["err"] != nil {
		return errors.New(data["err"].(string))
	}
	return nil
}

func InsertTableId(sql string) (int, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.InsertTableId",
		Param:        sql}, &data)
	if data == nil || data["state"] == nil {
		return -1, errors.New("数据库服务已离线，请联系管理员")
	}

	if data["err"] != nil {
		return 0, errors.New(data["err"].(string))
	}
	var list_sql interface{}
	err := json.Unmarshal(data["data"].([]byte), &list_sql)
	if err != nil {
		return 0, err
	}
	return int(list_sql.(float64)), nil
}

func QueryIdlimit1(tableName string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.QueryIdlimit1",
		Param:        tableName}, &data)
	if data == nil || data["state"] == nil {
		return []map[string]interface{}{}, errors.New("数据库服务已离线，请联系管理员")
	}
	var list_sql []map[string]interface{}
	json.Unmarshal(data["data"].([]byte), &list_sql)
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}

	if fmt.Sprintf("%T", data["err"]) == "error" {
		return list_sql, data["err"].(error)
	}

	return list_sql, nil
}

func CreateTable(sql string) error {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.CreateTable",
		Param:        sql}, &data)

	if data == nil {
		return errors.New("数据库服务已离线，请联系管理员")
	}

	if data["err"] != nil {
		return errors.New(data["err"].(string))
	}
	return nil
}

func UpdateTable(sql string) error {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.UpdateTable",
		Param:        sql}, &data)

	if data == nil {
		return errors.New("数据库服务已离线，请联系管理员")
	}

	if data["err"] != nil {
		return errors.New(data["err"].(string))
	}
	return nil
}

func DeleteTable(sql string) error {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.DeleteData",
		Param:        sql}, &data)
	if data == nil {
		return errors.New("数据库服务已离线，请联系管理员")
	}

	if data["err"] != nil {
		return errors.New(data["err"].(string))
	}
	return nil
}

// 特定的调用 Rpc 不需要在传，本身意义不大
func CallOther(method RpcMethod) (interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", method, &data)

	if data == nil || data["state"] == nil {
		return []map[string]interface{}{}, errors.New("数据异常")
	}

	var list_sql interface{}
	if data["data"] != nil {
		_, ok := data["data"].([]byte)
		if ok {
			json.Unmarshal(data["data"].([]byte), &list_sql)
			data["data"] = list_sql
		}
	}

	if data["err"] == nil {
		return data, nil
	}
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return data, errors.New(data["err"].(string))
	}
	return data, data["err"].(error)
}
