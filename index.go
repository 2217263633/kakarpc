package myrpc

import (
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
