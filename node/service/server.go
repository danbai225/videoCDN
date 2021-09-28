package service

import (
	"bytes"
	"encoding/json"

	"errors"

	"encoding/gob"
	"github.com/aceld/zinx/znet"
	logs "github.com/danbai225/go-logs"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shopspring/decimal"
	"io"
	"net"
	"p00q.cn/video_cdn/node/config"
	"p00q.cn/video_cdn/node/model"
	m3u8Server "p00q.cn/video_cdn/node/service/m3u8"
	"time"
)

const (
	pingPong = iota
	authentication
	newCache
	Friday
	Saturday
	Sunday
)

type Msg struct {
	SessionCode uint64      `json:"session_code"`
	Err         error       `json:"err"`
	Data        interface{} `json:"data"`
}

var server = serverConn{}

type serverConn struct {
	conn net.Conn
	Addr string
	dp   *znet.DataPack
}

//连接认证
func (s *serverConn) connect() error {
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.conn = conn
	s.dp = znet.NewDataPack()
	//开始认证
	err = s.sendMsg(authentication, msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        config.GlobalConfig.Token,
	}))
	if err != nil {
		return err
	}
	msg, err := s.readMsg()
	if err != nil {
		return err
	}
	m := byte2Msg(msg.GetData())
	if msg.GetMsgId() != authentication || m.Err != nil {
		_ = conn.Close()
		return m.Err
	}
	return nil
}

//读取消息
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
		return msg, nil
	}
	return nil, errors.New("GetDataLen() < 0")
}

//发送消息
func (s *serverConn) sendMsg(id uint32, data []byte) error {
	pack, err := s.dp.Pack(znet.NewMsgPackage(id, data))
	if err != nil {
		return err
	}
	_, err = s.conn.Write(pack)
	return err
}

// Ping 发送ping
func Ping() {
	if server.conn != nil {
		server.sendMsg(pingPong, PingData())
	}
}

// PingData ping数据
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
		data.CPUPercent, _ = decimal.NewFromFloat(percent[0]).Round(2).Float64()
	}
	usage, err := disk.Usage("./")
	if err == nil {
		data.TotalDiskSpace = usage.Total
		data.AvailableDiskSpace = usage.Free
		data.DiskSpaceUsed = usage.Used
	}
	data.Port = config.GlobalConfig.Port
	data.Time = time.Now()
	marshal, _ := json.Marshal(data)
	return msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        marshal,
	})
}

// Run 运行服务
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
	for err == nil {
		var msg *znet.Message
		msg, err = server.readMsg()
		if err != nil {
			logs.Err(err)
			break
		}
		go messageHandling(msg)
	}
}

//自定义消息转换
func msg2byte(m Msg) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m)
	if err != nil {
		logs.Err(err)
	}
	return b.Bytes()
}

//自定义消息转换
func byte2Msg(data []byte) Msg {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	var m Msg
	err := dec.Decode(&m)
	if err != nil {
		logs.Err(err)
	}
	return m
}

//消息处理
func messageHandling(msg *znet.Message) {
	m := byte2Msg(msg.GetData())
	switch msg.GetMsgId() {
	case newCache:
		rUrl, err := m3u8Server.NewTransit(m.Data.(string))
		server.sendMsg(newCache, msg2byte(Msg{
			SessionCode: m.SessionCode,
			Err:         err,
			Data:        rUrl,
		}))
	}
}
