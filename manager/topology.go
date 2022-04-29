/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package manager

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/dao"
	"mepmgr/models/mepmd"
	"mepmgr/util"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MepTopologyManager interface {
	GetTopology(authMepIds []string, topologys *[]mepmd.RespTopology) error
	GetTopologyByLocation(provincial, city string, authMepIds []string, topologys *[]mepmd.RespTopologyByProvincial) error
	GetTopologyByDefaultLocation(authMepIds []string, meps *[]mepmd.RespTopologyByProvincial) error
	SearchTopology(key string, authMepIds []string, topologys *[]mepmd.RespTopology) error
	GetMepDetailByGroupId(group_id int64, authMepIds []string, data *[]mepmd.RespMepDetail) error
	GetMepDetailByMepId(mep_id string, authMepIds []string, data *[]mepmd.RespMepDetail) error
	GetMepDetailByLocation(provincial, city string, authMepIds []string, data *[]mepmd.RespMepDetail) error
}

var mepTopologyOnce sync.Once
var mepTopologyMg defaultMepTopologyMg

type NetworkNames struct {
	NetworkNames []string `json:"networkNames"`
}

func NewDefaultMepTopologyManager() MepTopologyManager {
	mepTopologyOnce.Do(func() {
		mepTopologyMg = defaultMepTopologyMg{}
	})
	return mepTopologyMg
}

type defaultMepTopologyMg struct{}

