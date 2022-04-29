/**
 @author: yuefei
 @date: 2021/2/3
 @note:MEPM/MEP的操作日志
**/

package logmd

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var MepLog *logs.BeeLogger

func init() {
	MepLog = logs.NewLogger()
	logPath := beego.AppConfig.String("mep_log_path")
	maxsize := beego.AppConfig.DefaultInt("maxsize", 10) * 1024 * 1024
	maxDays := beego.AppConfig.DefaultInt("max_days", 7)
	logConf := fmt.Sprintf(`{"filename":"%s","maxsize":%d,"daily":true,"maxdays":%d,"rotate":true,"level":%d}`, logPath, maxsize, maxDays, logs.LevelInfo)
	err := MepLog.SetLogger(logs.AdapterFile, logConf)
	if err != nil {
		panic(err)
	}
}

func Log(mepName string, logType string) {
	MepLog.Info("mep_name=[%s] | log_type=[%s]", mepName, logType)
}
