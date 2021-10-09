package utils

import (
	"github.com/go-ping/ping"
	"log"
)

func Ping(addr string) int64 {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		log.Println(err)
		return 999
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Println(err)
		return 999
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	if err != nil {
		log.Println(err)
		return 999
	}
	return stats.AvgRtt.Milliseconds()
}
