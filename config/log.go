/*
 */

// beego logs config
package config

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func InitSysLog() {
	//logConf record code file and line
	logs.EnableFuncCallDepth(true)

	//enable async logConf
	//logs.Async()

	//get logConf conf
	logPath := beego.AppConfig.String("sys_log_path")
	maxsize := beego.AppConfig.DefaultInt("maxsize", 10) * 1024 * 1024
	maxDays := beego.AppConfig.DefaultInt("max_days", 7)

	//set common logger that saves all logConf
	logConf := fmt.Sprintf(`{"filename":"%s","maxsize":%d,"daily":true,"maxdays":%d,"rotate":true,"level":%d}`, logPath, maxsize, maxDays, logs.LevelInfo)
	err := logs.SetLogger(logs.AdapterFile, logConf)
	if err != nil {
		panic(err)
	}
}
