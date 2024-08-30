package rpc

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"time"
)

type ModeType int

const (
	Debug   ModeType = 0
	Release ModeType = 1
)

type ServerStruct struct {
	Chinese_name string                 `yaml:"chinese_name"` // 中文名
	Name         string                 `yaml:"name"`
	Port         int                    `yaml:"port"`
	Swag_port    int                    `yaml:"swag_port"`
	Router       map[string]interface{} `yaml:"router"` //把自己的方法和需要传递的参数写在这里
	Path         string                 `yaml:"path"`
	Mode         ModeType               `yaml:"mode"` // 运行模式，可选值为 "debug" 或 "release"
}

type YamlStruct struct {
	Server ServerStruct `yaml:"server"`
}

type RpcMethod struct {
	Chinese_name string `json:"chinese_name"`
	Method       string `json:"method"`
	Param        any    `json:"param"`
}

type RpcClientType struct {
	Client map[string]*rpc.Client `json:"client"` // rpc客户端
	Heart  time.Time              `json:"heart"`  // 心跳时间
	Addr   string                 `json:"addr"`   // rpc地址
	Name   string                 `json:"name"`   // rpc名称
	Online bool                   `json:"online"` // 是否在线
	Router map[string]interface{} `json:"router"` // 路由表
}

var RpcServer = map[string]*YamlStruct{}

var RpcClient = map[string]*RpcClientType{}

type SqlStruct struct {
	Values       string // 不要写 select
	Tabel_name   string // 不要写 from
	Where        string // 自己写 where 或者 on
	Order        string // 要写 order by
	Page         int    // 自己写页码
	Size         int    // 自己写size大小
	Company_id   int    // 公司id
	Params       string // 自己写参数
	Insert_value string // 插入数据的 insert_value
	Update_value string // 更新数据的 update_value
}

func (_sql SqlStruct) ToString() string {
	if _sql.Values == "" {
		_sql.Values = "*"
	}
	sql := fmt.Sprintf(`select %s from %s `,
		_sql.Values, _sql.Tabel_name)

	if _sql.Where != "" {
		sql += _sql.Where
	}

	if _sql.Order != "" {
		sql += " " + _sql.Order
	}

	return sql
}

func (_sql SqlStruct) ToInsert() string {
	sql := fmt.Sprintf(`insert into %s (%s) values (%s)`, _sql.Tabel_name, _sql.Params, _sql.Insert_value)
	return sql
}

func (_sql SqlStruct) ToUpdate() string {
	// 这里规定更新必须加where
	sql := fmt.Sprintf(`update %s set %s `, _sql.Tabel_name, _sql.Update_value)
	if _sql.Where != "" {
		sql += " " + _sql.Where
	} else {
		sql = "这里规定更新必须加where"
	}

	return sql
}

func (_sql SqlStruct) ToDelete() string {
	sql := fmt.Sprintf(`delete from %s `, _sql.Tabel_name)
	if _sql.Where != "" {
		sql += " " + _sql.Where
	} else {
		sql = "这里规定删除必须加where"
	}

	return sql
}

func (_sql SqlStruct) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"values":       _sql.Values,
		"tabel_name":   _sql.Tabel_name,
		"where":        _sql.Where,
		"order":        _sql.Order,
		"page":         _sql.Page,
		"size":         _sql.Size,
		"company_id":   _sql.Company_id,
		"params":       _sql.Params,
		"insert_value": _sql.Insert_value,
		"update_value": _sql.Update_value,
	}
}

func (_sql SqlStruct) ToJson() []byte {
	j, _ := json.Marshal(_sql)
	return j
}

func (_sql SqlStruct) JsonTo(_json []byte) SqlStruct {
	newMap := SqlStruct{}
	json.Unmarshal(_json, &newMap)
	return newMap
}

func (_sql SqlStruct) MapTo(_map map[string]interface{}) SqlStruct {
	sql := SqlStruct{}
	if _map["values"] != nil {
		sql.Values = _map["values"].(string)
	}
	if _map["tabel_name"] != nil {
		sql.Tabel_name = _map["tabel_name"].(string)
	}
	if _map["where"] != nil {
		sql.Where = _map["where"].(string)
	}
	if _map["order"] != nil {
		sql.Order = _map["order"].(string)
	}
	if _map["page"] != nil {
		sql.Page = _map["page"].(int)
	}
	if _map["size"] != nil {
		sql.Size = _map["size"].(int)
	}
	if _map["company_id"] != nil {
		sql.Company_id = _map["company_id"].(int)
	}
	if _map["params"] != nil {
		sql.Params = _map["params"].(string)
	}
	if _map["insert_value"] != nil {
		sql.Insert_value = _map["insert_value"].(string)
	}
	if _map["update_value"] != nil {
		sql.Update_value = _map["update_value"].(string)
	}
	return sql

}
