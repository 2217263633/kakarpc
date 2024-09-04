package sql

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/2217263633/kakarpc/gal"
	"github.com/2217263633/kakarpc/types"
	"github.com/wonderivan/logger"
)

type SqlStruct struct{}

func Init() *SqlStruct {
	return &SqlStruct{}
}

func (u *SqlStruct) QuerySql(sql string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) JudgeTable(table string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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
func (u *SqlStruct) PageSql(sql types.SqlStruct) ([]map[string]interface{}, int, int, error) {
	map_sql := sql.ToMap()
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) CallToken(token string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) CallAny(method string, param string, chinese_name string) (interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) ErrDeal(err error, res *map[string]interface{}) {
	if err != nil {
		(*res)["state"] = 401
		// 因为rpc的原因 必须传string
		(*res)["err"] = err.Error()
	} else {
		(*res)["state"] = 200
	}
}

func (u *SqlStruct) StructToSql(fileStruct any) ([]string, []interface{}) {
	t := reflect.TypeOf(fileStruct)
	value := reflect.ValueOf(fileStruct)

	fieleNum := t.NumField()
	var typeStr []string = make([]string, 0)
	var valueStr []interface{} = make([]interface{}, 0)
	for i := 0; i < fieleNum; i++ {
		if t.Field(i).Tag.Get("show") != "" {
			continue
		}
		// logger.Info(t.Field(i).Type.String(), t.Field(i).Name)
		// logger.Info(t.Field(i).Type.String(), t.Field(i).Name)
		// logger.Info(t.Field(i).Type.String(), t.Field(i).Type.Name(), value.Field(i))
		if t.Field(i).Type.Name() == "int" {
			if value.Field(i).Int() != 0 {
				valueStr = append(valueStr, value.Field(i).Int())
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
			}
		} else if t.Field(i).Type.Name() == "float64" {
			if value.Field(i).Float() > 0 {
				valueStr = append(valueStr, value.Field(i).Float())
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
			}
		} else if t.Field(i).Type.Name() == "float32" {
			if value.Field(i).Float() > 0 {
				valueStr = append(valueStr, value.Field(i).Float())
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
			}
		} else if t.Field(i).Type.Name() == "string" {
			if len(value.Field(i).String()) > 0 {
				valueStr = append(valueStr, value.Field(i).String())
				typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
			}
		} else if t.Field(i).Type.Name() == "Time" {

			if value.Field(i).Interface().(time.Time).Year() < 2000 {
				continue
			}
			valueStr = append(valueStr, value.Field(i).Interface().(time.Time).Local().Format("2006-01-02 15:04:05"))
			// logger.Info(value.Field(i), value.Field(i).Interface().(time.Time).Local().Format("2006-01-02 15:04:05"))
			typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
		} else if t.Field(i).Type.String() == "[]int" {
			var intArr = value.Field(i).Interface().([]int)
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
			typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
		} else if t.Field(i).Type.String() == "[]string" {
			var strArr = value.Field(i).Interface().([]string)
			intStr := ""
			for index, intVal := range strArr {
				if index == len(strArr)-1 {
					intStr += intVal
				} else {
					intStr += intVal + ","
				}
			}
			valueStr = append(valueStr, intStr)
			typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))
		} else if t.Field(i).Type.String() == "bool" {

			var strArr = value.Field(i).Interface().(bool)
			intStr := 0
			if strArr {
				intStr = 1

			}
			valueStr = append(valueStr, intStr)
			typeStr = append(typeStr, strings.ToLower(t.Field(i).Name))

		}
	}

	return typeStr, valueStr
}

func (u *SqlStruct) StructToSql2(fileStruct any, obj map[string]interface{}) ([]string, []interface{}) {
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

func (u *SqlStruct) InsertTable(sql string) error {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) InsertTableId(sql string) (int, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) QueryIdlimit1(tableName string) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) CreateTable(sql string) error {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) UpdateTable(sql string) error {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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

func (u *SqlStruct) DeleteTable(sql string) error {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", types.RpcMethod{
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
func (u *SqlStruct) CallOther(method types.RpcMethod) (interface{}, error) {
	var data map[string]interface{}
	gal.MyRpc.Client.Call("RPC.Call", method, &data)

	if data == nil || data["state"] == nil {
		return []map[string]interface{}{}, errors.New("数据异常")
	}

	var list_sql interface{}
	if data["data"] != nil {
		json.Unmarshal(data["data"].([]byte), &list_sql)
	}
	if data["err"] == nil {
		return list_sql, nil
	}
	if fmt.Sprintf("%T", data["err"]) == "string" {
		return list_sql, errors.New(data["err"].(string))
	}
	return list_sql, data["err"].(error)
}
