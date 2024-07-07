package myrpc

import (
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
}

var RpcServer = map[string]*YamlStruct{}

var RpcClient = map[string]*RpcClientType{}

type SqlStruct struct {
	Values     string // 不要写 select
	Tabel_name string // 不要写 from
	Where      string // 自己写 where 或者 on
	Order      string // 要写 order by
	Page       int    // 自己写页码
	Size       int    // 自己写size大小
}

func (_sql SqlStruct) ToString() string {
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

func (_sql SqlStruct) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"values":     _sql.Values,
		"tabel_name": _sql.Tabel_name,
		"where":      _sql.Where,
		"order":      _sql.Order,
		"page":       _sql.Page,
		"size":       _sql.Size,
	}
}

func (_sql SqlStruct) MapTo(_map map[string]interface{}) SqlStruct {
	return SqlStruct{
		Values:     _map["values"].(string),
		Tabel_name: _map["tabel_name"].(string),
		Where:      _map["where"].(string),
		Order:      _map["order"].(string),
		Page:       _map["page"].(int),
		Size:       _map["size"].(int),
	}

}
