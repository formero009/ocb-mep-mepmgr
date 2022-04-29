/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package dao

import (
	"fmt"
	"mepmgr/models"
	"mepmgr/models/mepmd"
)

var MepDao mepMeta

type mepMeta struct{}

func (m mepMeta) Create(value *mepmd.MepMeta) error {
	return models.PostgresDB.Create(value).Error
}

func (m mepMeta) Update(cols []string, bean, value mepmd.MepMeta) error {
	tx := models.PostgresDB.Model(&mepmd.MepMeta{}).Where(bean)
	if cols != nil && len(cols) > 0 {
		tx = tx.Select(cols)
	}
	return tx.Updates(value).Error
}

func (m mepMeta) Get(bean *mepmd.MepMeta) error {
	//return models.PostgresDB.Take(bean, bean).Error
	return models.PostgresDB.Where(bean).First(bean).Error
}

func (m mepMeta) GetMepInfosByMepNames(mepNames []string, metas *[]mepmd.MepMeta) error {
	return models.PostgresDB.Where("mep_name IN ?", mepNames).Find(metas).Error
}

func (m mepMeta) List(from, size int, mepName string, mepId string, endPoint string, bean mepmd.MepMeta, mepIds, exMepIds, authMepIds []string, metas *[]mepmd.MepMeta) error {
	tx := models.PostgresDB
	if from > 0 && size > 0 {
		tx = tx.Limit(size).Offset((from - 1) * size)
	}
	if mepName != "" {
		tx = tx.Where("mep_name like ?", mepName)
	}
	if mepId != "" {
		tx = tx.Where("mep_id like ?", mepId)
	}
	if endPoint != "" {
		tx = tx.Where("end_point like ?", endPoint)
	}
	if mepIds != nil {
		tx = tx.Where("mep_id in (?)", mepIds)
	}
	if exMepIds != nil && len(exMepIds) > 0 {
		tx = tx.Where("mep_id not in (?)", exMepIds)
	}
	if authMepIds != nil {
		tx = tx.Where("mep_id in (?)", authMepIds)
	}
	tx = tx.Order("created_at DESC")
	return tx.Find(metas, bean).Error
}

func (m mepMeta) ListByStatus(from, size int, matchStr string, bean mepmd.MepMeta, mepIds, exMepIds []string, metas *[]mepmd.MepMeta) error {
	tx := models.PostgresDB
	if from > 0 && size > 0 {
		tx = tx.Limit(size).Offset((from - 1) * size)
	}
	if matchStr != "" {
		tx = tx.Where(matchStr)
	}

	if mepIds != nil {
		tx = tx.Where("mep_id in (?)", mepIds)
	}
	if exMepIds != nil && len(exMepIds) > 0 {
		tx = tx.Where("mep_id not in (?)", exMepIds)
	}
	tx = tx.Order("created_at DESC")
	return tx.Find(metas, bean).Error
}

func (m mepMeta) Count(mepName string, mepId string, endPoint string, bean mepmd.MepMeta, mepIds, exMepIds, authMepIds []string, count *int64) error {
	tx := models.PostgresDB.Model(&mepmd.MepMeta{})
	if mepName != "" {
		tx = tx.Where("mep_name like ?", mepName)
	}
	if mepId != "" {
		tx = tx.Where("mep_id like ?", mepId)
	}
	if endPoint != "" {
		tx = tx.Where("end_point like ?", endPoint)
	}
	if mepIds != nil {
		tx = tx.Where("mep_id in (?)", mepIds)
	}
	if exMepIds != nil && len(exMepIds) > 0 {
		tx = tx.Where("mep_id not in (?)", exMepIds)
	}
	if authMepIds != nil {
		tx = tx.Where("mep_id in (?)", authMepIds)
	}
	err := tx.Where(bean).Count(count).Error
	return err
}

func (m mepMeta) Delete(bean mepmd.MepMeta) error {
	return models.PostgresDB.Where(bean).Delete(mepmd.MepMeta{}).Error
}

func (m mepMeta) ListMepIdWithGroup(groupId int64, mepIds *[]string) error {
	var relations []mepmd.MepGroupRelation
	if err := models.PostgresDB.Find(&relations, mepmd.MepGroupRelation{MepGroupId: groupId}).Error; err != nil {
		return err
	}

	for _, relation := range relations {
		*mepIds = append(*mepIds, relation.MepId)
	}
	return nil
}

func CreateMepGroupRelation(value *[]mepmd.MepGroupRelation) error {
	return models.PostgresDB.Create(value).Error
}

func DeleteMepGroupRelation(bean mepmd.MepGroupRelation) error {
	return models.PostgresDB.Where(bean).Delete(mepmd.MepGroupRelation{}).Error
}

func GetMepGroupRelation(bean *mepmd.MepGroupRelation) error {
	return models.PostgresDB.Where(bean).First(bean).Error
}

func GetMepGroupRelations(matchStr string, bean *[]mepmd.MepGroupRelation) error {
	return models.PostgresDB.Raw(fmt.Sprintf("SELECT * FROM mep_group_relation WHERE (mep_id, mep_group_id) IN %v;", matchStr)).Scan(bean).Error
}
