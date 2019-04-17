package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Roles struct {
	Id          int
	UserId     int
	RoleName   string
	Description string
	CreateTime  int64

	//RoleResource []*RoleResource `orm:"rel(m2m);column(role_id)"`  // 设置一对多的反向关系
	//RoleResource *RoleResource `orm:"rel(fk);column(role_id)"`  // 设置一对多的反向关系
}

func (t *Roles) TableName() string {
	return TableName("roles")
}

func (t *Roles) Update(fields ...string) error {
	if t.RoleName == "" {
		return fmt.Errorf("角色名不能为空")
	}
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func RoleAdd(obj *Roles) (int64, error) {
	if obj.RoleName == "" {
		return 0, fmt.Errorf("角色名不能为空")
	}
	return orm.NewOrm().Insert(obj)
}

func RoleGetById(id int) (*Roles, error) {
	obj := &Roles{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func RoleNameGetById(id int) (*Roles, error) {
	//var role Roles
	u := new(Roles)
	err := orm.NewOrm().QueryTable(TableName("roles")).Filter("id", id).One(u, "roleName")
	if err != nil {
		return nil, err
	}
	return u, nil
}

func RoleDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("roles")).Filter("id", id).Delete()
	return err
}

func RoleGetList(page, pageSize int) ([]*Roles, int64) {
	offset := (page - 1) * pageSize

	list := make([]*Roles, 0)
	query := orm.NewOrm().QueryTable(TableName("roles")).RelatedSel()
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func RoleList() ([]*Roles, error) {
	list := make([]*Roles, 0)
	query := orm.NewOrm().QueryTable(TableName("roles"))
	query.OrderBy("-id").All(&list)

	return list, nil
}