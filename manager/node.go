/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package manager

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/guonaihong/gout"
	"mepmgr/common"
	"mepmgr/dao"
	"mepmgr/models/certmd"
	"mepmgr/models/logmd"
	"mepmgr/models/mepmd"
	"mepmgr/util"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type MepManager interface {
	Create(mep *mepmd.MepMeta, userId string) error
	Delete(mepId string) error
	Update(mep *mepmd.MepMeta) error
	List(from, size int, groupId, excludeGroupId int64, matchStr string, bean mepmd.MepMeta) ([]mepmd.MepMeta, int64, error)
	Get(mepId string) (mepmd.MepMeta, error)
	UpdateStatus(mepId, status string) error
	Topology() (map[string][]string, error)
}

var DefaultMepMg mepMg

type mepMg struct {
	NodeChecker *nodeCheckKeeper
}

type nodeCheckKeeper struct {
	timer   *time.Timer
	endChan chan bool
	wg      *sync.WaitGroup
}

var maxCon = 20
var httpTimeoutMs int64 = 2000
var mepSalt = "OpenSigma@10086"
var DefaultHeartBeatInterval = 60
var MinHeartBeatInterval = 10
var MaxHeartBeatInterval = 100
var LastHeartBeatInterval = DefaultHeartBeatInterval
var HeaderMepId = "X-Mepids"

func InitChecker() {
	logs.Info("node init")
	n := nodeCheckKeeper{}
	period, err := getCheckPeriod()
	if err != nil {
		panic(fmt.Sprintf("getCheckPeriod fail %s", err.Error()))
	}
	LastHeartBeatInterval = period
	n.timer = time.NewTimer(time.Second * time.Duration(period))
	n.endChan = make(chan bool, 1)
	n.wg = &sync.WaitGroup{}
	DefaultMepMg.NodeChecker = &n
	mepSalt = beego.AppConfig.DefaultString("mepSalt", mepSalt)
	logs.Info("node init slat %s checker %d", mepSalt, period)
}

func getCheckPeriod() (int, error) {
	bean := mepmd.ConfigParameter{ParaGroup: mepmd.MaintainConfig, ParaName: mepmd.HeartBeatInterval}
	err := dao.ConfigParamDao.Get(&bean)
	if err != nil {
		return -1, err
	}
	//period := beego.AppConfig.DefaultInt("MepCheckerDuration", 60)
	checkPeriod, err := strconv.Atoi(bean.ParaValue)
	if err != nil {
		return -1, err
	}
	return checkPeriod, nil
}

func SetAuthForMepId(method, path, userId, mepId string) error {
	addr := beego.AppConfig.String("AuthAddress")
	reqUrl := fmt.Sprintf("%s%s", addr, path)
	var reqBody []byte
	var err error
	if method == "POST" {
		req := struct {
			UserId string `json:"userId"`
			MepId  string `json:"mepId"`
		}{UserId: userId, MepId: mepId}
		reqBody, err = json.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshal req %+v fail %s", req, err.Error())
		}
	}

	//向auth发起请求
	code, body, _, err := util.DoRequest(method, reqUrl, "", nil, reqBody, httpTimeoutMs, &util.HttpsConf{HttpsSslType: util.HttpsAuthNone})
	if err != nil || code != http.StatusOK {
		logs.Error("DoRequest url [%s] fail %v %d %s", reqUrl, err, code, string(body))
		return common.NewError(common.ErrNetwork)
	}
	resp := common.ErrMsg{}
	if err = json.Unmarshal(body, &resp); err != nil {
		logs.Error("unmarshl auth body %s fail %s", string(body), err.Error())
		return common.NewError(common.ErrAuthResp)
	}
	if resp.Code != common.ErrSuccess {
		logs.Error("auth resp code error %s", resp.Code)
		return common.NewError(common.ErrAuthResp, resp.Code)
	}
	logs.Info("SetAuthForMepId method: %s url: %s userId: %s mepId: %s success", method, reqUrl, userId, mepId)

	return nil
}

