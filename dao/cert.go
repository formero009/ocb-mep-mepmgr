/*
@Time : 2022/4/27
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package dao

import (
	"mepmgr/models"
	"mepmgr/models/certmd"
)

var CertDao cert

type cert struct{}

func (m cert) Create(value *certmd.Cert) error {
	return models.PostgresDB.Create(value).Error
}

func (m cert) BatchCreate(values []*certmd.Cert) error {
	return models.PostgresDB.CreateInBatches(values, len(values)).Error
}

func (m cert) Update(value *certmd.Cert) error {
	return models.PostgresDB.Updates(value).Error
}
func (m cert) List(values *[]certmd.Cert) error {
	return models.PostgresDB.Find(values).Error
}

func (m cert) Get(cert *certmd.Cert) error {
	return models.PostgresDB.Take(cert, cert).Error
}
