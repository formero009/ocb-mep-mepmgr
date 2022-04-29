package main

import (
	"flag"
	"mepmgr/config"
	"mepmgr/manager"
	_ "mepmgr/models"
	_ "mepmgr/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func loadConfig() { //C:\Users\caigui\go\src\mep\mepserver\conf\app.conf  ./conf/app.conf
	configpath := flag.String("config", "./conf/app.conf", "config:default is app.conf")
	flag.Parse()
	beego.LoadAppConfig("ini", *configpath)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}

func main() {
	loadConfig()
	//note: beego.LoadAppConfig方法会执行beego日志的reset，必须将默认日志的手动初始化置于beego.loadAppConfig之后
	config.InitSysLog()
	manager.InitConfig()
	manager.InitChecker()
	if err := manager.LoadCerts(); err != nil {
		logs.Error("init mepm certs error, %s", err.Error())
		return
	}
	go manager.DefaultMepMg.NodeChecker.Start()
	beego.Run()
	manager.DefaultMepMg.NodeChecker.Stop()
}