func (m *mepMg) CheckReq(mep *mepmd.MepMeta) error {
	urlP, err := url.Parse(mep.EndPoint)
	if err != nil {
		logs.Error("parse endPoint %s fail %s", mep.EndPoint, err.Error())
		return common.NewError(common.ErrParaInvalid, fmt.Sprintf("endPoint %s %s", mep.EndPoint, common.MsgInvalidFormat))
	}
	if urlP.Scheme == "" || urlP.Host == "" {
		return common.NewError(common.ErrParaInvalid, "EndPoint中缺少协议或host信息")
	}

	if urlP.Scheme == "https" {
		//目前必须仅支持双向认证
		if mep.HttpsSslType != util.HttpsAuthBoth {
			return common.NewError(common.ErrParaInvalid, "Https协议需要双向认证")
		}
		if mep.RootCrt == "" {
			return common.NewError(common.ErrParaInvalid, "Https证书文件不能为空")
		}

	} else if urlP.Scheme == "http" {
		if mep.HttpsSslType != util.HttpsAuthNone {
			return common.NewError(common.ErrParaInvalid, "Http协议不需要认证")
		}
	} else {
		return common.NewError(common.ErrParaInvalid, "Endpoint协议错误")
	}

	return nil
}

func (m *mepMg) Create(mep *mepmd.MepMeta, userId string) error {
	httpsConf := util.HttpsConf{
		HttpsSslType: mep.HttpsSslType,
		HttpsRootCrt: mep.RootCrt,
	}
	bean := &mepmd.MepMeta{MepName: mep.MepName}
	if err := dao.MepDao.Get(bean); err == nil && bean.Id != 0 {
		return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("mepName %s %s", mep.MepName, common.MsgAlreadyExist))
	}
	bean = &mepmd.MepMeta{EndPoint: mep.EndPoint}
	if err := dao.MepDao.Get(bean); err == nil && bean.Id != 0 {
		return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("endPoint %s %s", mep.EndPoint, common.MsgAlreadyExist))
	}

	// 鉴权
	mep.PassWord = util.SaltEncode(mep.PassWord, mepSalt)
	auth, err := mepAuthentication(mep.EndPoint, mep.User, mep.PassWord, &httpsConf)
	if err == nil {
		mep.Token = auth.Token
		mep.RunStatus = mepmd.RunStatusOn

		bean = &mepmd.MepMeta{MepId: auth.MepId}
		if err := dao.MepDao.Get(bean); err == nil && bean.Id != 0 {
			return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("mepId %s %s", auth.MepId, common.MsgAlreadyExist))
		}
		mep.MepId = auth.MepId
	} else {
		logs.Error("mepAuthentication fail %s", err.Error())
		return common.NewError(common.ErrAuthFailed, fmt.Sprintf("mep.EndPoint %s %s", mep.EndPoint, err.Error()))
	}

	mep.SwVersion = mepmd.DefaultVersion
	if mepInfo, err := getMepInfo(mep.EndPoint, auth.Token, &httpsConf); err == nil {
		mep.SwVersion = mepInfo.MepVersion
	} else {
		logs.Error("getMepInfo %s with token %s fail %s", mep.EndPoint, auth.Token, err.Error())
	}

	// 设置权限
	if err := SetAuthForMepId("POST", common.AuthMepAuthorityPath, userId, mep.MepId); err != nil {
		logs.Error("SetAuthForMepId POST userId %s mepId %s fail %s", userId, mep.MepId, err.Error())
		return err
	}

	// 添加数据库
	err = dao.MepDao.Create(mep)
	if err != nil {
		logs.Error("crate mep [%+v] fail %s", mep, err.Error())
		return common.NewError(common.ErrDatabase)
	}

	//通知alarm
	var code int
	addr := beego.AppConfig.String("AlarmAddress")
	err = gout.POST(addr + "/mepm/alarm/v1/mepmgr/notify/add").
		SetQuery(map[string]string{"mepId": mep.MepId}).
		Code(&code).
		Do()
	if err != nil || code != http.StatusOK {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		logs.Error("notify mep %s add to alarm fail,code:%d,err:%s", mep.MepId, code, errMsg)
		return common.NewError(common.ErrInternal)
	}

	return nil
}

