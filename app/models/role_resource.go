package models

import (
	"github.com/astaxie/beego/orm"
)

type RoleResource struct {
	ResourceId        int
	Id 				int  `orm:"pk;column(role_id);"`

	//Roles []*Roles `orm:"reverse(many);"`  // 设置一对多的反向关系
	//Roles []*Roles `orm:"reverse(many)"`  // 设置一对多的反向关系
}

func (u *RoleResource) TableName() string {
	return TableName("role_resource")
}

func (u *RoleResource) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func RoleResourceAdd(user *RoleResource) (int64, error) {
	return orm.NewOrm().Insert(user)
}

//func RoleResourceAddAll(number int, fields ...string) (int64, error) {
//
//	num, err := dORM.InsertMulti(100, users)
//
//	num, err := orm.NewOrm().InsertMulti(number, fields)
//	if err != nil {
//		return 0, err
//	}
//	return num, nil
//}

func RoleResourceUrlGetByRoleId(role_id int) ([]*Resource) {
	//var list orm.ParamsList
	var list []*Resource
	//var list []*UserAndRole

	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb.Select("t_resource.url",
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

func RoleResourceGetByRoleId(role_id int) ([]interface{}, error) {
	var list orm.ParamsList
	_, err := orm.NewOrm().QueryTable(TableName("role_resource")).Filter("id", role_id).OrderBy("-id").ValuesFlat(&list, "ResourceId")

	return list, err
}

func RoleResourceUpdate(user *RoleResource, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func RoleResourceDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("role_resource")).Filter("id", id).Delete()
	return err
}

func RoleResourceGetList(page, pageSize int) ([]*RoleResource, int64) {
	offset := (page - 1) * pageSize

	list := make([]*RoleResource, 0)
	query := orm.NewOrm().QueryTable(TableName("role_resource"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func RoleResourceList(id int) ([]*RoleResource, error) {
	list := make([]*RoleResource, 0)
	query := orm.NewOrm().QueryTable(TableName("role_resource")).Filter("id", id)
	query.OrderBy("-id").All(&list)

	return list, nil
}