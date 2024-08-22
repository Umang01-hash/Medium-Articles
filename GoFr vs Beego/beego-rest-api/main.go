package main

import (
	_ "beego-rest-api/routers"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:password@tcp(127.0.0.1:2001)/test?charset=utf8")
}

func main() {
	beego.BConfig.RunMode = "dev"
	beego.Run()
}
