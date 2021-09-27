package node

import (
	"encoding/json"
	zutils "github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/model"
	"strings"
)

func Run() {
	zutils.GlobalObject.Name = "videoCDNServer"
	zutils.GlobalObject.TcpPort = 796

	//1 创建一个server句柄
	s := znet.NewServer()
	//2 配置路由
	s.AddRouter(ping, &pongRouter{})
	s.AddRouter(authentication, &authenticationRouter{})
	//3 开启服务
	s.Serve()
}

const (
	ping = iota
	pong
	authentication
	Thursday
	Friday
	Saturday
	Sunday
)
const (
	OK  = "ok"
	ERR = "err"
)

var ERRByte = []byte{101, 114, 114}
var OKByte = []byte{111, 107}

type pongRouter struct {
	znet.BaseRouter
}

func (r *pongRouter) Handle(request ziface.IRequest) {
	if !verification(request.GetConnection()) {
		return
	}
	data := request.GetData()
	node := model.Node{}
	_ = json.Unmarshal(data, &node)
	err := global.MySQL.Debug().Where("ip=?", getIP(request.GetConnection())).Updates(&node).Error
	if err != nil {
		global.Logs.Error(err)
	}
	err = request.GetConnection().SendBuffMsg(pong, []byte("pong"))
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
		request.GetConnection().SendMsg(authentication, ERRByte)
		request.GetConnection().Stop()
		return
	}
	token := string(data)
	var node model.Node
	err := global.MySQL.Model(&model.Node{}).Where("token=?", token).Take(&node).Error
	if err != nil {
		request.GetConnection().SendMsg(authentication, ERRByte)
		request.GetConnection().Stop()
		return
	}

	if !strings.Contains(request.GetConnection().RemoteAddr().String(), node.IP) {
		request.GetConnection().SendMsg(authentication, ERRByte)
		request.GetConnection().Stop()
		return
	}
	request.GetConnection().SendMsg(authentication, OKByte)
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
	return property.(string)
}
