package models

import (
	"github.com/astaxie/beego/orm"
)

type UserAndRole struct {
	Id        int
	UserName  string
	Password  string
	Salt      string
	Email     string
	LastLogin int64
	LastIp    string
	Status    int

	RoleName string
}

type User struct {
	Id        int
	UserName  string
	Password  string
	Salt      string
	Email     string
	LastLogin int64
	LastIp    string
	Status    int
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *User) UserExecForRoleId(roleId int, id int) error {
	if _, err := orm.NewOrm().Raw("UPDATE t_user SET role_id = ? where id = ?", roleId, id).Exec(); err != nil {
		return err
	}
	return nil
}

func UserAdd(user *User) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func UserGetById(id int) (*UserAndRole, error) {
	var user *UserAndRole

	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb.Select("t_user.id",
		"t_user.user_name",
		"t_user.email",
		"t_user.password",
		"t_user.last_login",
		"t_user.last_ip",
		"t_user.status",
		"t_roles.role_name").
		From("t_user").
		InnerJoin("t_user_role").On("t_user.id = t_user_role.user_id").
		InnerJoin("t_roles").On("t_user_role.role_id = t_roles.id").
		Where("t_user.id = ?").
		OrderBy("t_user.id").Desc()

	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o := orm.NewOrm()
	err := o.Raw(sql, id).QueryRow(&user)

	return user, err
}

func UserGetByIdOld(id int) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).RelatedSel().One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserGetByName(userName string) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func UserDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).Delete()
	return err
}

func UserGetList(page, pageSize int) ([]*User, int64) {
	offset := (page - 1) * pageSize

	list := make([]*User, 0)
	query := orm.NewOrm().QueryTable(TableName("user"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func UserGetListSql(page, pageSize int) ([]*UserAndRole, int64) {
	var user []*UserAndRole

	offset := (page - 1) * pageSize
	query := orm.NewOrm().QueryTable(TableName("user"))
	total, _ := query.Count()

	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb.Select("t_user.id",
		"t_user.user_name",
		"t_user.email",
		"t_user.password",
		"t_user.last_login",
		"t_user.last_ip",
		"t_user.status",
		"t_roles.role_name").
		From("t_user").
		InnerJoin("t_user_role").On("t_user.id = t_user_role.user_id").
		InnerJoin("t_roles").On("t_user_role.role_id = t_roles.id").
		Where("1 = ?").
		OrderBy("t_user.id").Desc().
		Limit(pageSize).Offset(offset)

	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o := orm.NewOrm()
	o.Raw(sql, 1).QueryRows(&user)

	return user, total
}
