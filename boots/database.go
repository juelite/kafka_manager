package boots

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func init()  {
	time.LoadLocation("Local")

	// 连接数据库
	database_host := beego.AppConfig.String("DB_HOST")
	database_name := beego.AppConfig.String("DB_DATABASE")
	database_port := beego.AppConfig.String("DB_PORT")
	database_user := beego.AppConfig.String("DB_USERNAME")
	database_pwd := beego.AppConfig.String("DB_PASSWORD")
	orm.DefaultRowsLimit = -1

	Database(database_host, database_name, database_port, database_user, database_pwd, false, "default")
}

/**
 * 初始化数据库连接
 * @param database_host string 链接地址
 * @param database_name string 数据库名称
 * @param database_port string 数据库端口
 * @param database_user string 数据库用户
 * @param database_pwd string 数据库密码
 * @param debug bool true 开始调试 false 关闭调试
 */
func Database(database_host string , database_name string , database_port string , database_user string , database_pwd string , debug bool, conn_name string) {

	//选择模式
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.Debug = debug
	orm.DefaultTimeLoc = time.Local
	conn := database_user + ":" + database_pwd + "@tcp(" + database_host + ":" + database_port + ")/" + database_name + "?charset=utf8&parseTime=true&loc=Local"
	//注册数据库连接
	err := orm.RegisterDataBase(conn_name, "mysql", conn)
	if err != nil {
		panic(err)
	} else {
		//输出结果
		fmt.Printf("数据库连接成功！%s\n", conn)
	}

}
