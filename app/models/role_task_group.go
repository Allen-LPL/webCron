package models

import (
	"github.com/astaxie/beego/orm"
)

type RoleTaskGroup struct {
	TaskGroupId        int
	Id 				int  `orm:"pk;column(role_id);"`
}

func (u *RoleTaskGroup) TableName() string {
	return TableName("role_task_group")
}

func (u *RoleTaskGroup) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func RoleTaskGroupAdd(user *RoleTaskGroup) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func RoleTaskGroupGetByUserId(user_id int) ([]*TaskGroup) {
	//var list orm.ParamsList
	var list []*TaskGroup
	//var list []*UserAndRole

	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb.Select("t_task_group.group_name",
		"t_task_group.id",
		"t_task_group.user_id").
		From("t_role_task_group").
		LeftJoin("t_roles").On("t_roles.id = t_role_task_group.role_id").
		LeftJoin("t_user_role").On("t_user_role.role_id = t_role_task_group.role_id").
		LeftJoin("t_task_group").On("t_role_task_group.task_group_id = t_task_group.id").
		Where("t_user_role.user_id = ?")

	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o := orm.NewOrm()
	o.Raw(sql, user_id).QueryRows(&list)

	return list
}

func RoleTaskGroupGetByRoleId(role_id int) ([]interface{}, error) {
	var list orm.ParamsList
	_, err := orm.NewOrm().QueryTable(TableName("role_task_group")).Filter("id", role_id).OrderBy("-id").ValuesFlat(&list, "TaskGroupId")

	return list, err
}

func RoleTaskGroupUpdate(user *RoleTaskGroup, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func RoleTaskGroupDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("role_task_group")).Filter("id", id).Delete()
	return err
}

func RoleTaskGroupGetList(page, pageSize int) ([]*RoleTaskGroup, int64) {
	offset := (page - 1) * pageSize

	list := make([]*RoleTaskGroup, 0)
	query := orm.NewOrm().QueryTable(TableName("role_resource"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func RoleTaskGroupList(id int) ([]*RoleTaskGroup, error) {
	list := make([]*RoleTaskGroup, 0)
	query := orm.NewOrm().QueryTable(TableName("role_resource")).Filter("id", id)
	query.OrderBy("-id").All(&list)

	return list, nil
}