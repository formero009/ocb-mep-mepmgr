/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package controllers

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/manager"
	"mepmgr/models/mepmd"
	"strconv"
)

type MepTopologyController struct {
	base.BaseController
	MepGTopologyMg manager.MepTopologyManager
}

// @Title Get mep topology
// @Description Get mep topology
// @Param	CALLERId   query string	false
// @Param	MEPMID     query string	false
// @Param   MepOrGroup query string false
// @Success 200 {object} mepmd.RespTopology
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router grouptopology/ [get]
func (s *MepTopologyController) Get() {
	topoReq := mepmd.TopologyReq{}
	topoReq.CALLERId = s.Ctx.Input.Query("CALLERId")
	topoReq.MEPMID = s.Ctx.Input.Query("MEPMID")
	topoReq.MepOrGroup = s.Ctx.Input.Query("MepOrGroup")
	logs.Info("topology: get topology: %v", topoReq.MepOrGroup)
	//get data
	var resq_topologys []mepmd.RespTopology
	if len(topoReq.MepOrGroup) != 0 {
		if err := s.MepGTopologyMg.SearchTopology(topoReq.MepOrGroup, s.AuthInfo.AuthMepInfo.MepIds, &resq_topologys); err != nil {
			logs.Error("Search topology error, err %v", err.Error())
			s.AbortInternalServerError(err)
			return
		}
	} else {
		if err := s.MepGTopologyMg.GetTopology(s.AuthInfo.AuthMepInfo.MepIds, &resq_topologys); err != nil {
			logs.Error("Get topology error, err %v", err.Error())
			s.AbortInternalServerError(err)
			return
		}
	}
	//log.Info("response: ",resq_topologys)
	s.Success(resq_topologys)
}

// @Title Get mep topology
// @Description Get mep topology
// @Param	CALLERId   query string	false
// @Param	MEPMID     query string	false
// @Param   Province query string false
// @Param   City       query string false
// @Success 200 {object} mepmd.RespTopology
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router locationtopology/ [get]
func (s *MepTopologyController) GetByLocation() {
	topoReq := mepmd.TopologyReq{}
	topoReq.CALLERId = s.Ctx.Input.Query("CALLERId")
	topoReq.MEPMID = s.Ctx.Input.Query("MEPMID")
	topoReq.Province = s.Ctx.Input.Query("Province")
	topoReq.City = s.Ctx.Input.Query("City")
	logs.Info("topology: get topology: %v", topoReq.MepOrGroup)
	//get data
	var resq_topologys []mepmd.RespTopologyByProvincial
	if (len(topoReq.Province) != 0 || len(topoReq.City) != 0){
		if err := s.MepGTopologyMg.GetTopologyByLocation(topoReq.Province, topoReq.City, s.AuthInfo.AuthMepInfo.MepIds, &resq_topologys);
			err != nil {
			logs.Error("Search topology error, err %v", err.Error())
			s.AbortInternalServerError(err)
			return
		}
	} else {
		if err := s.MepGTopologyMg.GetTopologyByDefaultLocation(s.AuthInfo.AuthMepInfo.MepIds, &resq_topologys); err != nil {
			logs.Error("Get topology error, err %v", err.Error())
			s.AbortInternalServerError(err)
			return
		}
	}
	//log.Info("response: ",resq_topologys)
	s.Success(resq_topologys)
}


// @Title Get mep topology
// @Description Get mep topology
// @Param   MepId      query string false
// @Param   MepGroupId query string false
// @Success 200 {array} mepmd.RespMepDetail
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router grouptopology/detail [get]
func (s *MepTopologyController) GetDetail() {
	topoReq := mepmd.TopologyReq{}
	topoReq.MepId = s.Ctx.Input.Query("mepId")
	topoReq.MepGroupId = s.Ctx.Input.Query("mepGroupId")
	//valid check
	if len(topoReq.MepId) == 0 && len(topoReq.MepGroupId) == 0 ||
		len(topoReq.MepId) != 0 && len(topoReq.MepGroupId) != 0 {
		logs.Error("topology: only one of mepId and mepGroupId must be selected %+v", topoReq)
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid, fmt.Sprintf("must and only select one of mepId and mepGroupId")))
		return
	}
	logs.Info("topology: get topology of mep: %v, mepGroup: %v", topoReq.MepId, topoReq.MepGroupId)
	//get data
	resq_details := []mepmd.RespMepDetail{}
	if len(topoReq.MepId) != 0 {
		if err := s.MepGTopologyMg.GetMepDetailByMepId(topoReq.MepId, s.MepIds, &resq_details); err != nil {
			s.AbortInternalServerError(err)
		}
	}
	if len(topoReq.MepGroupId) != 0 {
		mepGroupId, err := strconv.ParseInt(topoReq.MepGroupId, 10, 64)
		if err != nil {
			logs.Error("topology: mepGroupId %s ParseInt fail %s", topoReq.MepGroupId, err.Error())
			s.AbortBadRequest(common.NewError(common.ErrParaInvalid))
			return
		}
		if err := s.MepGTopologyMg.GetMepDetailByGroupId(mepGroupId, s.MepIds, &resq_details); err != nil {
			s.AbortInternalServerError(err)
		}
	}
	//log.Info("response: ",resq_detail)
	s.Success(resq_details)
}

// @Title Get mep topology
// @Description Get mep topology
// @Param   Province query string false
// @Param   City       query string false
// @Success 200 {array} mepmd.RespMepDetail
// @Failure 400 invalid request
// @Failure 500 server internal error
// @router locationtopology/detail [get]
func (s *MepTopologyController) GetDetailByLocation() {
	topoReq := mepmd.TopologyReq{}
	topoReq.Province = s.Ctx.Input.Query("Province")
	topoReq.City = s.Ctx.Input.Query("City")

	//valid check
	if len(topoReq.Province) == 0 && len(topoReq.City) == 0 ||
		len(topoReq.Province) != 0 && len(topoReq.City) != 0 {
		logs.Error("topology: only one of Province and City must be selected %+v", topoReq)
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid, fmt.Sprintf("must and only spefic one of Province and City")))
		return
	}
	logs.Info("topology: get topology of mep: %v, mepGroup: %v", topoReq.Province, topoReq.City)
	//get data
	resq_details := []mepmd.RespMepDetail{}

	if err := s.MepGTopologyMg.GetMepDetailByLocation(topoReq.Province, topoReq.City, s.MepIds, &resq_details); err != nil {
		logs.Error("topology: get detail by [%v/%v] failed",topoReq.Province,topoReq.City)
		s.AbortInternalServerError(err)
	}

	//log.Info("response: ",resq_detail)
	s.Success(resq_details)
}