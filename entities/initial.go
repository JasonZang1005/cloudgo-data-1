package entities

import "github.com/go-xorm/xorm"

var engine *xorm.Engine
var err error

func init() {
	engine, err = xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	checkErr(err)
	engine.Sync2(new(UserInfo))
}
