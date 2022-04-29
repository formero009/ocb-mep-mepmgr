/*
@Time : 2022/1/6
@Author : klp
@Project: mepmgr
*/
package mepmd

// param group
const (
	MaintainConfig = "maintainConfig"
)

// param name
const (
	HeartBeatInterval = "heartBeatInterval"

)

type ConfigParameter struct {
	Id        int64       `json:"-" gorm:"column:id;primary_key;auto_increment;not null"`
	ParaGroup string      `json:"paraGroup"`										// 参数组名称
	ParaName  string      `json:"paraName"`										// 参数组名称
	ParaValue string `json:"paraValue"`										// 参数值
	Status    string      `json:"status"`
}
