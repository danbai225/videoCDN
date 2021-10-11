package node

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	zutils "github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/os/gcache"
	"math/rand"
	"p00q.cn/video_cdn/comm/model"
	"p00q.cn/video_cdn/server/global"
	"strings"
	"time"
)

var nodeCache = gcache.New()
var server ziface.IServer
var nodeSet = gset.NewSet(true)
var nodeMap = make(map[string]model.Node)

func Run() {
	zutils.GlobalObject.Name = "videoCDNServer"
	zutils.GlobalObject.TcpPort = 7960
	zutils.GlobalObject.MaxPacketSize = 1024 * 1024 * 4

	//1 创建一个server句柄
	server = znet.NewServer()
	//2 配置路由
	server.AddRouter(model.PingPong, &pongRouter{})
	server.AddRouter(model.Authentication, &authenticationRouter{})
	server.AddRouter(model.NewCacheData, &newCacheDataRouter{})
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
	node := model.Node{}
	global.MySQL.Model(&node).Where("ip=?", ip).Take(&node)
	nodeMap[ip] = node
}
func stopHook(connection ziface.IConnection) {
	nodeSet.Remove(connection)
	ip := getIP(connection)
	global.MySQL.Model(&model.Node{}).Where("ip=?", ip).Update("on_line", false)
	global.Cache.Remove(fmt.Sprintf("ConnID-%s", ip))
}
func GetNodeInfoByIP(ip string) model.Node {
	if m, has := nodeMap[ip]; has {
		return m
	}
	return model.Node{}
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
		request.GetConnection().Stop()
		return
	}
	err = request.GetConnection().SendBuffMsg(model.PingPong, msg2byte(model.Msg{
		SessionCode: 0,
		Err:         "",
		Data:        nil,
	}))
	if err != nil {
		global.Logs.Error(err)
		request.GetConnection().Stop()
		return
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
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(model.Msg{
			SessionCode: 0,
			Err:         "数据len=0",
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}

	token := byte2Msg(data).Data.(string)
	var node model.Node
	err := global.MySQL.Model(&model.Node{}).Where("token=?", token).Take(&node).Error
	if err != nil {
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(model.Msg{
			SessionCode: 0,
			Err:         err.Error(),
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}
	if !strings.Contains(request.GetConnection().RemoteAddr().String(), node.IP) {
		_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(model.Msg{
			SessionCode: 0,
			Err:         "ip不存在",
			Data:        nil,
		}))
		request.GetConnection().Stop()
		return
	}
	_ = request.GetConnection().SendMsg(model.Authentication, msg2byte(model.Msg{
		SessionCode: 0,
		Err:         "",
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
func msg2byte(m model.Msg) []byte {
	var bufferr bytes.Buffer
	PerEncod := gob.NewEncoder(&bufferr) //1.创建一个编码器
	err := PerEncod.Encode(&m)           //编码
	if err != nil {
		global.Logs.Error(err)
	}
	return bufferr.Bytes()
}

func byte2Msg(data []byte) model.Msg {
	var msg model.Msg
	Decoder := gob.NewDecoder(bytes.NewReader(data)) //创建一个反编码器
	err := Decoder.Decode(&msg)
	if err != nil {
		global.Logs.Error(err)
	}
	return msg
}

func getNodeConnByIP(ip string) (ziface.IConnection, error) {
	id, err := global.Cache.Get(fmt.Sprintf("ConnID-%s", ip))
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, errors.New("null cache")
	}
	get, err := server.GetConnMgr().Get(id.(uint32))
	if err != nil {
		return nil, err
	}
	return get, nil
}

// NewCacheData 从数据库中查取缓存对应的url数据，给到节点
func NewCacheData(videoKey, ip string) {
	conn, err := getNodeConnByIP(ip)
	if err != nil {
		return
	}
	data := make([]model.Data, 0)
	global.MySQL.Model(&model.Data{}).Where("video_key=?", videoKey).Find(&data)
	conn.SendMsg(model.NewCacheData, msg2byte(model.Msg{
		SessionCode: 0,
		Err:         "",
		Data:        data,
	}))
}

//新缓存处理
type newCacheDataRouter struct {
	znet.BaseRouter
}

func (r *newCacheDataRouter) Handle(request ziface.IRequest) {
	if !verification(request.GetConnection()) {
		return
	}
	msg := byte2Msg(request.GetData())
	data := make([]model.Data, 0)
	global.MySQL.Model(&model.Data{}).Where("video_key=?", msg.Data.(string)).Find(&data)
	request.GetConnection().SendMsg(model.NewCacheData, msg2byte(model.Msg{
		SessionCode: msg.SessionCode,
		Err:         "",
		Data:        data,
	}))
}

// DelayTest 下发每个node对host测ping
func DelayTest(host string) {
	sid := rand.Uint64()
	nodeSet.Walk(func(item interface{}) interface{} {
		item.(ziface.IConnection).SendMsg(model.DelayTest, msg2byte(model.Msg{
			SessionCode: sid,
			Err:         "",
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
	cacheKey := fmt.Sprintf("delayTest-%d", msg.SessionCode)
	_, _ = nodeCache.GetOrSetFuncLock(cacheKey, func() (interface{}, error) {
		return delay, nil
	}, time.Minute)
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

var msgChanMap = gmap.New(true)

func sendAMessageAndWaitForAResponse(conn ziface.IConnection, msgID uint32, msg model.Msg, duration time.Duration) model.Msg {
	//conn.SendMsg(msgID,msg2byte(msg))
	tick := time.Tick(duration)
	msgC := make(chan model.Msg)
	msgChanMap.Set(fmt.Sprintf("%d%d", msgID, msg.SessionCode), msgC)
	defer func() {
		msgChanMap.Remove(fmt.Sprintf("%d%d", msgID, msg.SessionCode))
		close(msgC)
	}()
	select {
	case <-tick:
		return model.Msg{Err: "WaitTimeOUT"}
	case m := <-msgC:
		return m
	}
}
func whetherThereIsAWaitingRecipient(msgID uint32, msg model.Msg) bool {
	get := msgChanMap.Get(fmt.Sprintf("%d%d", msgID, msg.SessionCode))
	if get != nil {
		go func() {
			get.(chan model.Msg) <- msg
		}()
		return true
	}
	return false
}

//AssignANodeWithTheLeastLoad 分配一个负载最轻的节点
func AssignANodeWithTheLeastLoad() model.Node {
	node := model.Node{}
	global.MySQL.Model(&model.Node{}).Where("on_line=1").Order("cpu_percent ASC").First(&node)
	return node
}
