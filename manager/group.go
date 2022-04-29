/*
@Time : 2021/3/15
@Author : jzd
@Project: mepmgr
*/
package manager

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"gorm.io/gorm"
	"mepmgr/common"
	"mepmgr/dao"
	"mepmgr/models"
	"mepmgr/models/mepmd"
	"strings"
	"sync"
)

type MepGroupManager interface {
	Create(mepGroup *mepmd.MepGroup) error
	Delete(mepGroupName string) error
	Update(oldName string, mepGroupUpdate mepmd.MepGroupUpdate) error
	List(from, size int, matchStr string, bean mepmd.MepGroup) ([]mepmd.MepGroup, int64, error)
	AddMep(mepGroupName string, meps mepmd.MepGroupAddMep) error
	DelMep(mepGroupName, mepName string) error
}

var mepGroupOnce sync.Once
var mepGroupMg defaultMepGroupMg

func NewDefaultMepGroupManager() MepGroupManager {
	mepGroupOnce.Do(func() {
		mepGroupMg = defaultMepGroupMg{}
	})
	return mepGroupMg
}

type defaultMepGroupMg struct{}

func (m defaultMepGroupMg) Create(mepGroup *mepmd.MepGroup) error {
	bean := &mepmd.MepGroup{MepGroupName: mepGroup.MepGroupName}
	return models.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if err := dao.MepGroupDao.Get(bean); err == nil && bean.Id != 0 {
			logs.Error("mep group %s already exist", mepGroup.MepGroupName)
			return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("分组 %s %s", mepGroup.MepGroupName, common.MsgAlreadyExist))
		}
		if err := dao.MepGroupDao.Create(tx, mepGroup); err != nil {
			return common.NewError(common.ErrDatabase)
		}
		return nil
	})
}

