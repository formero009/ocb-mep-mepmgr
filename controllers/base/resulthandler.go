package base

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/siddontang/go/hack"
	"mepmgr/common"
	"mepmgr/models/logmd"
	erroresult "mepmgr/models/response/errors"
	"net/http"
)

type ResultHandlerController struct {
	OptLog *logmd.OptLog
	beego.Controller
}

type response struct {
	common.ErrMsg
	Result
}

type Result struct {
	Data interface{} `json:"data"`
}

// format BadRequest with param name.
func (c *ResultHandlerController) AbortBadRequestFormat(paramName string) {
	msg := fmt.Sprintf("Invalid param %s !", paramName)
	c.OptLog.LogFail(fmt.Errorf(msg))
	c.AbortBadRequest(msg)
}

func (c *ResultHandlerController) AbortBadRequest(err interface{}) {
	logs.Debug("client-[%s] request-[%s-%s] body-[%s]", c.Ctx.Request.RemoteAddr, c.Ctx.Request.Method, c.Ctx.Request.URL, string(c.Ctx.Input.RequestBody))
	switch err.(type){
	case string:
		logs.Info("Abort BadRequest error. %s", err.(string))
		c.OptLog.LogFail(fmt.Errorf(err.(string)))
		c.CustomAbort(http.StatusBadRequest, hack.String(c.errorResult(http.StatusBadRequest, err.(string))))
	case common.ErrMsg:
		resp := response{err.(common.ErrMsg), Result{}}
		errMsg, errTmp := json.Marshal(resp)
		logs.Error("error %v resp [%+v] msg %s", err, resp, errMsg)
		if errTmp != nil {
			logs.Error("marshal resp %+v fail %s", resp, errTmp.Error())
		}
		c.CustomAbort(http.StatusBadRequest, string(errMsg))
		c.OptLog.LogFail(err.(error))
	default:
		logs.Error("unknown error %v", err)
		errNew := common.NewError(common.ErrUnKnown)
		c.CustomAbort(http.StatusBadRequest, errNew.Error())
	}
}

func (c *BaseController) AbortInternalServerError(err interface{}) {
	logs.Debug("client-[%s] request-[%s-%s] body-[%s]", c.Ctx.Request.RemoteAddr, c.Ctx.Request.Method, c.Ctx.Request.URL, string(c.Ctx.Input.RequestBody))
	switch err.(type) {
	case string:
		logs.Error("Abort InternalServerError error. %s", err.(string))
		c.OptLog.LogFail(fmt.Errorf(err.(string)))
		c.CustomAbort(http.StatusInternalServerError, hack.String(c.errorResult(http.StatusInternalServerError, err.(string))))
	case common.ErrMsg:
		resp := response{err.(common.ErrMsg), Result{}}
		errMsg, errTmp := json.Marshal(resp)
		if errTmp != nil {
			logs.Error("marshal resp %+v fail %s", resp, errTmp.Error())
		}
		c.CustomAbort(http.StatusInternalServerError, string(errMsg))
		c.OptLog.LogFail(err.(error))
	default:
		logs.Error("unknown error %v", err)
		errNew := common.NewError(common.ErrInternal)
		c.CustomAbort(http.StatusInternalServerError, errNew.Error())
	}
}
/*
func (c *ResultHandlerController) AbortConflictRequest(msg string) {
	logs.Error("Abort Conflict error. %s", msg)
	c.OptLog.LogFail(fmt.Errorf(msg))
	c.CustomAbort(http.StatusBadRequest, hack.String(c.errorResult(http.StatusBadRequest, msg)))
}

func (c *ResultHandlerController) AbortNotFoundRequest(msg string) {
	logs.Error("Abort NotFound error. %s", msg)
	c.OptLog.LogFail(fmt.Errorf(msg))
	c.CustomAbort(http.StatusNotFound, hack.String(c.errorResult(http.StatusNotFound, msg)))
}
*/

func (c *ResultHandlerController) errorResult(code int, msg string) []byte {
	errorResult := erroresult.ErrorResult{
		Code: code,
		Msg:  msg,
	}
	body, err := json.Marshal(errorResult)
	if err != nil {
		logs.Error("Json Marshal error. %v", err)
		c.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return body
}

/*
func (c *ResultHandlerController) WriteErrorResponse(errMsg string, code int) {
	c.Data["json"] = errMsg
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.ServeJSON()
}

*/

func (c *ResultHandlerController) Success(data interface{}) {
	logs.Debug("client-[%s] request-[%s-%s] body-[%s]", c.Ctx.Request.RemoteAddr, c.Ctx.Request.Method, c.Ctx.Request.URL, string(c.Ctx.Input.RequestBody))
	c.Ctx.Output.SetStatus(http.StatusOK)
	resp := response{common.NewError(common.ErrSuccess).(common.ErrMsg), Result{data}}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ResultHandlerController) CreateSuccess(data interface{}) {
	logs.Debug("client-[%s] request-[%s-%s] body-[%s]", c.Ctx.Request.RemoteAddr, c.Ctx.Request.Method, c.Ctx.Request.URL, string(c.Ctx.Input.RequestBody))
	c.Ctx.Output.SetStatus(http.StatusCreated)
	resp := response{common.NewError(common.ErrSuccess).(common.ErrMsg), Result{data}}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ResultHandlerController) DeleteSuccess(data ...interface{}) {
	logs.Debug("client-[%s] request-[%s-%s] body-[%s]", c.Ctx.Request.RemoteAddr, c.Ctx.Request.Method, c.Ctx.Request.URL, string(c.Ctx.Input.RequestBody))
	c.Ctx.Output.SetStatus(http.StatusOK)
	resp := response{common.NewError(common.ErrSuccess).(common.ErrMsg), Result{data}}
	c.Data["json"] = resp
	c.ServeJSON()
}
