package dao

import (
	"errors"
	"gorm.io/gorm"
	"mepmgr/models"
	"mepmgr/models/mepmd"
)

type AlertDao interface {
	Create(alertDb *mepmd.AlertInfo) error
	GetByMepId(mepId string, alertInfo *mepmd.AlertInfo)
	DeleteByMepId(mepId string) error

	LogCreate(log *mepmd.MepProcessLog) error
	LogGet(q mepmd.MepQuery) (int64, []mepmd.MepProcessLog, error)
}

type DefaultAlertDao struct {
}

func NewDefaultDao() *DefaultAlertDao {
	return &DefaultAlertDao{}
}

func (d *DefaultAlertDao) Create(alertDb *mepmd.AlertInfo) error {
	var err error
	err = models.PostgresDB.Create(alertDb).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DefaultAlertDao) GetByMepId(mepId string, alertInfo *mepmd.AlertInfo) error {
	if err := models.PostgresDB.Model(mepmd.AlertInfo{}).Where("mep_id = ?", mepId).First(&alertInfo).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil
		}
		return err
	}
	return nil
}

func (d *DefaultAlertDao) DeleteByMepId(mepId string) error {
	return models.PostgresDB.Model(&mepmd.AlertInfo{}).Where("mep_id = ?", mepId).Delete(&mepmd.AlertInfo{}).Error
}

func (d *DefaultAlertDao) LogCreate(log *mepmd.MepProcessLog) error {
	return models.PostgresDB.Model(&mepmd.MepProcessLog{}).Create(log).Error
}

func (d *DefaultAlertDao) LogGet(q mepmd.MepQuery) (int64, []mepmd.MepProcessLog, error) {
	tempDb := models.PostgresDB
	var mepProcessLog []mepmd.MepProcessLog
	var count int64

	if q.MepName != "" {
		tempDb = tempDb.Where("mep_name = ?", q.MepName)
	}

	if q.LogType != "" {
		tempDb = tempDb.Where("log_type = ?", q.LogType)
	}

	if q.StartTime != 0 && q.EndTime != 0 {
		tempDb = tempDb.Where("log_time >= ? and log_time <= ?", q.StartTime, q.EndTime)
	}

	err := tempDb.Model(&mepmd.MepProcessLog{}).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}
	err = tempDb.Model(&mepmd.MepProcessLog{}).Limit(int(q.Limit())).Offset(int(q.Offset())).Find(&mepProcessLog).Error
	return count, mepProcessLog, err
}
