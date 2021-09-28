package service

import (
	netps "github.com/shirou/gopsutil/v3/net"
	"time"
)

var NetWorkState = struct {
	Send    uint64
	Receive uint64
}{}

func init() {
	go FlowCountStart()
}
func FlowCountStart() {
	for {
		n1, _ := netps.IOCounters(false)
		time.Sleep(time.Second)
		n2, _ := netps.IOCounters(false)
		if len(n1) > 0 && len(n2) > 0 {
			stat1 := n1[0]
			stat2 := n2[0]
			NetWorkState.Send = stat2.BytesSent - stat1.BytesSent
			NetWorkState.Receive = stat2.BytesRecv - stat1.BytesRecv
		}
	}
}
