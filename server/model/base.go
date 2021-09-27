package model

import (
	"time"
)

var ExportModel = []interface{}{
	&Node{},
}

type base struct {
	CreatedAt time.Time `gorm:"comment:'创建时间'" json:"created_at" `
	UpdatedAt time.Time `gorm:"comment:'更新时间'" json:"updated_at"`
}
