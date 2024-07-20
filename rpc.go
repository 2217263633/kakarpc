package myrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/2217263633/kakarpc/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wonderivan/logger"
)

type RPC struct {
	Client    *rpc.Client `json:"client"`
	Count     int         `json:"count"`     // 重联计数器
	R         *gin.Engine `json:"r"`         // gin框架
	Conn      any         `json:"conn"`      // 连接注册中心
	Swag_port int         `json:"swag_port"` // swagger端口
}

// 是否存活
func (con *RPC) IsAlive(req string, res *bool) error {
	for rpcName := range RpcClient {
		_rpcName := strings.Split(rpcName, ".")
		if _rpcName[0] == req {
			RpcClient[rpcName].Heart = time.Now()
		}
	}
	*res = true
	return nil
}

// 注册
func (r *RPC) Register(req ServerStruct, res *ServerStruct) error {
	RpcServer[req.Chinese_name] = &YamlStruct{Server: req}

	time.AfterFunc(time.Second*5, func() {

		if RpcClient[req.Chinese_name] == nil {
			RpcClient[req.Chinese_name] = &RpcClientType{
				Heart: time.Now(),
				Addr:  "127.0.0.1:" + strconv.Itoa(req.Swag_port),
				Name:  strings.Split(req.Name, ".")[0],
			}
			logger.Info("转发服务", "/"+req.Name)
			centor.R.GET("/"+req.Name+"/*any", func(c *gin.Context) {
				target := "http://127.0.0.1:" + strconv.Itoa(req.Swag_port)
				url, _ := url.Parse(target)
				proxy := httputil.NewSingleHostReverseProxy(url)
				proxy.ServeHTTP(c.Writer, c.Request)
			})
			centor.R.POST("/"+req.Name+"/*any", func(c *gin.Context) {
				target := "http://127.0.0.1:" + strconv.Itoa(req.Swag_port)
				url, _ := url.Parse(target)
				proxy := httputil.NewSingleHostReverseProxy(url)
				proxy.ServeHTTP(c.Writer, c.Request)
			})
			centor.R.PUT("/"+req.Name+"/*any", func(c *gin.Context) {
				target := "http://127.0.0.1:" + strconv.Itoa(req.Swag_port)
				url, _ := url.Parse(target)
				proxy := httputil.NewSingleHostReverseProxy(url)
				proxy.ServeHTTP(c.Writer, c.Request)
			})
			centor.R.DELETE("/"+req.Name+"/*any", func(c *gin.Context) {
				target := "http://127.0.0.1:" + strconv.Itoa(req.Swag_port)
				url, _ := url.Parse(target)
				proxy := httputil.NewSingleHostReverseProxy(url)
				proxy.ServeHTTP(c.Writer, c.Request)
			})
		}

		for f := range RpcServer[req.Chinese_name].Server.Router {
			files := r.getConfigList()
			isFind := false
			for _, file := range files {
				if file.Name() == req.Chinese_name+".yaml" {
					isFind = true
					break
				}
			}

			if !isFind {
				os.Create("./config/" + req.Chinese_name + ".yaml")
				os.WriteFile("./config/"+req.Chinese_name+".yaml", []byte(fmt.Sprintf(
					"server:\n name: %s \n port: %d\n swag_port: %d \n path: %s \n mode: %d \n ",
					req.Name, req.Port, req.Swag_port, req.Path, req.Mode)), 0644)

				logger.Info("创建配置文件", "./config/"+req.Chinese_name+".yaml")
			}
			cli, err := rpc.DialHTTP("tcp", "127.0.0.1:"+strconv.Itoa(req.Port))
			if err == nil {
				var ff map[string]*rpc.Client = map[string]*rpc.Client{
					req.Name + "." + f: cli,
				}

				if RpcClient[req.Chinese_name].Client == nil {
					RpcClient[req.Chinese_name].Client = ff
				} else {
					RpcClient[req.Chinese_name].Client[req.Name+"."+f] = cli
				}

				logger.Info("连接服务成功", req.Chinese_name, req.Name+"."+f)
			} else {
				logger.Error("连接服务失败", req.Name+"."+f, err)
			}
		}
	})
	// logger.Info("注册服务", req)
	return nil
}

var initPort = 9100
var centor *RPC

