package model

import "time"

type Node struct {
	base
	ID                 uint      `gorm:"comment:'ID'" json:"id"`
	IP                 string    `gorm:"size:32;comment:'IP';uniqueIndex" json:"ip"`
	Port               uint16    `gorm:"comment:'服务端口'" json:"port"`
	Domain             string    `gorm:"size:128comment:'域名'" json:"domain"`
	Area               string    `gorm:"size:128comment:'地区'" json:"area"`
	Bandwidth          uint      `gorm:"comment:'带宽'" json:"bandwidth"`
	Token              string    `gorm:"size:128comment:'访问令牌'" json:"token"`
	UseOfMemory        uint64    `gorm:"comment:'已使用的内存'" json:"use_of_memory"`
	AvailableMemory    uint64    `gorm:"comment:'可使用的内存'" json:"available_memory"`
	TotalMemory        uint64    `gorm:"comment:'总运行内存'" json:"total_memory"`
	CPUPercent         float64   `gorm:"comment:'cpu使用率'" json:"cpu_percent"`
	TotalDiskSpace     uint64    `gorm:"comment:'总磁盘空间'" json:"total_disk_space"`
	DiskSpaceUsed      uint64    `gorm:"comment:'已使用的磁盘空间'" json:"disk_space_used"`
	AvailableDiskSpace uint64    `gorm:"comment:'可使用的磁盘空间'" json:"available_disk_space"`
	OnLine             bool      `gorm:"comment:'在线状态';default:false" json:"on_line"`
	Time               time.Time `gorm:"-" json:"time"`
}
