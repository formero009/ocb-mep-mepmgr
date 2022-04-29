package base

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/manager"
	"mepmgr/models/logmd"
	"mepmgr/models/usermd"
	"mepmgr/util"
	"net/url"
	"strings"
	"time"
)

type BaseController struct {
	ParamBuilderController
	AuthInfo
}

type AuthInfo struct {
	AuthMepInfo
}

type AuthMepInfo struct {
	MepIds []string
	IsAuth bool
}

var method2Op = map[string]string{
	"GET":    "查询",
	"POST":   "新增",
	"PUT":    "更新",
	"DELETE": "删除",
}

var path2MetaOp = map[string]map[string]string{
	"mepGroup": mepGroupPath2Op,
	"topology": mepTopologyPath2Op,
}

var mepPath2Op = map[string]string{
	":mepId/status": "mep状态",
	":mepId":        "mep",
	"/":             "mep",
	"":              "mep",
}

var mepGroupPath2Op = map[string]string{
	":mepGroupName/mep/:mepName": "mep分组下的mep",
	":mepGroupName/mep":          "mep分组下的mep",
	":mepGroupName":              "mep分组",
	"/":                          "mep分组",
	"":                           "mep分组",
}

var mepTopologyPath2Op = map[string]string{
	"detail": "mep详细信息",
	"/":      "mep拓扑",
	"":       "mep拓扑",
}

func (c *BaseController) analyzeOperation() {
	method := method2Op[c.Ctx.Request.Method]
	operation := "未知"
	mepId := ""

	inputUrls := strings.Split(strings.TrimSuffix(strings.TrimPrefix(c.Ctx.Input.URL(), "/mepm/v1/"), "/"), "/")
	if len(inputUrls) < 2 {
		return
	}

	urlIndex := 3
	if len(inputUrls) < 3 {
		urlIndex = 2
	}
	inputPaths := inputUrls[urlIndex:]
	if len(inputPaths) == 0 {
		inputPaths = []string{""}
	}

	path2Op, f := path2MetaOp[inputUrls[urlIndex-1]]
	if !f {
		path2Op = mepPath2Op
	}

	for path, op := range path2Op {
		paths := strings.Split(path, "/")
		if len(paths) != len(inputPaths) {
			continue
		}
		match := true
		for k := 0; k < len(paths); k++ {
			if paths[k] == ":mepId" {
				mepId = inputPaths[k]
			}
			if !strings.HasPrefix(paths[k], ":") && paths[k] != inputPaths[k] {
				match = false
				break
			}
		}
		if match {
			operation = method + op
			c.OptLog.MepId = mepId
			break
		}
	}
	c.OptLog.Content = operation
	return

}

func (c *BaseController) Prepare() {
	//init OptLog
	c.OptLog = &logmd.OptLog{StartTime: time.Now()}

	//get ip from req
	ip, err := util.GetIP(c.Ctx.Request)
	if err != nil {
		logs.Error("failed to get ip from request,", err)
	}
	c.OptLog.Ip = ip
	c.analyzeOperation()

	c.AnalyzeAuthInfo()
	//get userInfo from header
	userInfoEncode := c.Ctx.Request.Header.Get(common.UserInfo)
	if userInfoEncode == "" {
		return
	}
	userJsonStr, err := url.PathUnescape(userInfoEncode)
	if err != nil {
		logs.Error("failed to decode userInfo,", err)
	}
	userInfo := &usermd.UserInfo{}
	err = json.Unmarshal([]byte(userJsonStr), userInfo)
	if err != nil {
		logs.Error("failed to unmarshal userInfo,", err)
	}

	//fill userInfo to OptLog
	c.OptLog.Account = userInfo.Username
	c.OptLog.UserId = userInfo.UserId
	c.OptLog.Name = userInfo.Name
	c.OptLog.Role = strings.Join(userInfo.Roles, ",")
}

func (c *BaseController) Finish() {
	//init OptLog
	c.OptLog.LogSuccess()
}

//page struct
type Page struct {
	PageNo     int         `json:"currentPage"`
	PageSize   int         `json:"pageSize"`
	TotalPage  int64       `json:"totalPage"`
	TotalCount int64       `json:"totalCount"`
	List       interface{} `json:"content"`
}

func (c *BaseController) NewPage(count int64, list interface{}, pageNo int, pageSize int) *Page {
	tp := count / int64(pageSize)
	if count%int64(pageSize) > 0 {
		tp = count/int64(pageSize) + 1
	}
	return &Page{
		PageNo:     pageNo,
		PageSize:   pageSize,
		TotalPage:  tp,
		TotalCount: count,
		List:       list,
	}
}

func (c *BaseController) AnalyzeAuthInfo () {
	mepIdFromHeader, fMepIdFromH := c.Ctx.Request.Header[manager.HeaderMepId]
	c.AuthInfo.AuthMepInfo.IsAuth = fMepIdFromH
	if fMepIdFromH {
		if mepIdFromHeader == nil {
			mepIdFromHeader = []string{""}
		}
		c.AuthInfo.AuthMepInfo.MepIds = make([]string, 0, 10)
		for _, id := range mepIdFromHeader {
			c.MepIds = append(c.MepIds, strings.Split(id, ",")...)
		}
	}
}
