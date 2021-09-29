package model

type Data struct {
	base
	Key      string `gorm:"primaryKey;comment:'键名'" json:"key"`
	VideoKey string `gorm:"size:128;comment:'video标记';index" json:"video_key"`
	Type     string `gorm:"size:128;comment:'类型'" json:"type"`
	Data     string `gorm:"comment:'数据'" json:"data"`
}