func (m *mepMg) Update(mep *mepmd.MepMeta) error {
	httpsConf := util.HttpsConf{
		HttpsSslType: mep.HttpsSslType,
		HttpsRootCrt: mep.RootCrt,
	}

	mepMdl := mepmd.MepMeta{MepName: mep.MepName}
	if err := dao.MepDao.Get(&mepMdl); err == nil && mepMdl.MepId != mep.MepId {
		return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("mepName %s %s", mep.MepName, common.MsgAlreadyExist))
	}
	mepMdl = mepmd.MepMeta{EndPoint: mep.EndPoint}
	if err := dao.MepDao.Get(&mepMdl); err == nil && mepMdl.MepId != mep.MepId {
		return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("endPoint %s %s", mep.EndPoint, common.MsgAlreadyExist))
	}
	mepMdl = mepmd.MepMeta{MepId: mep.MepId}
	err := dao.MepDao.Get(&mepMdl)
	if err != nil {
		logs.Error("update mep, Get mep meta for %s fail %s", mep.MepId, err.Error())
		if err == dao.NOTFOUNDGET {
			return common.NewError(common.ErrNotFound, fmt.Sprintf("mepId %s %s", mep.MepId, common.MsgNotFound))
		}
		return common.NewError(common.ErrDatabase)
	}

	if mepMdl.MgrStatus == mepmd.MgrStatusLock || mepMdl.MgrStatus == mepmd.MgrStatusClose {
		logs.Error("couldn't update mep %s under status %s", mepMdl.MepId, mepMdl.MgrStatus)
		return common.NewError(common.ErrMepStatus, "当前状态的mep不能更新")
	}

	encodedPassWord := util.SaltEncode(mep.PassWord, mepSalt)
	if mep.PassWord != mepMdl.PassWord {
		mep.PassWord = encodedPassWord
	}
	mep.Token = mepMdl.Token
	mep.RunStatus = mepMdl.RunStatus

	if mep.User != mepMdl.User || mep.PassWord != mepMdl.PassWord {
		if auth, err := mepAuthentication(mepMdl.EndPoint, mep.User, mep.PassWord, &httpsConf); err == nil {
			mep.Token = auth.Token
			mep.RunStatus = mepmd.RunStatusOn
		} else {
			logs.Error("mepAuthentication %s fail %s", mep.MepId, err.Error())
			mep.Token = ""
			mep.RunStatus = mepmd.RunStatusOff
		}
	}

	cols := []string{"mep_name", "user", "pass_word", "province", "city", "user_tag", "longitude", "latitude",
		"contractor", "token", "run_status", "httpsSslType", "mepRootCrt", "mepmClientCrt", "mepmClientKey"}
	err = dao.MepDao.Update(cols, mepmd.MepMeta{MepId: mep.MepId}, *mep)
	if err != nil {
		logs.Error("update mep [%+v] fail %s", mep, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m *mepMg) UpdateStatus(mepId, status string) error {
	bean := mepmd.MepMeta{MepId: mepId}
	if err := dao.MepDao.Get(&bean); err != nil {
		logs.Error("get mep_meta by [%+v] fail %s", bean, err.Error())
		if err == dao.NOTFOUNDGET {
			return common.NewError(common.ErrNotFound)
		}
		return common.NewError(common.ErrDatabase)
	}
	bean = mepmd.MepMeta{MepId: mepId}
	val := mepmd.MepMeta{MgrStatus: status}
	if err := dao.MepDao.Update([]string{"mgr_status"}, bean, val); err != nil {
		logs.Error("update mep mgr_status fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m *mepMg) CheckUpdateStatusReq(v mepmd.MepStatus) error {
	if v.MgrStatus != mepmd.MgrStatusLock && v.MgrStatus != mepmd.MgrStatusUnLock && v.MgrStatus != mepmd.MgrStatusClose {
		return common.NewError(common.ErrParaInvalid, fmt.Sprintf("mgrStatus %s %s", v.MgrStatus, common.MsgInvalidPara))
	}
	return nil
}

func (m *mepMg) Get(mepId string) (mepmd.MepMeta, error) {
	bean := mepmd.MepMeta{MepId: mepId}
	if err := dao.MepDao.Get(&bean); err != nil {
		logs.Error("get mep_meta by [%+v] fail %s", bean, err.Error())
		if err == dao.NOTFOUNDGET {
			return bean, common.NewError(common.ErrNotFound, fmt.Sprintf("mepId %s %s", mepId, common.MsgNotFound))
		}
		return bean, common.NewError(common.ErrDatabase)
	}
	return bean, nil
}

func (m *mepMg) List(from, size int, groupId, excludeGroupId int64, mepName string, mepId string, endPoint string, authMepIds []string, bean mepmd.MepMeta) ([]mepmd.MepMeta, int64, error) {
	var metas []mepmd.MepMeta
	count := int64(0)
	var mepIds []string
	var exMepIds []string
	if groupId != 0 {
		if err := dao.MepDao.ListMepIdWithGroup(groupId, &mepIds); err != nil {
			logs.Error("ListMepIdWithGroup %s fail %s", groupId, err.Error())
			return nil, 0, common.NewError(common.ErrDatabase)
		}
		if len(mepIds) == 0 {
			return metas, count, nil
		}
	}
	if excludeGroupId != 0 {
		if err := dao.MepDao.ListMepIdWithGroup(excludeGroupId, &exMepIds); err != nil {
			logs.Error("ListMepIdWithGroup %s %s", excludeGroupId, err.Error())
			return nil, 0, common.NewError(common.ErrDatabase)
		}
	}

	if err := dao.MepDao.Count(mepName, mepId, endPoint, bean, mepIds, exMepIds, authMepIds, &count); err != nil {
		logs.Error("Count mep fail %s", err.Error())
		return nil, 0, common.NewError(common.ErrDatabase)
	}

	if err := dao.MepDao.List(from, size, mepName, mepId, endPoint, bean, mepIds, exMepIds, authMepIds, &metas); err != nil {
		logs.Error("List mep fail %s", err.Error())
		return nil, 0, common.NewError(common.ErrDatabase)
	}

	//set client cert
	cert := certmd.Cert{Type: certmd.ClientCert}
	if err := dao.CertDao.Get(&cert); err != nil {
		logs.Error("get mepm client cert fail %s", err.Error())
		return nil, 0, common.NewError(common.ErrDatabase)
	}
	//set root cert
	ca := certmd.Cert{Type: certmd.RootCert}
	if err := dao.CertDao.Get(&ca); err != nil {
		logs.Error("get root cert fail %s", err.Error())
		return nil, 0, common.NewError(common.ErrDatabase)
	}
	for _, v := range metas {
		v.RootCrt = ca.Cert
		v.MepmClientCrt = cert.Cert
		v.MepmClientKey = cert.Key
	}
	return metas, count, nil
}

func (m *mepMg) Delete(mepId string) error {
	if err := dao.MepDao.Get(&mepmd.MepMeta{MepId: mepId}); err != nil {
		logs.Error("get mep by mepId %s fail %s", mepId, err.Error())
		if err == dao.NOTFOUNDGET {
			return common.NewError(common.ErrNotFound)
		}
		return common.NewError(common.ErrDatabase)
	}

	if err := dao.MepDao.Delete(mepmd.MepMeta{MepId: mepId}); err != nil {
		logs.Error("delete mep fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}

	if err := dao.DeleteMepGroupRelation(mepmd.MepGroupRelation{MepId: mepId}); err != nil {
		logs.Error("delete relation fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	// 删除权限
	if err := SetAuthForMepId("DELETE", fmt.Sprintf("%s/%s", common.AuthMepAuthorityPath, mepId), "", mepId); err != nil {
		logs.Error("SetAuthForMepId DELETE mepId %s fail %s", mepId, err.Error())
		return err
	}

	//通知alarm
	var code int
	addr := beego.AppConfig.String("AlarmAddress")
	err := gout.POST(addr + "/mepm/alarm/v1/mepmgr/notify/delete").
		SetQuery(map[string]string{"mepId": mepId}).
		Code(&code).
		Do()
	if err != nil || code != http.StatusOK {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		logs.Error("notify mep %s delete to alarm fail,code:%d,err:%s", mepId, code, errMsg)
		return common.NewError(common.ErrInternal)
	}

	return nil
}

func (n *nodeCheckKeeper) Start() {
	logs.Info("start node status check")
	n.wg.Add(1)
	defer n.wg.Done()

	// 启动后第一次更新mep状态
	n.Worker()

	for {
		select {
		case <-n.timer.C:
			n.Worker()
			// 读db，获取当前period
			period, err := getCheckPeriod()
			if err != nil {
				period = LastHeartBeatInterval
			} else {
				LastHeartBeatInterval = period
			}
			n.timer.Reset(time.Second * time.Duration(period))
		case <-n.endChan:
			logs.Info("exit node status checker")
			return
		}
	}
}

func (n *nodeCheckKeeper) Stop() {
	n.timer.Stop()
	close(n.endChan)
	n.wg.Wait()
}

func (n *nodeCheckKeeper) Worker() {
	meps := make([]mepmd.MepMeta, 0, 100)
	match := fmt.Sprintf("mgr_status != '%s'", mepmd.MgrStatusClose)
	if err := dao.MepDao.ListByStatus(0, 0, match, mepmd.MepMeta{}, nil, nil, &meps); err != nil {
		logs.Error("List mep fail %s", err.Error())
		return
	}
	if len(meps) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	con := make(chan interface{}, maxCon)
	defer close(con)
	for _, mep := range meps {
		con <- 1
		wg.Add(1)
		go func(mep mepmd.MepMeta) {
			httpsConf := util.HttpsConf{
				HttpsSslType: mep.HttpsSslType,
				HttpsRootCrt: mep.RootCrt,
			}

			defer func() {
				<-con
				wg.Done()
			}()

			token := mep.Token
			status := mepmd.RunStatusOff
			version := mepmd.DefaultVersion

			for retry := 0; retry <= 1; retry++ {
				if info, err := getMepInfo(mep.EndPoint, token, &httpsConf); err == nil {
					status = mepmd.RunStatusOn
					version = info.MepVersion
					break
				} else if err.Error() == common.ErrAuthFailed {
					if auth, err := mepAuthentication(mep.EndPoint, mep.User, mep.PassWord, &httpsConf); err == nil {
						token = auth.Token
						continue
					}
				}
				break
			}
			updateCols := make([]string, 0, 3)
			if status != mep.RunStatus {
				updateCols = append(updateCols, "run_status")
			}
			if version != mep.SwVersion {
				updateCols = append(updateCols, "sw_version")
			}
			if token != mep.Token {
				updateCols = append(updateCols, "token")
			}

			if len(updateCols) != 0 {
				if err := dao.MepDao.Update(updateCols, mepmd.MepMeta{MepId: mep.MepId}, mepmd.MepMeta{Token: token, RunStatus: status, SwVersion: version}); err != nil {
					logs.Error("Update %v for %s fail %s", updateCols, mep.MepId, err.Error())
				}

				if status != mep.RunStatus && status == mepmd.RunStatusOff {
					//产生告警
					ReportAlarm(&mep)

					//产生MEP上下线记录
					logmd.Log(mep.MepName, "Disconnect")
				} else if status != mep.RunStatus && status == mepmd.RunStatusOn {
					//消除告警
					ResolveAlarm(mep.MepId)
					//调用performance组件采集离线数据
					SupplyOfflineData(mep.MepId)
					//产生MEP上下线记录
					logmd.Log(mep.MepName, "Reconnect")
				}
			}
		}(mep)
	}
	wg.Wait()
}

func mepDoRequest(method, url, host string, header map[string]string, req []byte, httpsConf *util.HttpsConf) ([]byte, error) {
	resp := mepmd.MepCommonResp{}
	code, body, _, err := util.DoRequest(method, url, host, header, req, httpTimeoutMs, httpsConf)
	if err != nil {
		logs.Error("DoRequest url [%s]  fail %v %d %s", url, err, code, string(body))
		return nil, common.NewError(common.ErrNetwork)
	}
	if code == http.StatusUnauthorized {
		logs.Error("mepDoRequest code %s", code)
		return nil, common.NewError(common.ErrAuthFailed)
	}
	if code != http.StatusOK {
		logs.Error("DoRequest url [%s]  code [%s] invalid [%s]", url, code, string(body))
		return nil, common.NewError(common.ErrMepResp)
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		logs.Error("unmarshal resp %s fail %s", string(body), err.Error())
		return nil, common.NewError(common.ErrMepResp)
	}
	if resp.Code != common.ErrSuccess {
		logs.Error("mepDoRequest code %s", resp.Code)
		return nil, common.NewError(resp.Code, resp.Message)
	}

	data, err := json.Marshal(resp.Data)
	if err != nil {
		logs.Error("marshal data [%+v] fail %s", resp.Data, err.Error())
		return nil, common.NewError(common.ErrInternal)
	}
	return data, nil
}

func getMepInfo(endPoint, token string, httpsConf *util.HttpsConf) (mepmd.MepInfo, error) {
	info := mepmd.MepInfo{}
	reqUrl := fmt.Sprintf("%s%s", endPoint, common.MepInfoPath)
	header := map[string]string{common.AuthHeader: "Bearer " + token}

	/*
		code, body, _, err := util.DoRequest("GET", url, "", header, nil, httpTimeoutMs)
		if err != nil || code != http.StatusOK {
			logs.Error("DoRequest url [%s] header [%v] fail %v %d %s", url, header, err, code, string(body))
			if code == http.StatusUnauthorized {
				return info, errors.New(common.ErrAuthFailed)
			} else {
				return info, errors.New("get mep info err")
			}
		}
	*/
	data, err := mepDoRequest("GET", reqUrl, "", header, nil, httpsConf)
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return info, err
	}
	return info, nil
}

func mepAuthentication(endPoint, userPin, passWord string, httpsConf *util.HttpsConf) (mepmd.MepAuth, error) {
	auth := mepmd.MepAuth{}
	uapUrl := beego.AppConfig.String("AlarmAddress")
	performanceUrl := beego.AppConfig.String("PerformanceAddress")
	networkmgrUrl := beego.AppConfig.String("NetworkmgrAddress")
	reqUrl := fmt.Sprintf("%s%s?username=%s&password=%s&mepmUapUrl=%s&mepmPerformanceUrl=%s&mepmNetworkmgrUrl=%s", endPoint, common.MepAuthPath, userPin, passWord, uapUrl, performanceUrl, networkmgrUrl)
	data, err := mepDoRequest("GET", reqUrl, "", nil, nil, httpsConf)
	if err != nil {
		return auth, err
	}
	err = json.Unmarshal(data, &auth)
	if err != nil {
		logs.Error("unMarshal data [%s] fail %s", string(data), err.Error())
		return auth, common.NewError(common.ErrInternal)
	}
	return auth, nil
}

/*
func mepKeepAlive(endPoint, userPin, password string) error {
	encode := util.SaltEncode(password, mepSalt)
	url := fmt.Sprintf("%s%s", endPoint, common.MepKeepAlivePath)
	header := map[string]string{
		"username":     userPin,
		"password":     encode,
		"Content-Type": "application/json;charset=UTF-8",
	}
	code, _, _, err := util.DoRequest("HEAD", url, "", header, nil, httpTimeoutMs)
	if err != nil || code != http.StatusOK {
		logs.Error("DoRequest url [%s] header [%v] fail %s %s", url, header, err, code)
		return fmt.Errorf("mepKeepAlive fail err [%v] code [%s]", err, code)
	}
	return nil
}
*/
