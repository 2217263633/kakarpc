package cuswebsocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"
)

type RpcMethod struct {
	Chinese_name string `json:"chinese_name"`
	Method       string `json:"method"`
	Param        any    `json:"param"`
}

type RPC struct {
	Client      *rpc.Client  `json:"client"`
	Count       int          `json:"count"`       // 重联计数器
	R           *gin.Engine  `json:"r"`           // gin框架
	Conn        any          `json:"conn"`        // 连接注册中心
	Swag_port   int          `json:"swag_port"`   // swagger端口
	Conect_port int          `json:"conect_port"` // 远程方调用接口时候，需要从自己这里进行返回
	Srv         *http.Server `json:"srv"`         // 关闭gin的端口
}

// 特定的调用
func CallOther(Rpc *RPC, method RpcMethod) (interface{}, error) {
	var data map[string]interface{}
	Rpc.Client.Call("RPC.Call", method, &data)

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