func (m defaultMepTopologyMg) GetTopology(authMepIds []string, topologys *[]mepmd.RespTopology) error {
	err, data1 := dao.TopologyDao.GetMepGroupRelation(authMepIds)
	if err != nil {
		logs.Error("get mep topologys err.%v", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	err, data2 := dao.TopologyDao.GetMepNoRelationGroup(authMepIds)
	if err != nil {
		logs.Error("get mep topologys err.%v", err.Error())
		return common.NewError(common.ErrDatabase)
	}

	err, data3 := dao.TopologyDao.GetGroupNoRelationMep()
	if err != nil {
		logs.Error("get mep  topologys err.%v", err.Error())
		return common.NewError(common.ErrDatabase)
	}

	m_data := make(map[string][]mepmd.TopologyMep)
	for _, v := range *data1 {
		if v.Id != 0 && v.MepGroupName != "" {
			key := strconv.FormatInt(v.Id, 10) + ":" + v.MepGroupName
			mep := mepmd.TopologyMep{
				MepId:   v.MepId,
				MepName: v.MepName,
			}
			m_data[key] = append(m_data[key], mep)
		}
	}

	var resq_topologys []mepmd.RespTopology
	for k, v := range m_data {
		strs := strings.Split(k, ":")
		mep_id, _ := strconv.ParseInt(strs[0], 10, 64)
		resq := mepmd.RespTopology{
			MepGroupName: strs[1],
			MepGroupId:   mep_id,
			Meps:         v,
		}
		resq_topologys = append(resq_topologys, resq)
	}
	resq := mepmd.RespTopology{
		Meps: *data2,
	}
	resq_topologys = append(resq_topologys, resq)

	for _, v := range *data3 {
		if v.Id != 0 && v.MepGroupName != "" {
			resq := mepmd.RespTopology{
				MepGroupName: v.MepGroupName,
				MepGroupId:   v.Id,
				Meps:         []mepmd.TopologyMep{},
			}
			resq_topologys = append(resq_topologys, resq)
		}
	}

	*topologys = resq_topologys
	logs.Info("get mep topologys data:", resq_topologys)
	return nil
}

func (m defaultMepTopologyMg) SearchTopology(key string, authMepIds []string, topologys *[]mepmd.RespTopology) error {
	err, data1 := dao.TopologyDao.SearchMepGroupbykey(key, authMepIds)
	if err != nil {
		logs.Error("search mep group topologys err.%v", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	err, data2 := dao.TopologyDao.SearchMepbykey(key, authMepIds)
	if err != nil {
		logs.Error("search mep topologys err.%v", err.Error())
		return common.NewError(common.ErrDatabase)
	}

	//merge
	data := *data1
	data = append(data, *data2...)
	/*
		for _, v := range *data2 {
			data = append(data, v)
		}
	*/
	//Remove Repeated
	resp := RemoveRepeatedElement(data)

	m_data := make(map[string][]mepmd.TopologyMep)
	for _, v := range resp {
		if v.Id != 0 && v.MepGroupName != "" {
			key := strconv.FormatInt(v.Id, 10) + ":" + v.MepGroupName
			mep := mepmd.TopologyMep{
				MepId:   v.MepId,
				MepName: v.MepName,
			}
			m_data[key] = append(m_data[key], mep)
		} else {
			key := "-"
			mep := mepmd.TopologyMep{
				MepId:   v.MepId,
				MepName: v.MepName,
			}
			m_data[key] = append(m_data[key], mep)
		}
	}

	var resq_topologys []mepmd.RespTopology
	for k, v := range m_data {
		if k != "-" {
			strs := strings.Split(k, ":")
			mep_id, _ := strconv.ParseInt(strs[0], 10, 64)
			resq := mepmd.RespTopology{
				MepGroupName: strs[1],
				MepGroupId:   mep_id,
				Meps:         v,
			}
			resq_topologys = append(resq_topologys, resq)
		} else {
			resq := mepmd.RespTopology{
				Meps: v,
			}
			resq_topologys = append(resq_topologys, resq)
		}
	}
	*topologys = resq_topologys
	logs.Info("search mep topologys data:", topologys)
	return nil
}

func RemoveRepeatedElement(arr []mepmd.TopologyMepGroupRe) (newArr []mepmd.TopologyMepGroupRe) {
	newArr = make([]mepmd.TopologyMepGroupRe, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i].Id == arr[j].Id &&
				arr[i].MepGroupName == arr[j].MepGroupName &&
				arr[i].MepId == arr[j].MepId &&
				arr[i].MepName == arr[j].MepName {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func (m defaultMepTopologyMg) GetMepDetailByGroupId(group_id int64, authMepIds []string, meps *[]mepmd.RespMepDetail) error {
	err, mep_list := dao.TopologyDao.GetMepByGroupId(group_id, authMepIds)
	if err != nil {
		logs.Error("topology: get meps by group id[%d]  err.%v", group_id, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	var mepAlarmMap map[string]mepmd.AlarmResp
	if len(*mep_list) != 0 {
		mepAlarmMap = getAlarmsByMepName(mep_list)
	} else {
		logs.Error("topology: no group id is [%d]", group_id)
	}

	meps_deatil := []mepmd.RespMepDetail{}
	for _, v := range *mep_list {
		alarm, ok := mepAlarmMap[v.MepName]
		if !ok {
			alarm = mepmd.AlarmResp{}
		}
		mep_detail := mepmd.RespMepDetail{
			MepId:      v.MepId,
			MepName:    v.MepName,
			EndPoint:   v.EndPoint,
			User:       v.User,
			UserTag:    v.UserTag,
			Type:       v.Type,
			Province:   v.Province,
			City:       v.City,
			Latitude:   v.Latitude,
			Longitude:  v.Longitude,
			RunStatus:  v.RunStatus,
			MgrStatus:  v.MgrStatus,
			SwVersion:  v.SwVersion,
			Contractor: v.Contractor,
			AlarmMsg:   alarm,
		}
		meps_deatil = append(meps_deatil, mep_detail)
	}
	*meps = meps_deatil
	return nil
}

func getAlarmsByMepName(mep_list *[]mepmd.MepMeta) map[string]mepmd.AlarmResp {
	var mepNames NetworkNames
	for _, v := range *mep_list {
		mepNames.NetworkNames = append(mepNames.NetworkNames, v.MepName)
	}
	mepAlarmMap := map[string]mepmd.AlarmResp{}

	addr := beego.AppConfig.String("AlarmAddress")
	url := addr + common.AlarmStatisticsPath
	reqBody, err := json.Marshal(mepNames)
	if err != nil {
		logs.Error("topology: Marshal body error. %v", err)
	}
	headers := map[string]string{"Content-Type": "application/json"}
	code, body, _, err := util.DoRequest("GET", url, "", headers, reqBody, httpTimeoutMs, &util.HttpsConf{HttpsSslType: util.HttpsAuthNone})
	if err != nil || code != http.StatusOK {
		logs.Error("topology: DoRequest fail %v %d %v", err, code, string(body))
		return mepAlarmMap
	}
	var v mepmd.AlarmStatistics
	if err := json.Unmarshal(body, &v); err != nil {
		logs.Error("topology: parse body from alarm statistics error. %v", err)
		return mepAlarmMap
	}
	logs.Info("alarm:", v)
	if v.Code != "opensigma.common.common.Success" || v.Message != "成功" {
		logs.Error("topology: get alarm data fail %v %d", v.Message, v.Code)
		return mepAlarmMap
	}
	for _, mep := range v.Data {
		mepAlarmMap[mep.NetworkName] = mep.AlarmResp
	}
	return mepAlarmMap
}

func (m defaultMepTopologyMg) GetMepDetailByMepId(mep_id string, authMepIds []string, meps *[]mepmd.RespMepDetail) error {
	err, mep_list := dao.TopologyDao.GetMepByMepId(mep_id, authMepIds)
	if err != nil {
		logs.Error("topology: get meps by mep id[%v]  err.%v", mep_id, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	var mepAlarmMap map[string]mepmd.AlarmResp
	if len(*mep_list) != 0 {
		mepAlarmMap = getAlarmsByMepName(mep_list)
	} else {
		logs.Error("topology: no mep id is [%v]", mep_id)
	}

	meps_deatil := []mepmd.RespMepDetail{}
	for _, v := range *mep_list {
		alarm, ok := mepAlarmMap[v.MepName]
		if !ok {
			alarm = mepmd.AlarmResp{}
		}
		mep_detail := mepmd.RespMepDetail{
			MepId:      v.MepId,
			MepName:    v.MepName,
			EndPoint:   v.EndPoint,
			User:       v.User,
			UserTag:    v.UserTag,
			Type:       v.Type,
			Province:   v.Province,
			City:       v.City,
			Latitude:   v.Latitude,
			Longitude:  v.Longitude,
			RunStatus:  v.RunStatus,
			MgrStatus:  v.MgrStatus,
			SwVersion:  v.SwVersion,
			Contractor: v.Contractor,
			AlarmMsg:   alarm,
		}
		meps_deatil = append(meps_deatil, mep_detail)
	}
	*meps = meps_deatil
	return nil
}

func (m defaultMepTopologyMg) GetTopologyByLocation(provincial, city string, authMepIds []string, meps *[]mepmd.RespTopologyByProvincial) error {
	if len(provincial) == 0 {
		logs.Error("topology: province is empty ")
		return common.NewError(common.ErrDatabase)
	}
	err, mep_list := dao.TopologyDao.GetMepByLocation(provincial, city, authMepIds)
	if err != nil {
		logs.Error("topology: get meps by provincial [%v] city [%v] err.%v", provincial, city, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	logs.Info(mep_list)
	if len(*mep_list) == 0 {
		logs.Info("topology: get meps by provincial [%v] city [%v] is null", provincial, city)
		return nil
	}
	map_c := make(map[string][]mepmd.TopologyMep)
	for _, mep := range *mep_list {
		t_mep := mepmd.TopologyMep{
			MepId:   mep.MepId,
			MepName: mep.MepName,
		}
		map_c[mep.City] = append(map_c[mep.City], t_mep)
	}
	var resq_ps []mepmd.RespTopologyByProvincial
	var resq_cs []mepmd.RespTopologyByCity
	for k, v := range map_c {
		resq_c := mepmd.RespTopologyByCity{
			City: k,
			Meps: v,
		}
		resq_cs = append(resq_cs, resq_c)
	}
	resq_p := mepmd.RespTopologyByProvincial{
		Province: provincial,
		Cites:    resq_cs,
	}

	resq_ps = append(resq_ps, resq_p)
	*meps = resq_ps
	return err
}

func (m defaultMepTopologyMg) GetTopologyByDefaultLocation(authMepIds []string, meps *[]mepmd.RespTopologyByProvincial) error {
	provincial, city := "", ""
	err, mep_list := dao.TopologyDao.GetMepByLocation(provincial, city, authMepIds)
	if err != nil {
		logs.Error("topology: get meps by provincial [%v] city [%v] err.%v", provincial, city, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	logs.Info(mep_list)
	map_c := make(map[string][]mepmd.TopologyMep)
	map_p := make(map[string][]string)
	for _, mep := range *mep_list {
		t_mep := mepmd.TopologyMep{
			MepId:   mep.MepId,
			MepName: mep.MepName,
		}
		_, ok := map_c[mep.City]
		if !ok {
			map_p[mep.Province] = append(map_p[mep.Province], mep.City)
		}
		map_c[mep.City] = append(map_c[mep.City], t_mep)
	}

	var resq_ps []mepmd.RespTopologyByProvincial
	for k, v := range map_p {
		var resq_cs []mepmd.RespTopologyByCity
		for _, s := range v {
			resq_c := mepmd.RespTopologyByCity{
				City: s,
				Meps: map_c[s],
			}
			resq_cs = append(resq_cs, resq_c)
		}
		resq_p := mepmd.RespTopologyByProvincial{
			Province: k,
			Cites:    resq_cs,
		}
		resq_ps = append(resq_ps, resq_p)
	}

	*meps = resq_ps
	return err
}

func (m defaultMepTopologyMg) GetMepDetailByLocation(provincial, city string, authMepIds []string, meps *[]mepmd.RespMepDetail) error {
	err, mep_list := dao.TopologyDao.GetMepByLocation(provincial, city, authMepIds)
	if err != nil {
		logs.Error("topology: get meps by provincial [%v] city [%v] err.%v", provincial, city, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	var mepAlarmMap map[string]mepmd.AlarmResp
	if len(*mep_list) != 0 {
		mepAlarmMap = getAlarmsByMepName(mep_list)
	} else {
		logs.Error("topology: no mep with [%v/%v]", provincial, city)
	}

	meps_deatil := []mepmd.RespMepDetail{}
	for _, v := range *mep_list {
		alarm, ok := mepAlarmMap[v.MepName]
		if !ok {
			alarm = mepmd.AlarmResp{}
		}
		mep_detail := mepmd.RespMepDetail{
			MepId:      v.MepId,
			MepName:    v.MepName,
			EndPoint:   v.EndPoint,
			User:       v.User,
			UserTag:    v.UserTag,
			Type:       v.Type,
			Province:   v.Province,
			City:       v.City,
			Latitude:   v.Latitude,
			Longitude:  v.Longitude,
			RunStatus:  v.RunStatus,
			MgrStatus:  v.MgrStatus,
			SwVersion:  v.SwVersion,
			Contractor: v.Contractor,
			AlarmMsg:   alarm,
		}
		meps_deatil = append(meps_deatil, mep_detail)
	}
	*meps = meps_deatil
	return nil
}
