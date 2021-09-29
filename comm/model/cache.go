package model

type Cache struct {
	base
	ID       uint   `gorm:"comment:'ID'" json:"id"`
	NodeIP   string `gorm:"comment:'size:32;NodeIP;index'" json:"node_ip"`
	Host     string `gorm:"size:128;comment:'源资源host ts文件'" json:"host"`
	NodeUrl  string `gorm:"size:256;comment:'缓存url'" json:"node_url"`
	Url      string `gorm:"size:256;comment:'被缓存的url;index'" json:"url"`
	Visits   uint   `gorm:"comment:'资源被访问次数'" json:"visits"`
	Valid    bool   `gorm:"comment:'是否有效';default:false" json:"valid"`
	Flow     uint64 `gorm:"comment:'流量'" json:"flow"`
	VideoKey string `gorm:"size:128;comment:'video标记';index" json:"video_key"`
}
