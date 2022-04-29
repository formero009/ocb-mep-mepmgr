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

/*

查询mep_group 包含mep 的表
select re.mep_id, re.mep_group_id from mep_group as gr,mep_group_relation as re
where re.mep_group_id = gr.id;

select re.mep_id, mp.mep_name,re.mep_group_id,gr.mep_group_name
from mep_group as gr,mep_group_relation as re ,mep_meta as mp
where re.mep_group_id = gr.id and mp.mep_id = re.mep_id;

查询mep不关联mep_group 的表
方法一
select mp.mep_id,mp.mep_name from mep_meta mp left join mep_group_relation re
on mp.mep_id = re.mep_id where re.mep_id is null;

方法二
select mp.mep_id,mp.mep_name from mep_meta mp
where mp.mep_id not in (select re.mep_id from mep_group_relation as re);


搜索关键字，在mep_group表中
select  distinct  gr.id,gr.mep_group_name,re.mep_id ,mp.mep_name
 from mep_group gr
left join mep_group_relation re on gr.id = re.mep_group_id
left join mep_meta mp on re.mep_id = mp.mep_id
 where gr.mep_group_name LIKE '%g%';


搜索关键字，在mep_meta中
 select  distinct  gr.id,gr.mep_group_name,re.mep_id ,mp.mep_name
 from mep_meta mp
left join mep_group_relation re on mp.mep_id = re.mep_id
left join mep_group gr  on re.mep_group_id = gr.id
 where mp.mep_name LIKE '%m%';
*/

var TopologyDao topologyMeta

type topologyMeta struct{}

func (m topologyMeta) GetMepGroupRelation(authMepIds []string) (error, *[]mepmd.TopologyMepGroupRe) {
	sql_s := "select re.mep_id, mp.mep_name,gr.id,gr.mep_group_name " +
		"from mep_group as gr,mep_group_relation as re ,mep_meta as mp " +
		"where re.mep_group_id = gr.id and mp.mep_id = re.mep_id"

	if authMepIds != nil {
		sql_s = fmt.Sprintf("%s AND mp.mep_id IN ?", sql_s)
	}
	var data []mepmd.TopologyMepGroupRe

	err := models.PostgresDB.Raw(sql_s, authMepIds).Scan(&data).Error
	fmt.Println(err, data)
	return err, &data
}

func (m topologyMeta) GetMepNoRelationGroup(authMepIds []string) (error, *[]mepmd.TopologyMep) {
	sql_s := "select mp.mep_id,mp.mep_name from mep_meta mp " +
		"where mp.mep_id not in (select re.mep_id from mep_group_relation as re)"

	if authMepIds != nil {
		sql_s = fmt.Sprintf("%s AND mp.mep_id IN ?", sql_s)
	}
	var data []mepmd.TopologyMep
	err := models.PostgresDB.Raw(sql_s, authMepIds).Scan(&data).Error
	fmt.Println(err, data)
	return err, &data
}

func (m topologyMeta) GetGroupNoRelationMep() (error, *[]mepmd.TopologyMepGroupRe) {
	sql_s := "select id ,mep_group_name from mep_group " +
		"where id not in (select mep_group_id from mep_group_relation);"

	var data []mepmd.TopologyMepGroupRe
	err := models.PostgresDB.Raw(sql_s).Scan(&data).Error
	fmt.Println(err, data)
	return err, &data
}

func (m topologyMeta) SearchMepGroupbykey(key string, authMepIds []string) (error, *[]mepmd.TopologyMepGroupRe) {
	sql_s := "select gr.mep_group_name,mp.mep_id ,mp.mep_name ,gr.id " +
		"from mep_group gr " +
		"left join mep_group_relation re on gr.id = re.mep_group_id " +
		"left join mep_meta mp on re.mep_id = mp.mep_id " +
		"where gr.mep_group_name LIKE ?"

	var data []mepmd.TopologyMepGroupRe
	key = "%" + key + "%"
	if authMepIds != nil {
		sql_s = fmt.Sprintf("%s AND mp.mep_id IN ?", sql_s)
	}
	err := models.PostgresDB.Raw(sql_s, key, authMepIds).Scan(&data).Error
	fmt.Println(err, data)
	return err, &data
}

func (m topologyMeta) SearchMepbykey(key string, authMepIds []string) (error, *[]mepmd.TopologyMepGroupRe) {
	sql_s := " select  distinct  gr.id,gr.mep_group_name,mp.mep_id ,mp.mep_name " +
		"from mep_meta mp " +
		"left join mep_group_relation re on mp.mep_id = re.mep_id " +
		"left join mep_group gr on re.mep_group_id = gr.id " +
		"where mp.mep_name LIKE ?"

	var data []mepmd.TopologyMepGroupRe
	key = "%%" + key + "%%"
	if authMepIds != nil {
		fmt.Sprintf("%s and mp.mep_id IN ?", sql_s)
	}
	err := models.PostgresDB.Raw(sql_s, key).Scan(&data).Error
	fmt.Println(err, data)
	return err, &data
}
func (m topologyMeta) GetMepByGroupId(group_id int64, authMepIds []string) (error, *[]mepmd.MepMeta) {
	var err error
	var data []mepmd.MepMeta
	var sql_s string
	if group_id != 0 {
		sql_s = "select * from mep_meta mp where mp.mep_id in " +
			"(select mep_id from mep_group_relation re where re.mep_group_id = ?)"
	}else if group_id == 0 {
		sql_s = "select * from mep_meta mp where mp.mep_id not in (select mep_id from mep_group_relation)"
	}
	if authMepIds != nil {
		sql_s = fmt.Sprintf("%s AND mp.mep_id In ?", sql_s)
	}
	err = models.PostgresDB.Raw(sql_s, group_id, authMepIds).Scan(&data).Error
	return err, &data
}

func (m topologyMeta) GetMepByMepId(mep_id string, authMepIds []string) (error, *[]mepmd.MepMeta) {
	sql_s := "select * from mep_meta mp where mp.mep_id = ?"

	if authMepIds != nil {
		sql_s = fmt.Sprintf("%s AND mp.mep_id IN (?)", sql_s)
	}
	var data []mepmd.MepMeta
	err := models.PostgresDB.Raw(sql_s, mep_id, authMepIds).Scan(&data).Error
	fmt.Println("************", data)
	return err, &data
}

func (m topologyMeta) GetMepByLocation(provincial, city string, authMepIds []string) (error, *[]mepmd.MepMeta) {
	var data []mepmd.MepMeta
	tx := models.PostgresDB
	if len(provincial) != 0 {
		tx = tx.Where("province = ?", provincial)
	}
	if len(city) != 0 {
		tx = tx.Where("city = ?", city)
	}
	if authMepIds != nil {
		tx = tx.Where("mep_id IN (?)", authMepIds)
	}
	tx = tx.Order("province DESC")

	err := tx.Find(&data).Error

	return err, &data
}