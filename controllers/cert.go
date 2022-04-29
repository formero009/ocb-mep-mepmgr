/*
@Time : 2022/4/27
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mepmgr/common"
	"mepmgr/controllers/base"
	"mepmgr/dao"
	"mepmgr/manager"
	"mepmgr/models/certmd"
	"mepmgr/util"
	"time"
)

type CertController struct {
	base.BaseController
	CertMg manager.CertManager
}

// @Title Update cert file
// @Description Update mep node
// @Success 201 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 403 body is empty
// @Failure 500 server internal error
// @router /:certId [put]
func (s *CertController) Update() {
	certId := s.Ctx.Input.Param(":certId")
	//查询是否存在
	bean := certmd.Cert{Id: certId}
	if err := dao.CertDao.Get(&bean); err != nil {
		s.AbortBadRequest(err.Error())
		return
	}

	//解析
	var params certmd.Cert
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &params); err != nil {
		logs.Error("parse body error. %v", err)
		s.AbortBadRequest(common.NewError(common.ErrParaJson))
		return
	}
	if params.Key == "" || params.Cert == "" {
		s.AbortBadRequestFormat("key or cert can't be empty")
		return
	}

	//非服务端证书跳过
	if params.Type != 0 {
		s.AbortBadRequest(common.NewError(common.ErrParaInvalid))
		return
	}

	//保存为文件
	serverCert := beego.AppConfig.String("ServerCert")
	serverKey := beego.AppConfig.String("ServerKey")
	err := util.WriteFile(params.Cert, serverCert)
	if err != nil {
		s.AbortInternalServerError(err.Error())
		return
	}
	err = util.WriteFile(params.Key, serverKey)
	if err != nil {
		s.AbortInternalServerError(err.Error())
		return
	}

	//reload nginx
	dockerSock := beego.AppConfig.String("DockerSock")
	dockerCommand := beego.AppConfig.String("DockerCommand")
	dockerContainerName := beego.AppConfig.String("DockerContainerName")
	err = util.ReloadNgninx(dockerSock, dockerCommand, dockerContainerName)
	if err != nil {
		logs.Error(err.Error())
		s.AbortInternalServerError(err.Error())
		return
	}

	//都成功更新数据库
	//获取证书有效期
	t, err := util.ParseCert(params.Cert)
	if err != nil {
		s.AbortInternalServerError(err.Error())
		return
	}

	var cert certmd.Cert
	cert.Cert = params.Cert
	cert.Key = params.Key
	cert.Id = certId
	cert.ValidTime = t.NotAfter //过期时间
	cert.UpdatedAt = time.Now()

	err = dao.CertDao.Update(&cert)
	if err != nil {
		logs.Error(err.Error())
		s.AbortInternalServerError(err.Error())
		return
	}

	s.Success(cert)
}

// @Title List cert info
// @Description List cert info
// @Success 200 {string} object msg
// @Failure 401 unauthorized
// @Failure 500 server internal error
// @router  /list [get]
func (s *CertController) List() {
	certs, err := s.CertMg.List()
	if err != nil {
		logs.Error("list mepm certs err, %s", err)
		s.AbortInternalServerError(err)
	}
	s.Success(certs)
}

// @Title Check mepm certs status
// @Description mepm certs status
// @Success 200 {string} object msg
// @Failure 401 unauthorized
// @Failure 400 app parse error
// @Failure 500 server internal error
// @router /check [get]
func (s *CertController) Check() {
	//rsp, err := s.CertMg.Check()
	//if err != nil {
	//	logs.Error("list mepm certs err, %s", err)
	//	s.AbortInternalServerError(err)
	//}
	//s.Success(rsp)
}
