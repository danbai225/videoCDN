package model

type Delay struct {
	base
	Host   string `gorm:"primaryKey;comment:'host地址'" json:"host"`
	NodeIP string `gorm:"size:32;index" json:"node_ip"`
	Val    uint   `gorm:"comment:'延迟但是ms'" json:"val"`
}
