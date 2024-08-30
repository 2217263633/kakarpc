package myrpc

import (
	"github.com/2217263633/kakarpc/cuswebsocket"
	"github.com/2217263633/kakarpc/rpc"
	"github.com/2217263633/kakarpc/utils"
)

var MyRpc = &rpc.RPC{}
var Utils *utils.CusUtils = utils.UtilsInit()
var Websocket *cuswebsocket.Cuswebsocket = cuswebsocket.Init()
