/*
@Time : 2022/4/27
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package manager

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"gorm.io/gorm"
	"io/ioutil"
	"mepmgr/common"
	"mepmgr/dao"
	"mepmgr/models"
	"mepmgr/models/certmd"
	"mepmgr/util"
	"strings"
	"sync"
	"time"
)

var alarmInfo interface{}

func init() {
	//[cert type][cert status][alarm info]
	//alarmInfo := make(map[int]map[int]string)
	//alarmInfo[certmd.RootCert][certmd.StatusToBeInvalid] = "根证书即将过期"
	//alarmInfo[certmd.RootCert][certmd.StatusInvalid] = "根证书已经过期"
	//
	//alarmInfo[certmd.ServerCert][certmd.StatusToBeInvalid] = "mepm服务端证书即将过期"
	//alarmInfo[certmd.ServerCert][certmd.StatusToBeInvalid] = "mepm服务端证书已经过期"
	//
	//alarmInfo[certmd.ClientCert][certmd.StatusToBeInvalid] = "mepm客户端证书即将过期"
	//alarmInfo[certmd.ClientCert][certmd.StatusToBeInvalid] = "mepm客户端证书已经过期"
}

type CertManager interface {
	TimerCheck()
	TimerCheckStop()
	Update() error
	List() ([]certmd.Cert, error)
	LoginCheck() (*CheckRsp, error)
}

type CheckRsp struct {
	MepmServerCert int
	MepmClientCert int
	RootCert       int
}

type certChecker struct {
	timer   *time.Timer
	endChan chan bool
}

var certMgOnce sync.Once
var certMg defaultCertMg

func NewDefaultCertManager() CertManager {
	certMgOnce.Do(func() {
		c := certChecker{}
		//reuse node check period
		period, err := getCheckPeriod()
		if err != nil {
			panic(fmt.Sprintf("getCheckPeriod fail %s", err.Error()))
		}
		LastHeartBeatInterval = period
		c.timer = time.NewTimer(time.Second * time.Duration(period))
		c.endChan = make(chan bool, 1)
		certMg = defaultCertMg{certChecker: c}
	})
	return certMg
}

type defaultCertMg struct {
	certChecker certChecker
}

func (d defaultCertMg) TimerCheck() {
	logs.Info("start cert status check timer")
	for {
		select {
		case <-d.certChecker.timer.C:
			d.certChecker.check()
			// 读db，获取当前period
			period, err := getCheckPeriod()
			if err != nil {
				period = LastHeartBeatInterval
			} else {
				LastHeartBeatInterval = period
			}
			d.certChecker.timer.Reset(time.Second * time.Duration(period))
		case <-d.certChecker.endChan:
			logs.Info("exit cert status checker")
			return
		}
	}
}
func (d defaultCertMg) TimerCheckStop() {
	d.certChecker.timer.Stop()
	close(d.certChecker.endChan)
}

func (d defaultCertMg) Update() error {
	panic("implement me")
}

func (d defaultCertMg) List() ([]certmd.Cert, error) {
	var certs []certmd.Cert
	if err := dao.CertDao.List(&certs); err != nil {
		logs.Error("list mepm certs err, %s", err.Error())
		return nil, common.NewError(common.ErrDatabase, fmt.Sprintf("list mepm certs err"))
	}
	return certs, nil
}

func (d defaultCertMg) LoginCheck() (*CheckRsp, error) {
	rsp := CheckRsp{}
	if certs, err := d.List(); err != nil {
		logs.Error("get mepm certs from database err, %s", err.Error())
		return nil, common.NewError(common.ErrDatabase, fmt.Sprintf("list mepm certs err"))
	} else {
		for _, cert := range certs {
			switch cert.Type {
			case certmd.ClientCert:
				rsp.MepmClientCert = getStatus(cert.ValidTime)
				break
			case certmd.ServerCert:
				rsp.MepmServerCert = getStatus(cert.ValidTime)
				break
			case certmd.RootCert:
				status := getStatus(cert.ValidTime)
				//update c/s cert to invalid
				if status == certmd.StatusInvalid {
					clientCert := &certmd.Cert{Type: certmd.ClientCert, Status: status}
					dao.CertDao.Update(clientCert)
					serverCert := &certmd.Cert{Type: certmd.ServerCert, Status: status}
					dao.CertDao.Update(serverCert)
				}
				rsp.RootCert = status
			}
		}
	}
	return &rsp, nil
}

func (c certChecker) check() {
	var certs []certmd.Cert
	if err := dao.CertDao.List(&certs); err != nil {
		logs.Error("get mepm cert list err, %s", err.Error())
		close(c.endChan)
		return
	}
	for _, cert := range certs {
		if getStatus(cert.ValidTime) != certmd.StatusValid {
			doAlarmReport(cert)
		}
	}
}

func doAlarmReport(cert certmd.Cert) {
	return
}

func LoadCerts() error {
	certConf := make(map[int]string, 3)
	certConf[certmd.ClientCert] = "ClientCert,ClientKey"
	certConf[certmd.ServerCert] = "ServerCert,ServerKey"
	certConf[certmd.RootCert] = "RootCert,RootKey"
	var certs []*certmd.Cert
	for types, conf := range certConf {
		confs := strings.Split(conf, ",")
		if cert, err := load(types, confs[0], confs[1]); err != nil {
			return err
		} else {
			certs = append(certs, cert)
		}
	}
	return models.PostgresDB.Transaction(func(tx *gorm.DB) error {
		tx.Exec("truncate table cert")
		return tx.CreateInBatches(certs, len(certs)).Error
	})
}

func load(types int, cert string, key string) (*certmd.Cert, error) {
	certPath := beego.AppConfig.String(cert)
	certContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		logs.Error("load %s err, %v", cert, err.Error())
		return nil, err
	}
	x509Cert, err := util.ParseCert(string(certContent))
	if err != nil {
		logs.Error("parse %s err, %v", cert, err.Error())
		return nil, err
	}
	keyPath := beego.AppConfig.String(key)
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logs.Error("load %s err, %v", key, err.Error())
		return nil, err
	}
	return &certmd.Cert{Cert: string(certContent),
		Key:       string(keyContent),
		Type:      types,
		ValidTime: x509Cert.NotAfter,
		Status:    getStatus(x509Cert.NotAfter)}, nil
}

func getStatus(validTime time.Time) int {
	if time.Now().After(validTime) {
		return certmd.StatusInvalid
	} else if time.Now().Add(time.Hour * 72).After(validTime) {
		return certmd.StatusToBeInvalid
	} else {
		return certmd.StatusValid
	}
}
