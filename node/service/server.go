package service

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/aceld/zinx/znet"
	logs "github.com/danbai225/go-logs"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shopspring/decimal"
	"io"
	"net"
	"p00q.cn/video_cdn/node/config"
	m3u8Server "p00q.cn/video_cdn/node/service/m3u8"
	"p00q.cn/video_cdn/node/utils"
	"p00q.cn/video_cdn/server/model/model"
	"time"
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
	err = s.sendMsg(model.Authentication, msg2byte(Msg{
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
	if msg.GetMsgId() != model.Authentication || m.Err != nil {
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
		_ = server.sendMsg(model.PingPong, PingData())
	}
}

// PingData ping数据
func PingData() []byte {
	data := model.Node{}
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
	data.Port = uint16(config.GlobalConfig.Port)
	data.Time = time.Now()
	data.Send = NetWorkState.Send
	data.Receive = NetWorkState.Receive
	return msg2byte(Msg{
		SessionCode: 0,
		Err:         nil,
		Data:        data,
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

func msg2byte(m Msg) []byte {
	var bufferr bytes.Buffer
	PerEncod := gob.NewEncoder(&bufferr) //1.创建一个编码器
	err := PerEncod.Encode(&m)           //编码
	if err != nil {
		logs.Err(err)
	}
	return bufferr.Bytes()
}

func byte2Msg(data []byte) Msg {
	var msg Msg
	Decoder := gob.NewDecoder(bytes.NewReader(data)) //创建一个反编码器
	err := Decoder.Decode(&msg)
	if err != nil {
		logs.Err(err)
	}
	return msg
}

//消息处理
func messageHandling(msg *znet.Message) {
	m := byte2Msg(msg.GetData())
	switch msg.GetMsgId() {
	case model.NewCache:
		rUrl, err := m3u8Server.NewTransit(m.Data.(string))
		_ = server.sendMsg(model.NewCache, msg2byte(Msg{
			SessionCode: m.SessionCode,
			Err:         err,
			Data:        rUrl,
		}))
	case model.DelayTest:
		//测速ping
		ping := utils.Ping(m.Data.(string))
		_ = server.sendMsg(model.DelayTest, msg2byte(Msg{
			SessionCode: m.SessionCode,
			Err:         nil,
			Data: model.Delay{
				Host:   m.Data.(string),
				NodeIP: "",
				Val:    uint(ping),
			},
		}))
	}
}
