/*
@Time : 2022/1/6
@Author : klp
@Project: mepmgr
*/
package dao

import (
	"mepmgr/models"
	"mepmgr/models/mepmd"
)

var ConfigParamDao configParam

type configParam struct{}

func (c configParam) Create(param *mepmd.ConfigParameter) error {
	return models.PostgresDB.Create(param).Error
}

func (c configParam) Update(bean, param mepmd.ConfigParameter) error {
	return models.PostgresDB.Model(&mepmd.ConfigParameter{}).Where(&bean).Update("para_value", param.ParaValue).Error
}

func (c configParam) Get(param *mepmd.ConfigParameter) error {
	return models.PostgresDB.Take(param, param).Error
}

func (c configParam) List(from, size int, order, matchStr string, bean mepmd.ConfigParameter, params *[]mepmd.ConfigParameter) error {
	db := models.PostgresDB
	if size > 0 {
		db = db.Limit(size)
		if from > 0 {
			db = db.Offset((from - 1) * size)
		}
	}
	if matchStr != "" {
		db = db.Where(matchStr)
	}

	if order != "" {
		db = db.Order(order)
	}
	return db.Where(bean).Find(params).Error
}

func (c configParam) Count(bean mepmd.ConfigParameter, count *int64) error {
	return models.PostgresDB.Model(&mepmd.ConfigParameter{}).Where(bean).Count(count).Error
}
