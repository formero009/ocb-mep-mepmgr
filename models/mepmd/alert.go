package mepmd

import (
	"mepmgr/common"
	"strconv"
)

type Alarm struct {
	AlarmLevel  int    `json:"alarmLevel"`
	AlarmLogo   string `json:"alarmLogo"`
	AlarmName   string `json:"alarmName"`
	AlarmNumber string `json:"alarmNumber"`
	AlarmObject string `json:"alarmObject"`
	AlarmSource string `json:"alarmSource"`
	AlarmTime   int64  `json:"alarmTime"`
	AlarmType   int    `json:"alarmType"` // 告警类型 1:物理资源告警,2:中间件告警,3:服务质量告警 3性能
	NetworkLogo string `json:"networkLogo"`
	NetworkName string `json:"networkName"`
	IsPrompt    bool   `json:"isPrompt"`
}

type AlertInfo struct {
	Id          string `json:"id" gorm:"pk;type(char);size(36);not null"`
	MepId       string `json:"mepId" gorm:"type(char);size(36);not null"`
	AlertObject string `json:"alertObject" gorm:"type(char);size(36);not null"`
	AlertNumber string `json:"alertNumber" gorm:"type(char);size(12);not null"`
	AlertName   string `json:"alertName" gorm:"type(varchar);size(128);not null"`
	AlertType   int    `json:"alertType" gorm:"type(int);not null"`
	Source      string `json:"source" gorm:"type(varchar);size(255);not null"`
	Level       string `json:"level" gorm:"type(char);size(16);not null"`
	StartAt     int64  `json:"startAt" gorm:"type(int64);not null"`
}

func NewAlarm(alertInfo *AlertInfo) *Alarm {
	alarmName := alertInfo.AlertName
	alarmNumber := alertInfo.AlertNumber
	alarmObject := alertInfo.AlertObject
	alarmSource := alertInfo.Source
	alarmTime := alertInfo.StartAt
	alarmType := alertInfo.AlertType
	networkLogo := alertInfo.MepId
	networkName := "MEPName"
	level, _ := strconv.Atoi(alertInfo.Level)

	return &Alarm{
		AlarmLevel:  level,
		AlarmLogo:   alertInfo.Id,
		AlarmName:   alarmName,
		AlarmNumber: alarmNumber,
		AlarmObject: alarmObject,
		AlarmSource: alarmSource,
		AlarmTime:   alarmTime,
		AlarmType:   alarmType,
		NetworkLogo: networkLogo,
		NetworkName: networkName,
		IsPrompt:    true,
	}
}

type MepProcessLog struct {
	MepName string `json:"mepName"`
	LogType string `json:"logType"`
	LogTime int64  `json:"logTime"`
}

type MepQuery struct {
	common.QueryParam
	StartTime int64
	EndTime   int64
	MepName   string
	LogType   string
}
