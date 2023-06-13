package mysql

import (
	"fmt"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8sManagerApi/config"
)

// 定义全局的DB，是因为在其它包里面需要使用

var (
	isInit bool
	DB     *gorm.DB
	err    error
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

	fmt.Println("dsn:", dsn)
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("初始化MySQL连接失败," + err.Error())
	}
	// 创建表
	//db.AutoMigrate(model.Chart{})
	// 设置是否打印sql语句
	db.Logger.LogMode(logger.Silent)
	fmt.Println("初始化MySQL连接成功。")
	DB = db
}