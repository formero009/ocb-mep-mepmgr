package manager

import (
	"mepmgr/dao"
	"mepmgr/models/mepmd"
)

type PerferenceManager interface {
	SavePerference(*mepmd.Preference) error
	GetPerference(userName string) (*mepmd.Preference, error)
}

var defaultPerferenceMg PerferenceMg

type PerferenceMg struct {
}

func NewDefaultPreferenceManager() PerferenceManager {
	mepTopologyOnce.Do(func() {
		defaultPerferenceMg = PerferenceMg{}
	})
	return defaultPerferenceMg
}

func (m PerferenceMg) SavePerference(p *mepmd.Preference) error {
	exist, err := dao.PreferenceDao.GetByUserName(p.UserName)
	if err != nil && err.Error() != "record not found" {
		return err
	}
	if exist.UserName == "" {
		err = dao.PreferenceDao.Create(p)
	} else {
		err = dao.PreferenceDao.Update(p)
	}
	return err
}

func (m PerferenceMg) GetPerference(userName string) (*mepmd.Preference, error) {
	exist, err := dao.PreferenceDao.GetByUserName(userName)
	if err != nil && err.Error() == "record not found" {
		return &mepmd.Preference{}, nil
	}
	return exist, err
}
