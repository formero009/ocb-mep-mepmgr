/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
	"mepmgr/util"
	"reflect"
	"strings"
)

type MepGroupController struct {
	base.BaseController
	MepGroupMg manager.MepGroupManager
}

// @Title create mep group
// @Description create mep group information
// @Param	mepGroup  body	models.mepmd.MepGroup  true  "mep group information"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup [post]
func (s *MepGroupController) Create() {
	var v mepmd.MepGroup
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("parse body error. %v", err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}
	valid, msg := s.createValidCheck(v)
	if !valid {
		logs.Error(msg)
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid, msg))
		return
	}

	if err := s.MepGroupMg.Create(&v); err != nil {
		logs.Error("create mep group: error. %v", err)
		s.AbortInternalServerError(err)
		return
	}

	s.Success(v)
}
func (s *MepGroupController) createValidCheck(request mepmd.MepGroup) (bool, string) {
	if util.IsBlank(reflect.ValueOf(request.CALLERId)) {
		logs.Error("field CALLERId is empty")
		return false, fmt.Sprintf("%s不能为空", "CALLERId")
	}
	if util.IsBlank(reflect.ValueOf(request.MEPMID)) {
		logs.Error("field MEPMID is empty")
		return false, fmt.Sprintf("%s不能为空", "MEPMID")
	}
	if util.IsBlank(reflect.ValueOf(request.MepGroupName)) {
		logs.Error("field mepGroupName is empty")
		return false, fmt.Sprintf("%s不能为空", "分组名称")
	}
	return true, ""
}

// @Title Delete mep group
// @Description delete mep group
// @Param	mepGroupName  path 	string	true  "mep group name"
// @Success 204 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup/:mepGroupName [delete]
func (s *MepGroupController) Delete() {
	name := s.Ctx.Input.Param(":mepGroupName")

	if err := s.MepGroupMg.Delete(name); err != nil {
		logs.Error("delete trafficPolicy: error. %v", err)
		s.AbortInternalServerError(err.Error())
		return
	}
	s.Success(nil)
}

// @Title Update mep group
// @Description update mep group
// @Param	mepGroupName  path 	string	true  "mep group name"
// @Param	mepGroupUpdate  body	models.mepmd.MepGroupUpdate  true  "mep group information"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup/:mepGroupName [put]
func (s *MepGroupController) Update() {
	name := s.Ctx.Input.Param(":mepGroupName")
	var v mepmd.MepGroupUpdate
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("parse body error. %v", err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}

	valid, msg := s.putValidCheck(v)
	if !valid {
		logs.Error(msg)
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid, msg))
		return
	}

	if err := s.MepGroupMg.Update(name, v); err != nil {
		logs.Error("update error. %v", err)
		s.AbortInternalServerError(err)
		return
	}
	s.Success(v)
}

func (s *MepGroupController) putValidCheck(request mepmd.MepGroupUpdate) (bool, string) {
	if util.IsBlank(reflect.ValueOf(request.CALLERId)) {
		logs.Error("field CALLERId is empty")
		return false, fmt.Sprintf("%s不能为空", "CALLERId")
	}
	if util.IsBlank(reflect.ValueOf(request.MEPMID)) {
		logs.Error("field MEPMID is empty")
		return false, fmt.Sprintf("%s不能为空", "MEPMID")
	}
	if util.IsBlank(reflect.ValueOf(request.MepGroupName)) {
		logs.Error("field mepGroupName is empty")
		return false, fmt.Sprintf("%s不能为空", "分组名称")
	}
	return true, ""
}

// @Title List mep group
// @Description list mep group
// @Param	mepGroupName	query 	string	 false	"mep group name"
// @Param	MEPMID	        query 	string	 false	"mepm id"
// @Param	CALLERId	    query 	string	 false	"caller id"
// @Param	pageSize	    query 	int 	 false	"page size"
// @Param	currentPage  	    query 	int 	 false	"page number"
// @Success 200 {object} []models.mepmd.MepGroup
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup [get]
func (s *MepGroupController) List() {
	name := s.GetString("mepGroupName")
	mepmId := s.GetString("MEPMID")
	callerId := s.GetString("CALLERId")

	size, _ := s.GetInt("pageSize", common.DefaultPageSize)
	index, _ := s.GetInt("currentPage", common.DefaultPageNumber)

	if size > common.MaxPageSize {
		size = common.MaxPageSize
	}

	if strings.ContainsAny(name, "~!@#$%^&*(){}|<>\\\\/+\\-=【】:\"?'：；‘’“”，。、《》\\]\\[`") {
		logs.Error("invalid search character")
		var metas []interface{}
		s.Success(s.NewPage(0, metas, index, size))
		return
	}

	match := ""
	if len(name) != 0 {
		match = fmt.Sprintf("%%%s%%", name)
	}
	bean := mepmd.MepGroup{MEPMID: mepmId, CALLERId: callerId}
	logs.Debug(name)
	meps, count, err := s.MepGroupMg.List(index, size, match, bean)
	if err != nil {
		logs.Error("list mep group meta fail %s", err)
		s.AbortInternalServerError(err)
	}

	s.Success(s.NewPage(count, meps, index, size))
}

// @Title Add mep to mep group
// @Description add mep to mep group
// @Param	mepGroupName  path 	string	true  "mep group name"
// @Param	addMep  body	models.mepmd.MepGroupAddMep  true  "mepNames"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup/:mepGroupName/mep [post]
func (s *MepGroupController) AddMep() {
	name := s.Ctx.Input.Param(":mepGroupName")
	var v mepmd.MepGroupAddMep
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("parse body error. %v", err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}

	valid, msg := s.addMepValidCheck(v)
	if !valid {
		logs.Error(msg)
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid, msg))
		return
	}

	if err := s.MepGroupMg.AddMep(name, v); err != nil {
		logs.Error("add mep to group error. %v", err)
		s.AbortInternalServerError(err)
		return
	}
	s.Success(v)
}

func (s *MepGroupController) addMepValidCheck(request mepmd.MepGroupAddMep) (bool, string) {
	if len(request.MepNames) == 0 {
		logs.Error("field MepNames is empty")
		return false, fmt.Sprintf("%s不能为空", "Mep名称")
	}
	return true, ""
}

// @Title Delete mep from mep group
// @Description delete mep from mep group
// @Param	mepGroupName  path 	string	true  "mep group name"
// @Param	mepName       path 	string	true  "mep name"
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 mep group parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /mepGroup/:mepGroupName/mep/:mepName [delete]
func (s *MepGroupController) DeleteMep() {
	mepGroupName := s.Ctx.Input.Param(":mepGroupName")
	mepName := s.Ctx.Input.Param(":mepName")

	if err := s.MepGroupMg.DelMep(mepGroupName, mepName); err != nil {
		logs.Error("delete mep from group error. %v", err)
		s.AbortInternalServerError(err)
		return
	}
	s.Success(nil)
}
