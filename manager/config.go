/*
@Time : 2022/1/6
@Author : klp
@Project: mepmgr
*/
package manager

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/dao"
	"mepmgr/models/mepmd"
	"strconv"
)

var DefaultConfigMg configMg

type configMg struct{}

var (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

func InitConfig() {
	params := make([]mepmd.ConfigParameter, 0, 10)
	if err := dao.ConfigParamDao.List(0, 0, "", "", mepmd.ConfigParameter{}, &params); err != nil {
		panic(fmt.Sprintf("List Config Param fail %s", err.Error()))
	}
	name2Parm := make(map[string]map[string]mepmd.ConfigParameter)
	for _, param := range params {
		if _, f := name2Parm[param.ParaGroup]; !f {
			name2Parm[param.ParaGroup] = make(map[string]mepmd.ConfigParameter)
		}
		name2Parm[param.ParaGroup][param.ParaName] = param
	}
	if _, f := name2Parm[mepmd.MaintainConfig][mepmd.HeartBeatInterval]; !f {
		// insert db
		para := mepmd.ConfigParameter{ParaGroup: mepmd.MaintainConfig, ParaName: mepmd.HeartBeatInterval, ParaValue: strconv.Itoa(DefaultHeartBeatInterval), Status: StatusSuccess}
		if err := dao.ConfigParamDao.Create(&para); err != nil {
			panic(fmt.Sprintf("Create Config Param %+v fail %s", para, err.Error()))
		}
	}
	logs.Info("config init success")
	return
}

func (m configMg) Create(config *mepmd.ConfigParameter) error {
	if err := dao.ConfigParamDao.Create(config); err != nil {
		logs.Error("failed to create mepm info: %v", err)
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m configMg) Update(config mepmd.ConfigParameter) error {
	bean := mepmd.ConfigParameter{ParaGroup: config.ParaGroup, ParaName: config.ParaName}
	if err := dao.ConfigParamDao.Get(&bean); err != nil {
		logs.Error("failed to get configParam %+v err : %v", bean, err)
		return common.NewError(common.ErrDatabase)
	}
	if err := dao.ConfigParamDao.Update(bean, config); err != nil {
		logs.Error("failed to update configParam bean %+v to %+v err %v", bean, config, err)
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

var validParamName = map[string]bool{
	mepmd.HeartBeatInterval: true,
}

func (m configMg) CheckReq(config mepmd.ConfigParameter) error {
	if !validParamName[config.ParaName] {
		return common.NewError(common.ErrParaInvalid, "invalid paramName")
	}
	if config.ParaName == mepmd.HeartBeatInterval {
		if val, err := strconv.Atoi(config.ParaValue); err != nil {
			logs.Error("paramVal %s atoi fail %s", config.ParaValue, err.Error())
			return common.NewError(common.ErrParaInvalid, "invalid paramValue")
		} else if val > MaxHeartBeatInterval || val < MinHeartBeatInterval {
			return common.NewError(common.ErrParaInvalid, "paramValue out of range")
		}
	}
	return nil
}

func (m configMg) List(from, size int) ([]mepmd.ConfigParameter, int64, error) {
	params := make([]mepmd.ConfigParameter, 0, 10)
	var total int64 = -1
	err := dao.ConfigParamDao.List(from, size, "", "", mepmd.ConfigParameter{}, &params)
	if err != nil {
		logs.Error("List config param fail %s", err.Error())
		return params, total, common.NewError(common.ErrDatabase)
	}
	err = dao.ConfigParamDao.Count(mepmd.ConfigParameter{}, &total)
	if err != nil {
		logs.Error("count config param fail %s", err.Error())
		return params, total, common.NewError(common.ErrDatabase)
	}
	return params, total, nil
}