func (r *RPC) CenterInit(_rpc *RPC) {
	gin.SetMode(gin.ReleaseMode)
	// conn := new(RPC)
	centor = _rpc
	rpc.Register(_rpc)
	rpc.HandleHTTP()
	_rpc.R.LoadHTMLFiles("./templates/index.html")

	_rpc.R.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "动态路由配置中心",
			"routerList": RpcClient,
		})
	})
	_rpc.Swag_port = 9101
	go _rpc.R.Run(fmt.Sprintf(":%d", _rpc.Swag_port))

	logger.Info("rpc server start at port: ", initPort)
	err := http.ListenAndServe(":"+strconv.Itoa(initPort), nil)
	if err != nil {
		logger.Error("error listening", err.Error())
		return
	}

}

func (r *RPC) init(conn any, port string) {
	rpc.Register(conn)
	rpc.HandleHTTP()
	logger.Info("rpc server start at port: ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		logger.Error("error listening", err.Error())
		return
	}
}

func (r *RPC) getConfigList() []fs.DirEntry {
	config_path := "./config"
	files, _ := os.ReadDir(config_path)
	return files
}

func (con *RPC) GetConfig(req string, res *map[string]interface{}) error {
	logger.Info("获取配置", req)
	if req == "" {
		(*res) = map[string]interface{}{
			"data": "127.0.0.1:1234"}
	} else if RpcServer[req] != nil {
		(*res) = RpcServer[req].Server.Router

	} else {
		(*res) = map[string]interface{}{
			"data": req + ": this service is not register"}
	}
	return nil
}

// 调用其他服务
func (con *RPC) Call(method RpcMethod, res *map[string]interface{}) error {
	logger.Info("调用服务", method.Chinese_name, method.Method, method.Param)
	// 因为编码原因 返回的 err 我们把它变成string  那边拿到后 做err 处理
	if RpcClient[method.Chinese_name] != nil && RpcClient[method.Chinese_name].Client[method.Method] != nil {
		err := RpcClient[method.Chinese_name].Client[method.Method].Call(method.Method, method.Param, res)
		if err != nil {
			(*res)["state"] = 401
			(*res)["err"] = err.Error()
			(*res)["data"] = []byte("[]")
		} else {
			if (*res)["err"] != nil {
				(*res)["state"] = 401
			} else {
				(*res)["state"] = 200
			}
		}
	} else {

		(*res)["state"] = 401
		(*res)["err"] = "this service is not online"
		(*res)["data"] = []byte("[]")
	}

	return nil
}

// 特定的调用
func (con *RPC) CallOther(method RpcMethod, res *map[string]interface{}) ([]map[string]interface{}, error) {
	var data map[string]interface{}
	con.Call(method, &data)

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

type Functype func(int)

// 连接注册中心
func (con *RPC) GoRpc(yaml *ServerStruct, _rpc *RPC) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9100")
	if err != nil {
		logger.Error("rpc.DialHTTP error: %v", err)
	} else {
		_rpc.Client = client

		structType := reflect.TypeOf(_rpc.Conn)

		for i := 0; i < structType.NumMethod(); i++ {
			method := structType.Method(i)
			yaml.Router[method.Name] = []string{method.Type.String()}
		}
		err = client.Call("RPC.Register", yaml, &yaml)
		if err != nil {
			logger.Error("rpc.Register error: %v", err)
		} else {
			logger.Info("连接注册中心成功", _rpc.Count)
		}
	}

	// 判断注册进程是否存活
	// 主进程不存活，则每隔10秒重连一次
	if _rpc.Count == 0 {
		_rpc.Count += 1
		ti := gocron.NewScheduler(time.UTC)
		logger.Info("判断主进程是否存在")
		ti.Every(10).Seconds().Do(func() {
			var reply bool
			if _rpc.Client == nil {
				logger.Info("主进程已掉线，等待重连")
				con.GoRpc(yaml, _rpc)
			} else {
				err := _rpc.Client.Call("RPC.IsAlive", yaml.Chinese_name, &reply)
				if err != nil {
					logger.Error("主进程已掉线，等待重连%v", err)
					con.GoRpc(yaml, _rpc)
				}
			}

		})
		ti.StartAsync()
		con.ginInit(_rpc.R, yaml.Swag_port, yaml.Name)
		con.init(_rpc.Conn, ":"+strconv.Itoa(yaml.Port))
	} else {
		_rpc.Count += 1
	}

}

func (con *RPC) ginInit(r *gin.Engine, port int, name string) {
	r.GET("/"+name+"/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DocExpansion("none")))
	logger.Info("gin run on port:", port)
	go r.Run(fmt.Sprintf(":%d", port))
}

var MyRpc = &RPC{}
var Utils *utils.CusUtils = utils.UtilsInit()
