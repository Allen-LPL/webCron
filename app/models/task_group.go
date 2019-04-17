package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type TaskGroup struct {
	Id          int
	UserId      int
	GroupName   string
	Description string
	CreateTime  int64
}

func (t *TaskGroup) TableName() string {
	return TableName("task_group")
}

func (t *TaskGroup) Update(fields ...string) error {
	if t.GroupName == "" {
		return fmt.Errorf("组名不能为空")
	}
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskGroupAdd(obj *TaskGroup) (int64, error) {
	if obj.GroupName == "" {
		return 0, fmt.Errorf("组名不能为空")
	}
	return orm.NewOrm().Insert(obj)
}

func TaskGroupGetById(id int) (*TaskGroup, error) {
	obj := &TaskGroup{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func TaskGroupDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task_group")).Filter("id", id).Delete()
	return err
}

func TaskGroupGetAll(page, pageSize int) ([]*TaskGroup, int64) {
	offset := (page - 1) * pageSize

	list := make([]*TaskGroup, 0)

	query := orm.NewOrm().QueryTable(TableName("task_group"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func TaskGroupGetList(page, pageSize int) ([]*TaskGroup, int64) {
	offset := (page - 1) * pageSize

	list := make([]*TaskGroup, 0)

	query := orm.NewOrm().QueryTable(TableName("task_group"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func RoleTaskGroupsGetByRoleId(role_id int) ([]*Resource) {
	//var list orm.ParamsList
	var list []*Resource
	//var list []*UserAndRole

	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb.Select( "t_resource.url",
		"t_resource.name",
		"t_resource.id").
		From("t_role_resource").
		InnerJoin("t_resource").On("t_resource.id = t_role_resource.resource_id").
		Where("t_role_resource.role_id = ?")

	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o := orm.NewOrm()
	o.Raw(sql, role_id).QueryRows(&list)

	return list
}