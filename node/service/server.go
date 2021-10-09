package service

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	zutils "github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/znet"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/container/gmap"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shopspring/decimal"
	"io"
	"math/rand"
	"net"
	"p00q.cn/video_cdn/comm/model"
	"p00q.cn/video_cdn/comm/utils"
	"p00q.cn/video_cdn/node/config"
	"time"
)

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
	err = s.sendMsg(model.Authentication, msg2byte(model.Msg{
		SessionCode: 0,
		Err:         "",
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
	if msg.GetMsgId() != model.Authentication || m.Err != "" {
		_ = conn.Close()
		return errors.New(m.Err)
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
		if v.UsedPercent > 85 {
			clearCacheMap()
		}
	}
	percent, err := cpu.Percent(0, false)
	if err == nil {
		data.CPUPercent, _ = decimal.NewFromFloat(percent[0]).Round(2).Float64()
	}
	usage, err := disk.Usage(config.GlobalConfig.CacheDir)
	if err == nil {
		data.TotalDiskSpace = usage.Total
		data.AvailableDiskSpace = usage.Free
		data.DiskSpaceUsed = usage.Used
		if usage.UsedPercent > 80 {
			clear()
		}
	}
	data.Port = uint16(config.GlobalConfig.Port)
	data.Time = time.Now()
	data.Send = NetWorkState.Send
	data.Receive = NetWorkState.Receive
	return msg2byte(model.Msg{
		SessionCode: 0,
		Err:         "",
		Data:        data,
	})
}

// Run 运行服务
func Run() {
	//连接初始化
	zutils.GlobalObject.MaxPacketSize = 1024 * 1024 * 4
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

func msg2byte(m model.Msg) []byte {
	var bufferr bytes.Buffer
	PerEncod := gob.NewEncoder(&bufferr) //1.创建一个编码器
	err := PerEncod.Encode(&m)           //编码
	if err != nil {
		logs.Err(err)
	}
	return bufferr.Bytes()
}

func byte2Msg(data []byte) model.Msg {
	var msg model.Msg
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
	if whetherThereIsAWaitingRecipient(msg.GetMsgId(), m) {
		return
	}
	switch msg.GetMsgId() {
	case model.NewCacheData:
		updateCache(m.Data.([]model.Data))
	case model.DelayTest:
		//测速ping
		ping := utils.Ping(m.Data.(string))
		_ = server.sendMsg(model.DelayTest, msg2byte(model.Msg{
			SessionCode: m.SessionCode,
			Err:         "",
			Data: model.Delay{
				Host:   m.Data.(string),
				NodeIP: "",
				Val:    uint(ping),
			},
		}))
	}
}

var msgChanMap = gmap.New(true)

func sendAMessageAndWaitForAResponse(msgID uint32, msg model.Msg, duration time.Duration) model.Msg {
	server.sendMsg(msgID, msg2byte(msg))
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

func GetVideoCacheData(videoKey string) []model.Data {
	msg := sendAMessageAndWaitForAResponse(model.NewCacheData, model.Msg{
		SessionCode: rand.Uint64(),
		Err:         "",
		Data:        videoKey,
	}, time.Second*5)
	if msg.Err != "" {
		logs.Err(msg.Err)
		return make([]model.Data, 0)
	}
	return msg.Data.([]model.Data)
}
