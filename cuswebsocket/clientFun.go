package cuswebsocket

import (
	"encoding/json"
	"fmt"

	"os"
	"strconv"

	"github.com/2217263633/kakarpc/cusrequest"
	"github.com/2217263633/kakarpc/rpc"
	"github.com/wonderivan/logger"
)

type Cuswebsocket struct{}

func Init() *Cuswebsocket {
	return &Cuswebsocket{}
}

// func FindUser(UserId int) *Client {
// 	var conn *Client
// 	for _, _conn := range Manager.Clients {
// 		if _conn.UserId == UserId {
// 			conn = _conn
// 			break
// 		}
// 	}
// 	return conn
// }

// 主要是发送通知  user_id 是发送人
func (c *Cuswebsocket) GetClient(_msg WsMessage, token string, user_id int) error {
	senUrl := "https://chat.kasiasafe.top:8091/api/v1/ws/sendMsg"
	if os.Args[len(os.Args)-1] == "test" {
		senUrl = "http://testqiye.kasiasafe.top:8091/api/v1/ws/sendMsg"
	} else {
		senUrl = "http://127.0.0.1:8091/api/v1/ws/sendMsg"
	}
	if _msg.Business == 0 {
		_msg.Business = 1
	}
	if _msg.Type == 0 {
		_msg.Type = 4
	}
	// logger.Error(senUrl, _msg.Data)
	_, err := cusrequest.Request(senUrl, cusrequest.Post, map[string]interface{}{
		"business": _msg.Business,
		"data":     _msg.Data,
		"userId":   _msg.UserId,
		"type":     _msg.Type,
		"callUrl":  _msg.CallUrl,
	}, token)
	logger.Error(err, _msg.Data)
	if err != nil {
		logger.Info(_msg.UserId, "不在线", err)
		resp, _ := json.Marshal(&WsMessage{
			Data:     _msg.Data,
			Business: _msg.Business,
			UserId:   _msg.UserId,
			CallUrl:  _msg.CallUrl,
		})
		_revice_user_id, _ := strconv.Atoi(_msg.UserId)
		c.NotFind(user_id, _revice_user_id, string(resp), "")
		logger.Info("已把离线消息存入数据库，等待他上线查看")
	}

	return err
}

var Rpc = &rpc.RPC{}

func (c *Cuswebsocket) NotFind(userId int, send_user_id int, data string, parameter string) error {

	sql := rpc.SqlStruct{}
	sql.Params = "user_id,data,parameter"
	sql.Insert_value = fmt.Sprintf("%d,'%s','%s'", send_user_id, data, parameter)
	sql.Company_id = userId
	_, err := rpc.CallOther(Rpc, rpc.RpcMethod{
		Chinese_name: "消息",
		Method:       "MsgService.PostItem",
		Param:        sql.ToMap(),
	})
	if err != nil {
		logger.Error(err)
	}
	return err
}
