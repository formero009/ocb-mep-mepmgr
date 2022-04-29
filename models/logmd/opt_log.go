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
	"time"
)

type OptLog struct {
	Account   string
	Name      string
	Role      string
	Ip        string
	StartTime time.Time
	Result    OptResult
	Exception string
	Content   string
	MepId     string
	UserId    string
}

type OptResult string

const (
	OptSuccess OptResult = "成功"
	OptFail    OptResult = "失败"
)

var log *logs.BeeLogger

func init() {
	log = logs.NewLogger()
	logPath := beego.AppConfig.String("opt_log_path")
	maxsize := beego.AppConfig.DefaultInt("maxsize", 10) * 1024 * 1024
	maxDays := beego.AppConfig.DefaultInt("max_days", 7)
	logConf := fmt.Sprintf(`{"filename":"%s","maxsize":%d,"daily":true,"maxdays":%d,"rotate":true,"level":%d}`, logPath, maxsize, maxDays, logs.LevelInfo)
	err := log.SetLogger(logs.AdapterFile, logConf)
	if err != nil {
		panic(err)
	}
}

func (optLog *OptLog) Log(result OptResult, exception, content, mepId string) {
	log.Info("account=[%s] | name=[%s] | role=[%s] | ip=[%s] | start_time=[%s] | ne_name=[%s] | opt_content=[%s] | result=[%s] | exception=[%s]",
		optLog.Account, optLog.Name, optLog.Role, optLog.Ip, optLog.StartTime.Format("2006/01/02 15:04:05.000"), mepId, content, result, exception)
}

func (optLog *OptLog) LogSuccess() {
	optLog.Log(OptSuccess, "", optLog.Content, optLog.MepId)
}

func (optLog *OptLog) LogFail(exception error) {
	optLog.Log(OptFail, exception.Error(), optLog.Content, optLog.MepId)
}
