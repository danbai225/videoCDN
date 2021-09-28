package model

import "time"

type PingData struct {
	UseOfMemory        uint64    `json:"use_of_memory"`        //使用的内存
	AvailableMemory    uint64    `json:"available_memory"`     //可用的内存
	TotalMemory        uint64    `json:"total_memory"`         //总运行内存 单位Bytes
	CPUPercent         float64   `json:"cpu_percent"`          //cpu使用率
	TotalDiskSpace     uint64    `json:"total_disk_space"`     //总磁盘空间
	DiskSpaceUsed      uint64    `json:"disk_space_used"`      //使用的磁盘空间
	AvailableDiskSpace uint64    `json:"available_disk_space"` //可使用的磁盘空间
	Port               int       `json:"port"`
	Time               time.Time `json:"time"`
}
