package myrpc

import (
	"github.com/2217263633/kakarpc/cuswebsocket"
	"github.com/2217263633/kakarpc/sql"
	"github.com/2217263633/kakarpc/types"
	"github.com/2217263633/kakarpc/utils"
)

var MyRpc = &RPC{}
var Utils *utils.CusUtils = utils.UtilsInit()
var Websocket *cuswebsocket.Cuswebsocket = cuswebsocket.Init()
var Mysql = sql.Init()
var RpcTypes = types.Init()
