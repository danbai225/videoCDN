package model

type Delay struct {
	Base
	ID     uint   `gorm:"comment:'ID'" json:"id"`
	Host   string `gorm:"index:node_delay;comment:'host地址'" json:"host"`
	NodeIP string `gorm:"size:32;index:node_delay" json:"node_ip"`
	Val    uint   `gorm:"comment:'延迟单位ms'" json:"val"`
}