func (m defaultMepGroupMg) Delete(mepGroupName string) error {
	group := mepmd.MepGroup{MepGroupName: mepGroupName}
	if err := dao.MepGroupDao.Get(&group); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get mep group fail %s, %s not exist", err.Error(), mepGroupName)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("分组 %s %s", mepGroupName, common.MsgNotFound))
		}
		logs.Error("get mep group by bean [%+v] fail %s", group, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	id := mepmd.MepGroupRelation{MepGroupId: group.Id}

	err := dao.GetMepGroupRelation(&id)
	if err == nil {
		logs.Error("can not delete non-empty mep group %s", mepGroupName)
		return common.NewError(common.ErrCannotDel, fmt.Sprintf("分组%s下有Mep，%s", mepGroupName, common.MsgCannotDel))
	}

	err = dao.MepGroupDao.Delete(group)
	if err != nil {
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m defaultMepGroupMg) List(from, size int, matchStr string, bean mepmd.MepGroup) ([]mepmd.MepGroup, int64, error) {

	var metas []mepmd.MepGroup
	count := int64(0)

	if err := dao.MepGroupDao.Count(matchStr, bean, &count); err != nil {
		return nil, 0, common.NewError(common.ErrDatabase)
	}

	if err := dao.MepGroupDao.List(from, size, matchStr, bean, &metas); err != nil {
		return nil, 0, common.NewError(common.ErrDatabase)
	}

	return metas, count, nil
}

func (m defaultMepGroupMg) Update(oldName string, mepGroupUpdate mepmd.MepGroupUpdate) error {
	oldInfo := mepmd.MepGroup{MepGroupName: oldName}
	if err := dao.MepGroupDao.Get(&oldInfo); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get mep group fail %s, %s not exist", err.Error(), oldName)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("分组 %s %s", oldName, common.MsgNotFound))
		}
		logs.Error("get mep group by bean [%+v] fail %s", oldInfo, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	updateInfo := mepmd.MepGroup{
		CALLERId:     mepGroupUpdate.CALLERId,
		MEPMID:       mepGroupUpdate.MEPMID,
		MepGroupName: mepGroupUpdate.MepGroupName,
		Description:  mepGroupUpdate.Description,
	}
	if err := dao.MepGroupDao.Update([]string{"caller_id", "mepm_id", "mep_group_name", "description"}, mepmd.MepGroup{MepGroupName: oldName}, updateInfo); err != nil {
		logs.Error("Update mep group fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m defaultMepGroupMg) AddMep(mepGroupName string, meps mepmd.MepGroupAddMep) error {
	//check group exist
	group := mepmd.MepGroup{MepGroupName: mepGroupName}
	if err := dao.MepGroupDao.Get(&group); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get mep group fail %s, %s not exist", err.Error(), mepGroupName)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("分组 %s %s", mepGroupName, common.MsgNotFound))
		}
		logs.Error("get mep group by bean [%+v] fail %s", group, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	//check mep exist
	var mepInfos []mepmd.MepMeta
	err := dao.MepDao.GetMepInfosByMepNames(meps.MepNames, &mepInfos)
	if err != nil {
		logs.Error("get mep info fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	existMeps := make(map[string]bool)
	mepId2Name := make(map[string]string)
	checkRelations := make([]string, 0, 10)
	var relations []mepmd.MepGroupRelation
	for _, v := range mepInfos {
		existMeps[v.MepName] = true
		mepId2Name[v.MepId] = v.MepName
		relationInfo := mepmd.MepGroupRelation{
			MepId:      v.MepId,
			MepGroupId: group.Id,
		}
		relations = append(relations, relationInfo)
		checkRelations = append(checkRelations, fmt.Sprintf("('%v', %v)", relationInfo.MepId, relationInfo.MepGroupId))
	}
	for _, mep := range meps.MepNames {
		if !existMeps[mep] {
			logs.Error(fmt.Sprintf("mep %s not exist", mep))
			return common.NewError(common.ErrNotFound, fmt.Sprintf("Mep %s %s", mep, common.MsgNotFound))
		}
	}
	// check relation between target group and mep not exist
	var existedRelations []mepmd.MepGroupRelation
	checkRelationsStr := fmt.Sprintf("(%v)", strings.Join(checkRelations, " , "))
	err = dao.GetMepGroupRelations(checkRelationsStr, &existedRelations)
	if err == nil && len(existedRelations) > 0 {
		existedRelations2str := make([]string, 0, 10)
		for _, existRelation := range existedRelations {
			existedRelations2str = append(existedRelations2str, fmt.Sprintf("(%v, %v)", mepGroupName, mepId2Name[existRelation.MepId]))
		}
		logs.Error("relation between %v already exists", strings.Join(existedRelations2str, ", "))
		return common.NewError(common.ErrAlreadyExist, fmt.Sprintf("所属关系 %v %s", strings.Join(existedRelations2str, ", "), common.MsgAlreadyExist))
	}

	// create relation
	err = dao.CreateMepGroupRelation(&relations)
	if err != nil {
		logs.Error("create relation between group and mep fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	return nil
}

func (m defaultMepGroupMg) DelMep(mepGroupName, mepName string) error {
	group := mepmd.MepGroup{MepGroupName: mepGroupName}
	if err := dao.MepGroupDao.Get(&group); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get mep group fail %s, %s not exist", err.Error(), mepGroupName)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("分组 %s %s", mepGroupName, common.MsgNotFound))
		}
		logs.Error("get mep group by bean [%+v] fail %s", group, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	mep := mepmd.MepMeta{MepName: mepName}
	if err := dao.MepDao.Get(&mep); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get mep fail %s, %s not exist", err.Error(), mepName)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("Mep %s %s", mepName, common.MsgNotFound))
		}
		logs.Error("get mep by bean [%+v] fail %s", mep, err.Error())
		return common.NewError(common.ErrDatabase)
	}

	relation := mepmd.MepGroupRelation{
		MepId:      mep.MepId,
		MepGroupId: group.Id,
	}
	if err := dao.GetMepGroupRelation(&relation); err != nil {
		if err == dao.NOTFOUNDGET {
			logs.Error("get relation fail %s, relation between %s and %s not exist", err.Error(), relation.MepGroupId, relation.MepId)
			return common.NewError(common.ErrNotFound, fmt.Sprintf("所属关系 (%s, %s) %s", mepGroupName, mepName, common.MsgNotFound))
		}
		logs.Error("get relation by bean [%+v] fail %s", relation, err.Error())
		return common.NewError(common.ErrDatabase)
	}
	if err := dao.DeleteMepGroupRelation(relation); err != nil {
		logs.Error("delete relation info fail %s", err.Error())
		return common.NewError(common.ErrDatabase)
	}
	return nil
}
