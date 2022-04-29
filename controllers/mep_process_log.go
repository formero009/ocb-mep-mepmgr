package controllers

import (
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
)

type MepProcessLogController struct {
	base.BaseController
	manager.MepProcessLogManager
}

// @Title get mep process log
// @Description get mep process log
//@Param	page_size	query 	int		false	"the count of data in a page"
//@Param	page_no		query 	int		false	"the current page"
//@Param	start_time  query 	int		false	"the start unix timestamp of log"
//@Param	end_time	query 	int		false	"the end unix timestamp of log"
//@Param	mep_name     query 	string	 false	"mep name"
//@Param	log_type     query 	string	 false	"log type"
// @Success 204 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /process/log [get]
func (m *MepProcessLogController) Get() {
	q := mepmd.MepQuery{}
	q.QueryParam = *m.BuildQueryParam()
	start, err := m.GetInt64("start_time", 0)
	if err != nil {
		logs.Error("error start time: %v", err)
		m.AbortBadRequestFormat("start_time")
	}
	end, err := m.GetInt64("end_time", 0)
	if err != nil {
		logs.Error("error end time: %v", err)
		m.AbortBadRequestFormat("end_time")
	}

	q.MepName = m.GetString("mep_name")
	q.LogType = m.GetString("log_type")
	q.StartTime = start
	q.EndTime = end

	logs.Debug("query: %+v", q)
	count, data, err := m.MepProcessLogManager.Get(q)
	if err != nil {
		logs.Error("failed to get mep process data: %v", err)
		m.AbortInternalServerError(common.NewError(common.ErrDatabase))
	}

	m.Success(m.NewPage(count, data, int(q.PageNo), int(q.PageSize)))
}
