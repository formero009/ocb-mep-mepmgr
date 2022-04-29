/*
@Time : 2022/4/28
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package util

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/astaxie/beego/logs"
)

func ParseCert(crt string) (*x509.Certificate, error) {
	certDERBlock, _ := pem.Decode([]byte(crt))
	if certDERBlock == nil {
		logs.Error("parse crt file error")
		return nil, errors.New("parse crt file error")
	}
	return x509.ParseCertificate(certDERBlock.Bytes)
}
