/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package mepmd

import (
	"time"
)

const (
	MgrStatusLock   = "lock"
	MgrStatusUnLock = "unLock"
	MgrStatusClose  = "close"

	RunStatusOn  = "on"
	RunStatusOff = "off"

	DefaultVersion = "-"
)

type MepMeta struct {
	Id            int64     `gorm:"column:id;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	CALLERId      string    `gorm:"column:caller_id" json:"CALLERId"`
	MEPMID        string    `gorm:"column:mepm_id" json:"MEPMID"`
	MepId         string    `gorm:"column:mep_id;INDEX:mepId,UNIQUE" json:"mepId"`
	MepName       string    `gorm:"column:mep_name;INDEX:mepName,UNIQUE" json:"mepName"`
	EndPoint      string    `gorm:"column:end_point;INDEX:endPoint,UNIQUE" json:"endPoint"`
	User          string    `gorm:"column:user" json:"user"`
	PassWord      string    `gorm:"column:pass_word" json:"passWord"`
	Type          string    `gorm:"column:type" json:"type"`
	Province      string    `gorm:"column:province" json:"province"`
	City          string    `gorm:"column:city" json:"city"`
	UserTag       string    `gorm:"column:user_tag" json:"userTag"`
	Longitude     string    `gorm:"column:longitude" json:"longitude"`
	Latitude      string    `gorm:"column:latitude" json:"latitude"`
	Contractor    string    `gorm:"column:contractor" json:"contractor"`
	RunStatus     string    `gorm:"column:run_status;default:'off'" json:"runStatus"`
	MgrStatus     string    `gorm:"column:mgr_status;default:'unLock'" json:"mgrStatus"`
	SwVersion     string    `gorm:"column:sw_version" json:"swVersion"`
	Token         string    `gorm:"column:token;default:''" json:"-"`
	HttpsSslType  int       `gorm:"column:https_ssl_type" json:"httpsSslType"`
	MepCrt        string    `gorm:"column:mep_crt" json:"mepCrt"`
	RootCrt       string    `gorm:"_" json:"rootCrt"`
	MepmClientCrt string    `gorm:"-" json:"mepmClientCrt"`
	MepmClientKey string    `gorm:"-" json:"mepmClientKey"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
}

type MepGroupRelation struct {
	Id         int64  `gorm:"column:id;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	MepId      string `gorm:"column:mep_id;INDEX:group_mep_relation,unique" json:"mepId"`
	MepGroupId int64  `gorm:"column:mep_group_id;INDEX:group_mep_relation,unique" json:"mepGroupId"`
}

type MepStatus struct {
	MgrStatus string `json:"mgrStatus"`
}

type MepInfo struct {
	MepVersion string `json:"mepVersion"`
}

type MepAuth struct {
	Token string `json:"token"`
	MepId string `json:"mepId"`
}

type MepCommonResp struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	ReferInfo interface{} `json:"referInfo"`
	Data      interface{} `json:"data"`
}
