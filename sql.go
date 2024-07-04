package myrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

func QuerySql(Rpc *RPC, sql string) ([]map[string]interface{}, error) {
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

func JudgeTable(Rpc *RPC, table string) ([]map[string]interface{}, error) {
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

func PageSql(Rpc *RPC, sql SqlStruct) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "数据库调用",
		Method:       "MysqlService.PageSql",
		Param:        sql}, &data)
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

func CallToken(Rpc *RPC, token string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", RpcMethod{
		Chinese_name: "用户服务",
		Method:       "UserService.ParseToken",
		Param:        token}, &data)
	logger.Info(data)
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

func CallAny(Rpc *RPC, method string, param string, chinese_name string) (interface{}, error) {
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
		(*res)["err"] = err.Error()
	} else {
		(*res)["state"] = 200
	}
}

func StructToSql(fileStruct any, obj map[string]interface{}) ([]string, []interface{}) {
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
