package node

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	zutils "github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/os/gcache"
	"math/rand"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/model"
	"strings"
	"time"
)

var nodeCache = gcache.New()
var server ziface.IServer
var nodeSet = gset.NewSet(true)

func Run() {
	zutils.GlobalObject.Name = "videoCDNServer"
	zutils.GlobalObject.TcpPort = 7960

	//1 创建一个server句柄
	server = znet.NewServer()
	//2 配置路由
	server.AddRouter(model.PingPong, &pongRouter{})
	server.AddRouter(model.Authentication, &authenticationRouter{})
	server.AddRouter(model.NewCache, &newCacheRouter{})
	server.AddRouter(model.DelayTest, &delayTestRouter{})
	server.SetOnConnStop(stopHook)
	server.SetOnConnStart(startHook)
	//3 开启服务
	server.Serve()
}

func startHook(connection ziface.IConnection) {
	nodeSet.Add(connection)
	ip := getIP(connection)
	global.MySQL.Model(&model.Node{}).Where("ip=?", ip).Update("on_line", true)
	global.Cache.Set(fmt.Sprintf("ConnID-%s", ip), connection.GetConnID(), 0)
}
func stopHook(connection ziface.IConnection) {
	nodeSet.Remove(connection)
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
	msg := byte2Msg(data)
	node := msg.Data.(model.Node)
	node.OnLine = true
	ip := getIP(request.GetConnection())
	err := global.MySQL.Where("ip=?", ip).Updates(&node).Error
	if err != nil {
		global.Logs.Error(err)
	}
	err = request.GetConnection().SendBuffMsg(model.PingPong, msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        nil,
	}))
	if err != nil {
		global.Logs.Error(err)
	}
	setNodeRate(ip, node.Send, node.Receive)
}

//认证处理
type authenticationRouter struct {
	znet.BaseRouter
}

func (r *authenticationRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	if len(data) == 0 {
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(Msg{
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
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(Msg{
			SessionCode: 0,
			Err:         err,
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}

	if !strings.Contains(request.GetConnection().RemoteAddr().String(), node.IP) {
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(Msg{
			SessionCode: 0,
			Err:         errors.New("ip不存在"),
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}
	_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(Msg{
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
	var bufferr bytes.Buffer
	PerEncod := gob.NewEncoder(&bufferr) //1.创建一个编码器
	err := PerEncod.Encode(&m)           //编码
	if err != nil {
		global.Logs.Error(err)
	}
	return bufferr.Bytes()
}

func byte2Msg(data []byte) Msg {
	var msg Msg
	Decoder := gob.NewDecoder(bytes.NewReader(data)) //创建一个反编码器
	err := Decoder.Decode(&msg)
	if err != nil {
		global.Logs.Error(err)
	}
	return msg
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
	conn.SendMsg(model.NewCache, msg2byte(Msg{
		SessionCode: SessionCode,
		Err:         nil,
		Data:        url,
	}))
	cacheKey := fmt.Sprintf("newCache-%d", SessionCode)
	nodeCache.Set(cacheKey, nil, time.Minute)
	for {
		time.Sleep(time.Millisecond * 50)
		get, err2 := nodeCache.Get(cacheKey)
		if err2 != nil {
			return "", err2
		}
		if get != nil {
			nodeCache.Remove(cacheKey)
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
	nodeCache.Set(cacheKey, msg, time.Minute)
}

// DelayTest 下发每个node对host测ping
func DelayTest(host string) {
	nodeSet.Walk(func(item interface{}) interface{} {
		item.(ziface.IConnection).SendMsg(model.DelayTest, msg2byte(Msg{
			SessionCode: 0,
			Err:         nil,
			Data:        host,
		}))
		return item
	})
}

//测ping处理
type delayTestRouter struct {
	znet.BaseRouter
}

func (r *delayTestRouter) Handle(request ziface.IRequest) {
	if !verification(request.GetConnection()) {
		return
	}
	ip := getIP(request.GetConnection())
	msg := byte2Msg(request.GetData())
	delay := msg.Data.(model.Delay)
	delay.NodeIP = ip
	var count int64
	global.MySQL.Model(&model.Delay{}).Where("node_ip=? AND host=?", ip, delay.Host).Count(&count)
	if count > 0 {
		global.MySQL.Model(&model.Delay{}).Where("node_ip=? AND host=?", ip, delay.Host).Update("val", delay.Val)
	} else {
		global.MySQL.Model(&model.Delay{}).Create(&delay)
	}
}

//节点速率相关方法
type rate struct {
	Send    uint64 `json:"send"`    //最近1s所发送的字节数
	Receive uint64 `json:"receive"` //最近1s所接收的字节数
}

func setNodeRate(ip string, send, receive uint64) {
	key := fmt.Sprintf("rate-%s", ip)
	//contains, err := nodeCache.Contains(key)
	//if !contains||err!=nil {
	//	nodeCache.Set()
	//}
	nodeCache.Set(key, rate{
		Send:    send,
		Receive: receive,
	}, 0)
}
func getNodeRate(ip string) rate {
	key := fmt.Sprintf("rate-%s", ip)
	get, err := nodeCache.Get(key)
	if err != nil {
		return rate{}
	}
	return get.(rate)
}
