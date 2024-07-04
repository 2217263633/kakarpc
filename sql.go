package myrpc

import (
	"encoding/json"
	"errors"
	"fmt"

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
		Method:       "MysqlService.QueryData",
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
