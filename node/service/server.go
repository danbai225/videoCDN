package service

import (
	"encoding/json"
	"errors"
	"github.com/aceld/zinx/znet"
	logs "github.com/danbai225/go-logs"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	"net"
	"p00q.cn/video_cdn/node/config"
	"p00q.cn/video_cdn/node/model"
	"time"
)

var server = serverConn{}

type serverConn struct {
	conn net.Conn
	Addr string
	dp   *znet.DataPack
}

func (s *serverConn) connect() error {
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.conn = conn
	s.dp = znet.NewDataPack()
	//开始认证
	err = s.sendMsg(authentication, []byte(config.GlobalConfig.Token))
	if err != nil {
		return err
	}
	msg, err := s.readMsg()
	if err != nil {
		return err
	}
	if msg.GetMsgId() != authentication || string(msg.Data) != OK {
		_ = conn.Close()
		return errors.New("认证失败错误的code")
	}
	return nil
}
func (s *serverConn) readMsg() (*znet.Message, error) {
	//先读出流中的head部分
	headData := make([]byte, s.dp.GetHeadLen())
	//ReadFull 会把msg填充满为止
	_, err := io.ReadFull(s.conn, headData)
	if err != nil {
		logs.Err("read head error", err)
		return nil, err
	}
	//将headData字节流 拆包到msg中
	msgHead, err := s.dp.Unpack(headData)
	if err != nil {
		logs.Err("server unpack err:", err)
		return nil, err
	}
	if msgHead.GetDataLen() > 0 {
		//msg 是有data数据的，需要再次读取data数据
		msg := msgHead.(*znet.Message)
		msg.Data = make([]byte, msg.GetDataLen())
		//根据dataLen从io中读取字节流
		_, err = io.ReadFull(s.conn, msg.Data)
		if err != nil {
			logs.Err("server unpack data err:", err)
			return nil, err
		}
		logs.Info("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		return msg, nil
	}
	return nil, errors.New("GetDataLen() < 0")
}
func (s *serverConn) sendMsg(id uint32, data []byte) error {
	pack, err := s.dp.Pack(znet.NewMsgPackage(id, data))
	if err != nil {
		return err
	}
	_, err = s.conn.Write(pack)
	return err
}

func Ping() {
	if server.conn != nil {
		server.sendMsg(0, PingData())
	}
}
func PingData() []byte {
	data := model.PingData{}
	v, err := mem.VirtualMemory()
	if err == nil {
		data.TotalMemory = v.Total
		data.UseOfMemory = v.Used
		data.AvailableMemory = v.Free
	}
	percent, err := cpu.Percent(0, false)
	if err == nil {
		data.CPUPercent = percent[0]
	}
	usage, err := disk.Usage("./")
	if err == nil {
		data.TotalDiskSpace = usage.Total
		data.AvailableDiskSpace = usage.Free
		data.DiskSpaceUsed = usage.Used
	}
	data.Time = time.Now()
	marshal, _ := json.Marshal(data)
	logs.Info(string(marshal))
	return marshal
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

func Run() {
	//连接初始化
	server.Addr = config.GlobalConfig.ServerAddress
	logs.Info("开始连接服务端...")
	err := server.connect()
	if err != nil {
		logs.Info("连接服务端失败", err)
		return
	}
	logs.Info("连接成功...")
	var msg *znet.Message
	for err == nil {
		msg, err = server.readMsg()
		if err != nil {
			logs.Err(err)
			break
		}
		switch msg.GetMsgId() {
		case pong:
			logs.Info("pong")
		}
	}
}
