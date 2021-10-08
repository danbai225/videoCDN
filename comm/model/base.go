package model

import (
	"encoding/gob"
	"time"
)

var ExportModel = []interface{}{
	&Node{},
	&Cache{},
	&Delay{},
	&Data{},
}

type Base struct {
	CreatedAt time.Time `gorm:"comment:'创建时间'" json:"created_at" `
	UpdatedAt time.Time `gorm:"comment:'更新时间'" json:"updated_at"`
}

func init() {
	gob.RegisterName("Node", Node{})
	gob.RegisterName("Delay", Delay{})
	gob.RegisterName("Data", Data{})
	gob.RegisterName("DataArray", []Data{})
	gob.RegisterName("Msg", Msg{})
}
