package cuswebsocket

import (
	"time"
)

// business
type BusinessType int

const (
	Chat        BusinessType = 0   // 聊天   0-99
	Dualprevent BusinessType = 100 // 双防  100-149 巡检点检   150-199 隐患排查
	Bigdanger   BusinessType = 200 // 重大危险源  200-299
	DailyList   BusinessType = 300 // 日常工作清单  300-399
	PostList    BusinessType = 400 // 岗位责任清单  400-499
	Work        BusinessType = 500 // 特殊作业 500-599
	Webrtc      BusinessType = 600 // webrtc 600-699
	Demind      BusinessType = 4   // button事件 =4
	CloseNotify BusinessType = 5   //	封闭化通知 = 5

	// 以下为业务自定义类型
	PushDailyRisk BusinessType = 10 // 推送日常风险
)

type WsMessage struct {
	Type     FishType     `json:"type"` // socket 交互使用的
	Data     interface{}  `json:"data"`
	UserId   string       `json:"userId"`
	Business BusinessType `json:"business"` // 0 chat  1 work 2 dualprevent 3 webrtc 4 demind 5 bigdanger
	CallUrl  interface{}  `json:"callUrl"`
}

type FishType int

const (
	Aicar FishType = 1
)

type Client struct {
	ID        int
	IpAddress string
	IpSource  string
	// Socket      *websocket.Conn
	Send        chan []byte
	Start       time.Time
	ExpireTime  time.Duration // 一段时间没有接收到心跳则过期
	UserId      int           // 用户ID
	Headimg     string        // 用户头像
	Person_sign string        // 用户签名
	Phone       string        // 用户电话
	Name        string        // 用户名称
	Business    FishType      `form:"business"`   // 业务ID,某些业务改了该值,才推送 aicar 1
	Company_id  int           `form:"company_id"` // 公司ID
}

type ClientManager struct {
	Clients    map[int]*Client // 记录在线用户
	Broadcast  chan []byte     // 触发消息广播
	SingleCast chan []byte     // 单信息
	SingId     string          //单消息的回复id
	Register   chan *Client    // 触发新用户登陆
	UnRegister chan *Client    // 触发用户退出
}

var Manager ClientManager = ClientManager{
	Clients:    make(map[int]*Client),
	Broadcast:  make(chan []byte),
	SingleCast: make(chan []byte),
	SingId:     "",
	Register:   make(chan *Client),
	UnRegister: make(chan *Client),
}

func InitManger() ClientManager {
	return Manager
}
