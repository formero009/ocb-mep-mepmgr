package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
)

type PerferenceController struct {
	base.BaseController
	PerferenceMg manager.PerferenceManager
}

// @Title save perference
// @Description save perference
// @Param   Preference body mepmd.Preference true
// @Success 200 {object} mepmd.Preference
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router /perference [post]
func (s *PerferenceController) Post() {
	userName := s.Ctx.Input.Header("userName")
	preference := &mepmd.Preference{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &preference); err != nil {
		logs.Error("parse body %s error. %v", string(s.Ctx.Input.RequestBody), err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}
	preference.UserName = userName
	if err := s.PerferenceMg.SavePerference(preference); err != nil {
		s.AbortBadRequest(err)
	}
	s.Success(preference)
}

// @Title get perference
// @Description get perference
// @Success 200 {object} mepmd.Preference
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router /perference [get]
func (s *PerferenceController) Get() {
	userName := s.Ctx.Input.Header("userName")
	preference, err := s.PerferenceMg.GetPerference(userName)
	if err != nil {
		s.AbortBadRequest(err)
	}
	s.Success(preference)
}
