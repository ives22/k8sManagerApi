package mysql

import (
	"fmt"
	"go.uber.org/zap"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8sManagerApi/config"
	"log"
	"os"
	"time"
)

// 定义全局的DB，是因为在其它包里面需要使用

var (
	isInit bool
	DB     *gorm.DB
)

func Init() {
	// 判断是否已经初始化了
	if isInit {
		return
	}

	//	连接MySQL数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Conf.MysqlInfo.Username,
		config.Conf.MysqlInfo.Password,
		config.Conf.MysqlInfo.Host,
		config.Conf.MysqlInfo.Port,
		config.Conf.MysqlInfo.Database,
		config.Conf.MysqlInfo.Charset)

	// SQL日志输出
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level  GORM 定义了这些日志级别：Silent、Error、Warn、Info
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	//fmt.Println("dsn:", dsn)
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("init mysql database failed, error: %s", err.Error()))
		panic("初始化MySQL连接失败," + err.Error())
	}
	// 创建表
	//db.AutoMigrate(model.Chart{})

	zap.L().Info("init mysql database connection successfully")

	DB = db
}