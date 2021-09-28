package node

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	zutils "github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/gogf/gf/os/gcache"
	"math/rand"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/model"
	"strings"
	"time"
)

var callMap = gcache.New()
var server ziface.IServer

func Run() {
	zutils.GlobalObject.Name = "videoCDNServer"
	zutils.GlobalObject.TcpPort = 7960

	//1 创建一个server句柄
	server = znet.NewServer()
	//2 配置路由
	server.AddRouter(pingPong, &pongRouter{})
	server.AddRouter(authentication, &authenticationRouter{})
	server.AddRouter(newCache, &newCacheRouter{})
	server.SetOnConnStop(stopHook)
	server.SetOnConnStart(startHook)
	//3 开启服务
	server.Serve()
}

const (
	pingPong = iota
	authentication
	newCache
	Friday
	Saturday
	Sunday
)

func startHook(connection ziface.IConnection) {
	ip := getIP(connection)
	global.MySQL.Model(&model.Node{}).Where("ip=?", ip).Update("on_line", true)
	global.Cache.Set(fmt.Sprintf("ConnID-%s", ip), connection.GetConnID(), 0)
}
func stopHook(connection ziface.IConnection) {
	ip := getIP(connection)
	global.MySQL.Model(&model.Node{}).Where("ip=?", ip).Update("on_line", false)
	global.Cache.Remove(fmt.Sprintf("ConnID-%s", ip))
}

type pongRouter struct {
	znet.BaseRouter
}

func (r *pongRouter) Handle(request ziface.IRequest) {
	if !verification(request.GetConnection()) {
		return
	}
	data := request.GetData()
	node := model.Node{}
	_ = json.Unmarshal(byte2Msg(data).Data.([]byte), &node)
	node.OnLine = true
	err := global.MySQL.Where("ip=?", getIP(request.GetConnection())).Updates(&node).Error
	if err != nil {
		global.Logs.Error(err)
	}
	err = request.GetConnection().SendBuffMsg(pingPong, msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        nil,
	}))
	if err != nil {
		global.Logs.Error(err)
	}
}

//认证处理
type authenticationRouter struct {
	znet.BaseRouter
}

func (r *authenticationRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	if len(data) == 0 {
		_ = request.GetConnection().SendMsg(authentication, msg2byte(Msg{
			SessionCode: 0,
			Err:         errors.New("数据len=0"),
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}

	token := byte2Msg(data).Data.(string)
	var node model.Node
	err := global.MySQL.Model(&model.Node{}).Where("token=?", token).Take(&node).Error
	if err != nil {
		_ = request.GetConnection().SendMsg(authentication, msg2byte(Msg{
			SessionCode: 0,
			Err:         err,
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}

	if !strings.Contains(request.GetConnection().RemoteAddr().String(), node.IP) {
		_ = request.GetConnection().SendMsg(authentication, msg2byte(Msg{
			SessionCode: 0,
			Err:         errors.New("ip不存在"),
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}
	_ = request.GetConnection().SendMsg(authentication, msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        nil,
	}))
	request.GetConnection().SetProperty("auth", true)
	request.GetConnection().SetProperty("ip", node.IP)
}

//判断是否完成认证
func verification(conn ziface.IConnection) bool {
	property, err := conn.GetProperty("auth")
	if err != nil {
		conn.Stop()
		return false
	}
	if property.(bool) != true {
		conn.Stop()
		return false
	}
	return true
}
func getIP(conn ziface.IConnection) string {
	property, _ := conn.GetProperty("ip")
	if property == nil {
		split := strings.Split(conn.RemoteAddr().String(), ":")
		return split[0]
	}
	return property.(string)
}
func msg2byte(m Msg) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m)
	if err != nil {
		global.Logs.Error(err)
	}
	return b.Bytes()
}

func byte2Msg(data []byte) Msg {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	var m Msg
	err := dec.Decode(&m)
	if err != nil {
		global.Logs.Error(err)
	}
	return m
}

type Msg struct {
	SessionCode uint64      `json:"session_code"`
	Err         error       `json:"err"`
	Data        interface{} `json:"data"`
}

func getNodeConnByIP(ip string) (ziface.IConnection, error) {
	id, err := global.Cache.Get(fmt.Sprintf("ConnID-%s", ip))
	if err != nil {
		return nil, err
	}
	get, err := server.GetConnMgr().Get(id.(uint32))
	if err != nil {
		return nil, err
	}
	return get, nil
}
func NewCache(url, ip string) (string, error) {
	conn, err := getNodeConnByIP(ip)
	if err != nil {
		return "", err
	}
	SessionCode := rand.Uint64()
	conn.SendMsg(newCache, msg2byte(Msg{
		SessionCode: SessionCode,
		Err:         nil,
		Data:        url,
	}))
	cacheKey := fmt.Sprintf("newCache-%d", SessionCode)
	callMap.Set(cacheKey, nil, time.Minute)
	for {
		time.Sleep(time.Millisecond * 50)
		get, err2 := callMap.Get(cacheKey)
		if err2 != nil {
			return "", err2
		}
		if get != nil {
			callMap.Remove(cacheKey)
			return get.(Msg).Data.(string), get.(Msg).Err
		}
	}
}

//新缓存处理
type newCacheRouter struct {
	znet.BaseRouter
}

func (r *newCacheRouter) Handle(request ziface.IRequest) {
	if !verification(request.GetConnection()) {
		return
	}
	msg := byte2Msg(request.GetData())
	cacheKey := fmt.Sprintf("newCache-%d", msg.SessionCode)
	callMap.Set(cacheKey, msg, time.Minute)
}
