package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/tylerjnettleton/motivator/models"
	_ "github.com/tylerjnettleton/motivator/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:tyleristhebest@/motivator?charset=utf8")
}

func main() {

	l := logs.GetLogger()
	l.Println("this is a message of http")

	// Create a new ORM instance
	o := orm.NewOrm()
	o.Using("default") // Using default, you can use other database

	name := "default"

	// Drop table and re-create.
	force := false

	// Print log.
	verbose := true

	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	orm.Debug = true

	// Start the web service
	beego.Run()
}
