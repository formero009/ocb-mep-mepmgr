/*
@Time : 2022/1/6
@Author : klp
@Project: mepmgr
*/
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
)

type ConfigController struct {
	base.BaseController
}

// @Title Update config param
// @Description update config param
// @Param       param  body   mepmd.ConfigParam  true   "config param"
// @Success 200 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router / [put]
func (s *ConfigController) Update() {
	logs.Debug("update param")

	config := mepmd.ConfigParameter{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &config); err != nil {
		logs.Error("parse body %s error. %v", string(s.Ctx.Input.RequestBody), err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}
	if err := manager.DefaultConfigMg.CheckReq(config); err != nil {
		s.AbortBadRequest(err)
	}
	if err := manager.DefaultConfigMg.Update(config); err != nil {
		logs.Error("update config param  %+v fail %s", config, err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(config)
}

// @Title List config param
// @Description List config param
// @Param	currentPage	query 	int 	 false	"page number"
// @Param	pageSize	query 	int 	 false	"page size"
// @Param	MEPMID  	query 	string	 false	"MEPMID"
// @Param	CALLERId  	query 	string	 false	"CALLERId"
// @Success 200 []mepmd.ConifgParam
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router  / [get]
func (s *ConfigController) List() {
	logs.Debug("list config param")

	pageSize, _ := s.GetInt(common.PageSize, common.DefaultPageSize)
	pageNo, _ := s.GetInt(common.CurrentPage, common.DefaultPageNumber)

	configParams, total, err := manager.DefaultConfigMg.List(pageNo, pageSize)
	if err != nil {
		logs.Error("list mep meta fail %s", err.Error())
		s.AbortInternalServerError(err)
	}

	s.Success(s.NewPage(total, configParams, pageNo, pageSize))
}
