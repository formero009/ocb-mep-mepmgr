package dao

import (
	"mepmgr/models"
	"mepmgr/models/mepmd"
)

type preferenceDao struct {
}

var PreferenceDao preferenceDao

func (preferenceDao) Create(value *mepmd.Preference) error {
	return models.PostgresDB.Create(value).Error
}

func (preferenceDao) GetByUserName(userName string) (*mepmd.Preference, error) {
	preference := &mepmd.Preference{}
	err := models.PostgresDB.Where("user_name=?", userName).First(preference).Error
	return preference, err
}

func (preferenceDao) Update(p *mepmd.Preference) error {
	err := models.PostgresDB.Model(p).
		Where("user_name=?", p.UserName).
		Update("node_type", p.NodeType).
		Update("node_id", p.NodeId).
		Update("icon_id", p.IconId).
		Error
	return err
}
