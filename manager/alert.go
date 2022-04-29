package manager

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"mepmgr/dao"
	"mepmgr/models/mepmd"
	"mepmgr/util"
	"net/http"
	"strings"
	"time"
)

const serviceAlert = 3

const (
	alertName   = "MEP异常离线"
	alertNumber = "202103160031"
	alertObject = "MEP"
	alertLevel  = "0" //严重告警
	alertType   = serviceAlert
)

type Client struct {
	addr string
}

type mepMeta struct {
	MepId string
}

type MepProcessLogManager interface {
	Get(q mepmd.MepQuery) (int64, []mepmd.MepProcessLog, error)
}

type defaultMepProcessLogManage struct {
	*dao.DefaultAlertDao
}

func NewDefaultMepProcessLogManager() *defaultMepProcessLogManage {
	return &defaultMepProcessLogManage{
		DefaultAlertDao: dao.NewDefaultDao(),
	}
}

func (m *defaultMepProcessLogManage) Get(q mepmd.MepQuery) (int64, []mepmd.MepProcessLog, error) {
	return m.DefaultAlertDao.LogGet(q)
}

func NewUapClient() (*Client, error) {
	addr := beego.AppConfig.String("AlarmAddress")
	logs.Info("uap_server: start connecting uap server, addr %v\n", addr)
	return &Client{addr: addr}, nil
}

func ReportAlarm(mep *mepmd.MepMeta) {
	alertInfoDb := &mepmd.AlertInfo{}
	if err := dao.NewDefaultDao().GetByMepId(mep.MepId, alertInfoDb); err != nil {
		logs.Error("failed to query in db: %v", err)
		return
	}
	if alertInfoDb.Id != "" {
		if err := dao.NewDefaultDao().DeleteByMepId(mep.MepId); err != nil {
			logs.Error("failed to delete alert info in db: %v", err)
			return
		}
	}

	alertInfo := &mepmd.AlertInfo{
		Id:          util.UUID(),
		MepId:       mep.MepId,
		AlertName:   alertName,
		AlertNumber: alertNumber,
		AlertObject: alertObject,
		AlertType:   alertType,
		Level:       alertLevel,
		Source:      mep.EndPoint,
		StartAt:     time.Now().Unix(),
	}
	if err := dao.NewDefaultDao().Create(alertInfo); err != nil {
		logs.Error("failed to create alert info in db.")
		return
	}

	client, err := NewUapClient()
	if err != nil {
		logs.Error("Create uap server err. %v", err.Error())
		return
	}

	alarm := mepmd.NewAlarm(alertInfo)
	resp, err := client.DoRequest("POST", "/v1/events", nil, alarm)
	logs.Info("get resp from uap server")
	if err != nil {
		logs.Error("Get resp from uap server err. %v", err.Error())
		return
	}

	defer resp.Body.Close()
	logs.Info("Check resp code from mep, %v", resp)
	if resp.StatusCode != 200 {
		logs.Error("Get wrong response code from mep. code %v, message %v", resp.StatusCode, resp.Status)
		return
	}
}

func ResolveAlarm(mepId string) {
	client, err := NewUapClient()
	if err != nil {
		logs.Error("Create uap server err. %v", err.Error())
		return
	}
	alertInfo := mepmd.AlertInfo{}
	if err := dao.NewDefaultDao().GetByMepId(mepId, &alertInfo); err != nil {
		logs.Error("failed to query in db: %v", err)
		return
	}

	if alertInfo.Id != "" {
		if err := dao.NewDefaultDao().DeleteByMepId(mepId); err != nil {
			logs.Error("failed to delete alert info in db: %v, still going to send delete request to alarm service.", err)
		}
		resp, err := client.DoRequest("DELETE", "/v1/event/"+alertInfo.Id+"/clear", nil, nil)
		logs.Info("get resp from uap server")
		if err != nil {
			logs.Error("Get resp from uap server err. %v", err.Error())
			return
		}

		defer resp.Body.Close()
		logs.Info("Check resp code from mep, %v", resp)
		if resp.StatusCode != 200 {
			logs.Error("Get wrong response code from mep. code %v, message %v", resp.StatusCode, resp.Status)
			return
		}
	} else {
		logs.Error("no alert info fund in db")
	}
}

func (c *Client) DoRequest(method string, path string, reqUrl map[string]interface{}, reqBody interface{}) (repo *http.Response, err error) {
	baseUrl := "/mepm/alarm"
	logs.Debug("Request to mepm method: %v, path: %v", method, baseUrl+path)

	client := &http.Client{Timeout: time.Second * 8}
	var reader io.Reader
	if reqBody != nil {
		bytesData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		logs.Debug("Request to mepm body: %v", string(bytesData))
		reader = bytes.NewReader(bytesData)
	}

	//process url request param
	var urlParam = "?"
	if reqUrl != nil {
		for k, v := range reqUrl {
			if v != "" {
				urlParam += k + "=" + v.(string) + "&"
			}
		}
	}

	var req *http.Request
	upMethod := strings.ToUpper(method)
	address := c.addr + baseUrl + path
	final := address + urlParam[0:len(urlParam)-1]
	if upMethod == "POST" || upMethod == "PUT" || upMethod == "DELETE" {
		req, err = http.NewRequest(method, final, reader)
	} else {
		req, err = http.NewRequest(method, final, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	repo, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
