package main

import (
	_ "api/routers"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	databaseconf := "database_" + beego.BConfig.RunMode
	dbHost := beego.AppConfig.String(databaseconf + "::host")
	dbPort := beego.AppConfig.String(databaseconf + "::port")
	dbUser := beego.AppConfig.String(databaseconf + "::user")
	dbPass := beego.AppConfig.String(databaseconf + "::passwd")
	dbName := beego.AppConfig.String(databaseconf + "::dbname")
	dbCharset := beego.AppConfig.String(databaseconf + "::charset")

	conn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=" + dbCharset
	err := orm.RegisterDataBase("default", "mysql", conn)
	if err != nil {
		fmt.Println("database error ", err.Error())
		log.Fatal(err)
	}
	//utils.SetDB(orm.GetDB())
	fmt.Printf("数据库连接成功！%s\n", conn)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {
	var currentDirectory string
	currentDirectory = getCurrentDirectory()

	// log := logs.NewLogger(10000)
	// defer log.Close()
	// // 设置配置文件
	// jsonConfig := `{
	// 	"filename" : "./logs/log.log",
	// 	"maxlines" : 1000,
	// 	"maxsize"  : 10240
	// }`
	// log.SetLogger("file", jsonConfig) // 设置日志记录方式：本地文件记录
	// log.EnableFuncCallDepth(true)     // 输出log时能显示输出文件名和行号（非必须）
	// beego.SetLevel(beego.LevelInformational)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = currentDirectory + "/swagger"
	}

	beego.Run()
}
