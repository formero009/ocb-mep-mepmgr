/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package dao

import (
	"gorm.io/gorm"
	"mepmgr/models"
	"mepmgr/models/mepmd"
)

var MepGroupDao mepGroup

type mepGroup struct{}

func (m mepGroup) Create(tx *gorm.DB, mg *mepmd.MepGroup) error {
	if tx == nil {
		tx = models.PostgresDB
	}
	return tx.Create(mg).Error
}

func (m mepGroup) Delete(bean mepmd.MepGroup) error {
	return models.PostgresDB.Where(bean).Delete(mepmd.MepGroup{}).Error
}

func (m mepGroup) Update(cols []string, bean, value mepmd.MepGroup) error {
	tx := models.PostgresDB.Model(&mepmd.MepGroup{}).Where(bean)
	if cols != nil && len(cols) > 0 {
		tx = tx.Select(cols)
	}
	return tx.Updates(value).Error
}

func (m mepGroup) Get(bean *mepmd.MepGroup) error {
	return models.PostgresDB.Where(bean).First(bean).Error
}

func (m mepGroup) List(index, size int, matchStr string, bean mepmd.MepGroup, meps *[]mepmd.MepGroup) error {
	db := models.PostgresDB
	if size > 0 {
		db = db.Limit(size)
		if index > 0 {
			db = db.Offset((index - 1) * size)
		}
	}

	if matchStr != "" {
		db = db.Where("mep_group_name like ?", matchStr)
	}
	db = db.Order("created_at DESC")
	return db.Where(bean).Find(meps).Error
}

func (m mepGroup) Count(matchStr string, mepg mepmd.MepGroup, count *int64) error {
	db := models.PostgresDB
	if len(matchStr) != 0 {
		db = db.Where("mep_group_name like ?", matchStr)
	}
	return db.Model(mepg).Where(mepg).Count(count).Error
}
