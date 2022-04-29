/*
@Time : 2022/1/11
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package manager

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/models/mepmd"
	"mepmgr/util"
	"net/http"
	"strings"
)

func SupplyOfflineData(mepId string) {
	//request POST [performance]/mepm/performance/v1/data/collect/supplyAll
	resp := mepmd.MepCommonResp{}
	perfAddr := beego.AppConfig.DefaultString("PerformanceAddress", "http://127.0.0.1:8014")
	url := strings.TrimSuffix(perfAddr, "/") + "/mepm/performance/v1/data/collect/supplyAll"
	mep := struct {
		Id string
	}{Id: mepId}
	req, _ := json.Marshal(mep)
	code, body, _, err := util.DoRequest("POST", url, perfAddr, nil, req, httpTimeoutMs, &util.HttpsConf{HttpsSslType: util.HttpsAuthNone})
	if err != nil {
		logs.Error("DoRequest url [%s]  fail %v %d %s", url, err, code, string(body))
	}
	if code != http.StatusOK {
		logs.Error("DoRequest url [%s]  code [%s] invalid [%s]", url, code, string(body))
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		logs.Error("unmarshal resp %s fail %s", string(body), err.Error())
	}
	if resp.Code != common.ErrSuccess {
		logs.Error("supply offline performance err, err code %v, msg %v, data [%+v] fail", resp.Code, resp.Message, resp.Data)
	}
}
