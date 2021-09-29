package global

import (
	"database/sql"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"p00q.cn/video_cdn/comm/model"
	"time"
)

var Cache = gcache.New()

//MySQL 定义MySql公共链接
var MySQL *gorm.DB
var Logs = g.Log()

func InitDB() {
	var ormConf = gorm.Config{
		SkipDefaultTransaction:                   false, //是否跳过默认事务
		NamingStrategy:                           nil,   //命名策略
		DryRun:                                   false, //生成 SQL 但不执行 针对整个系统
		PrepareStmt:                              true,  //数据库 预处理
		DisableAutomaticPing:                     true,  //检查数据库存活
		DisableForeignKeyConstraintWhenMigrating: false, //AutoMigrate 数据迁移的时候创建外键
		Logger: logger.New(
			log.New(Logs, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Microsecond, // 慢 SQL 阈值
				LogLevel:      logger.Error,     // Log level 可选值：Silent、Error、Warn、Info
				Colorful:      false,            // 禁用彩色打印
			},
		),
	}
	var err error
	if MySQL, err = gorm.Open(mysql.Open(g.Cfg().GetString("database.link")), &ormConf); err != nil {
		g.Log().Error(err.Error())
	}
	//下面是设置资源池 首先指定一个空的数据库进行设置连接数操作,初始化时生效一次
	var sqlDB *sql.DB
	if sqlDB, err = MySQL.DB(); err != nil {
		g.Log().Error(err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	initModel()
}

func initModel() {
	for _, i := range model.ExportModel {
		MySQL.AutoMigrate(i)
	}
}
