package pkg

import (
	"fmt"
	"sf-paas/config"
	"sf-paas/middleware"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Sql *gorm.DB

var log = middleware.Logger

//func InitDB() error {
func init() {
	var err error

	// 数据库的类型
	dbType := config.Conf.Mysql.Type

	// Mysql配置信息
	mysqlName := config.Conf.Mysql.Dbname
	mysqlUser := config.Conf.Mysql.User
	mysqlPwd := config.Conf.Mysql.Pwd
	mysqlPort := strconv.Itoa(config.Conf.Mysql.Port)
	mysqlCharset := config.Conf.Mysql.Dbcharset
	mysqlHost := config.Conf.Mysql.Host

	// sqlite3配置信息
	//sqliteName := conf.Section(DB_SQLITE3_TYPE).Key("db_name").String()

	var dataSource string
	//switch dbType {
	//case DB_MYSQL_TYPE:
	dataSource = mysqlUser + ":" + mysqlPwd + "@tcp(" + mysqlHost + ":" +
		mysqlPort + ")/" + mysqlName + "?charset=" + mysqlCharset +
		"&parseTime=" + "true" + "&loc=" + "Local"
	//log.Info("dbType:%v\n", dbType)
	//log.Info("dataSource:%v\n", dataSource)
	Sql, err = gorm.Open(dbType, dataSource)
	//case DB_SQLITE3_TYPE:
	//	dataSource = "database" + string(os.PathSeparator) + sqliteName
	//	if !gfile.Exists(dataSource) {
	//		os.MkdirAll(path.Dir(dataSource), os.ModePerm)
	//		os.Create(dataSource)
	//	}
	//	db, err = gorm.Open(dbType, dataSource)
	//}
	if Sql == nil {
		fmt.Sprintf("db engine invalid")
		return
	}
	if err != nil {
		fmt.Sprintf("connect mysql failed:%v\n", err)
		//db.Close()
		return
	}

	//defer db.Close()
	// 设置连接池，空闲连接
	Sql.DB().SetMaxIdleConns(config.Conf.Mysql.MaxIdleConns)
	// 打开链接
	Sql.DB().SetMaxOpenConns(config.Conf.Mysql.MaxOpenConns)
	//连接超时
	Sql.DB().SetConnMaxLifetime(time.Second * time.Duration(config.Conf.Mysql.MaxConnLifeTime))
	Sql.SetLogger(middleware.Logger)
	// 表明禁用后缀加s
	Sql.SingularTable(true)

	fmt.Printf("init db success\n")
}
