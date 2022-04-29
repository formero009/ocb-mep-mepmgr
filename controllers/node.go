package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
	"regexp"
)

type NodeController struct {
	base.BaseController
}

// @Title Create mep node
// @Description Create mep node
// @Param       mepMeta  body   mepmd.MepMeta   true   "mep meta info"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router / [post]
func (s *NodeController) Create() {
	logs.Debug("create mep")

	mepCfg := mepmd.MepMeta{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &mepCfg); err != nil {
		logs.Error("parse body %s error. %v", string(s.Ctx.Input.RequestBody), err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}

	err := manager.DefaultMepMg.CheckReq(&mepCfg)
	if err != nil {
		s.AbortBadRequest(err)
	}

	//FormatHttpsFile(s, &mepCfg)
	if err := manager.DefaultMepMg.Create(&mepCfg, s.OptLog.UserId); err != nil {
		logs.Error("create mep meta for %s fail %s", mepCfg, err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(mepCfg)
}

// @Title Update mep node
// @Description Update mep node
// @Param       mepMeta  body   mepmd.MepMeta true   "mep info"
// @Param       mepId    path    string        true   "mep id"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /:mepId [put]
func (s *NodeController) Update() {
	logs.Debug("update mep")

	mepId := s.Ctx.Input.Param(":mepId")
	mepCfg := mepmd.MepMeta{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &mepCfg); err != nil {
		logs.Error("parse body %s error. %v", string(s.Ctx.Input.RequestBody), err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}

	mepCfg.MepId = mepId
	err := manager.DefaultMepMg.CheckReq(&mepCfg)
	if err != nil {
		s.AbortBadRequest(err)
	}

	//FormatHttpsFile(s, &mepCfg)
	if err := manager.DefaultMepMg.Update(&mepCfg); err != nil {
		logs.Error("create mep meta for %s fail %s", mepCfg, err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(mepCfg)
}

// @Title Describe mep node
// @Description Get mep node information
// @Param	currentPage	    query 	int 	 false	"page number"
// @Param	pageSize	query 	int 	 false	"page size"
// @Param	MEPMID  	query 	string	 false	"MEPMID"
// @Param	CALLERId  	query 	string	 false	"CALLERId"
// @Param	mepGroupId  query 	int64    false	""
// @Param	mepName  	query 	bool	 false	""
// @Param	mepId  	    query 	string	 false	""
// @Param	endPoint  	query 	string	 false	""
// @Param	notMepGroupId  	query 	int64	 false	""
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router  / [get]
func (s *NodeController) List() {
	logs.Debug("list mep")
	pageSize, _ := s.GetInt("pageSize", common.DefaultPageSize)
	pageNo, _ := s.GetInt("currentPage", common.DefaultPageNumber)

	bean := mepmd.MepMeta{}
	if val := s.GetString("MEPMID"); len(val) != 0 {
		bean.MEPMID = val
	}
	if val := s.GetString("CALLERId"); len(val) != 0 {
		bean.CALLERId = val
	}

	mepGroupId, _ := s.GetInt64("mepGroupId", 0)
	notMepGroupId, _ := s.GetInt64("notMepGroupId", 0)

	mepName := s.GetString("mepName")
	mepId := s.GetString("mepId")
	endPoint := s.GetString("endPoint")

	logs.Debug("mepName %s, mepId %s, endPoint %s", mepName, mepId, endPoint)
	re := regexp.MustCompile("[~!@#$%^&*(){}|<>\\\\/+\\-=【】:\"?'：；‘’“”，。、《》\\]\\[`]")
	//mepNameUrl, _ := url.QueryUnescape(mepName)
	//mepIdUrl, _ := url.QueryUnescape(mepId)
	//endPointUrl, _ := url.QueryUnescape(endPoint)
	if re.MatchString(mepName) || re.MatchString(mepId) || re.MatchString(endPoint) {
		var metas []interface{}
		s.Success(s.NewPage(0, metas, pageNo, pageSize))
		logs.Info("contain speicail chars")
		return
	}
	if val := s.GetString("mepName"); len(val) != 0 {
		mepName = fmt.Sprintf("%%%s%%", val)
	}
	if val := s.GetString("mepId"); len(val) != 0 {
		mepId = fmt.Sprintf("%%%s%%", val)
	}
	if val := s.GetString("endPoint"); len(val) != 0 {
		endPoint = fmt.Sprintf("%%%s%%", val)
	}
	meps, total, err := manager.DefaultMepMg.List(pageNo, pageSize, mepGroupId, notMepGroupId, mepName, mepId, endPoint, s.AuthInfo.AuthMepInfo.MepIds, bean)
	if err != nil {
		logs.Error("list mep meta fail %s", err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(s.NewPage(total, meps, pageNo, pageSize))
}

// @Title Delete mep node
// @Description delete mep node
// @Param       mepId           path    string   true   "mep id"
// @Success 204 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /:mepId [delete]
func (s *NodeController) Delete() {
	logs.Debug("delete mep")
	mepId := s.Ctx.Input.Param(":mepId")

	if err := manager.DefaultMepMg.Delete(mepId); err != nil {
		logs.Error("delete mep meta for %s fail %s", mepId, err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(nil)
}

// @Title Describe mep node
// @Description Get mep node information
// @Param       mepId           path    string   true   "mep id"
// @Success 200 mepmd.MepMeta
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /detail/:mepId [get]
func (s *NodeController) Get() {
	mepId := s.Ctx.Input.Param(":mepId")
	data, err := manager.DefaultMepMg.Get(mepId)
	if err != nil {
		logs.Error("Get mep meta for %s fail %s", mepId, err.Error())
		s.AbortInternalServerError(err)
	}
	s.Success(data)
}

// @Title Update mep status
// @Description update mep manager status
// @Param       mepId           path    string   true   "mep id"
// @Param       status          body    mepmd.Status  true     "status information"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /:mepId/status [put]
func (s *NodeController) UpdateStatus() {
	mepId := s.Ctx.Input.Param(":mepId")
	v := mepmd.MepStatus{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("parse body %s error. %v", string(s.Ctx.Input.RequestBody), err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}

	if err := manager.DefaultMepMg.CheckUpdateStatusReq(v); err != nil {
		s.AbortBadRequest(err)
	}
	err := manager.DefaultMepMg.UpdateStatus(mepId, v.MgrStatus)
	if err != nil {
		logs.Error("UpdateStatus for mepId %s fail %s", mepId, err.Error())
		s.AbortInternalServerError(err)
	}
	s.Success(v)
}

//
//func FormatHttpsFile(s *NodeController, mepConf *mepmd.MepMeta) {
//	//获取根证书
//	if mepConf.HttpsSslType != util.HttpsAuthNone {
//		mepRootCrt, _, err := s.GetFile("mepRootCrt")
//		if err != nil {
//			logs.Error("failed to get mep root crt file: %v", err)
//			s.AbortInternalServerError(common.ErrParaInvalid)
//		}
//
//		mepConf.RootCrt, err = ioutil.ReadAll(mepRootCrt)
//		if err != nil {
//			logs.Error("failed to get mepm client crt file: %v", err)
//			s.AbortInternalServerError(common.ErrParaInvalid)
//		}
//
//		if mepConf.HttpsSslType == util.HttpsAuthBoth {
//			mepmClientCrt, _, err := s.GetFile("mepmClientCrt")
//			if err != nil {
//				logs.Error("failed to get mepm client crt file: %v", err)
//				s.AbortInternalServerError(common.ErrParaInvalid)
//			}
//
//			mepConf.MepmClientCrt, err = ioutil.ReadAll(mepmClientCrt)
//			if err != nil {
//				logs.Error("failed to get mepm client crt file: %v", err)
//				s.AbortInternalServerError(common.ErrParaInvalid)
//			}
//
//			mepmClientKey, _, err := s.GetFile("mepmClientKey")
//			if err != nil {
//				logs.Error("failed to get mepm client key file: %v", err)
//				s.AbortInternalServerError(common.ErrParaInvalid)
//			}
//
//			mepConf.MepmClientCrt, err = ioutil.ReadAll(mepmClientKey)
//			if err != nil {
//				logs.Error("failed to get mepm client crt file: %v", err)
//				s.AbortInternalServerError(common.ErrParaInvalid)
//			}
//		}
//	}
//}
