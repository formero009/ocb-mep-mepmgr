/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package mepmd

import "time"

type MepGroup struct {
	Id           int64     `gorm:"column:id;primary_key;auto_increment;not null" json:"id"`
	CALLERId     string    `gorm:"column:caller_id" json:"CALLERId"`
	MEPMID       string    `gorm:"column:mepm_id" json:"MEPMID"`
	MepGroupName string    `gorm:"column:mep_group_name;index:mep_group_name_idx;unique;not null" json:"mepGroupName"`
	Description  string    `gorm:"column:description" json:"description"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type MepGroupList struct {
	PageNo     int        `json:"pageNo"`
	PageSize   int        `json:"pageSize"`
	TotalCount int64      `json:"totalCount"`
	Meps       []MepGroup `json:"meps"`
}

type MepGroupUpdate struct {
	CALLERId     string `json:"CALLERId"`
	MEPMID       string `json:"MEPMID"`
	MepGroupName string `json:"mepGroupName"`
	Description  string `json:"description"`
}

type MepGroupAddMep struct {
	MepNames []string `json:"mepNames"`
}
