package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//全局复用连接对象[dbConn]
var (
	dbConn *sql.DB
	err    error
)

//【dbops】包被调用时，自动初始化连接sql
func init() {
	dbConn, err = sql.Open("mysql", "root:mmmm@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	err = dbConn.Ping() //open方法不会验证数据库是否连接成功，需要ping测试
	if err != nil {
		log.Println("连接失败")
		return
	}
	log.Println("连接成功")
}
