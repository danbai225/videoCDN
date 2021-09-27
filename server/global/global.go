package global

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
)

var Cache = gcache.New()

func DB() gdb.DB {
	return g.DB()
}
