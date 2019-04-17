package models

import (
	"github.com/astaxie/beego/orm"
)

type UserRole struct {
	Id 				int  `orm:"pk;column(user_id);"`
	RoleId        	int
}

func (u *UserRole) TableName() string {
	return TableName("user_role")
}

func (u *UserRole) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func UserRoleAdd(user *UserRole) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func UserRoleGetByUserId(UserId int) (*UserRole, error) {
	u := new(UserRole)

	err := orm.NewOrm().QueryTable(TableName("user_role")).Filter("user_id", UserId).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserRoleGetByRoleId(role_id int) (*UserRole, error) {
	u := new(UserRole)

	err := orm.NewOrm().QueryTable(TableName("user_role")).Filter("role_id", role_id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserRoleUpdate(user *UserRole, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func UserRoleDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("user_role")).Filter("id", id).Delete()
	return err
}

func UserRoleGetList(page, pageSize int) ([]*UserRole, int64) {
	offset := (page - 1) * pageSize

	list := make([]*UserRole, 0)
	query := orm.NewOrm().QueryTable(TableName("user_role"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}