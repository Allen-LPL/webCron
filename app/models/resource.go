package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Resource struct {
	Id          int
	Name   string
	Url string
}

func (t *Resource) TableName() string {
	return TableName("resource")
}

func (t *Resource) Update(fields ...string) error {
	if t.Id == 0 {
		return fmt.Errorf("角色规则ID不能为空")
	}
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func ResourceAdd(obj *Resource) (int64, error) {
	//if obj.Id == 0 {
	//	return 0, fmt.Errorf("角色规则ID不能为空")
	//}
	return orm.NewOrm().Insert(obj)
}

func ResourceGetById(id int) (*Resource, error) {
	obj := &Resource{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func ResourceDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("resource")).Filter("id", id).Delete()
	return err
}

func ResourceGetList(page, pageSize int) ([]*Resource, int64) {
	offset := (page - 1) * pageSize

	list := make([]*Resource, 0)
	query := orm.NewOrm().QueryTable(TableName("resource"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func ResourceList() ([]*Resource, error) {
	list := make([]*Resource, 0)
	query := orm.NewOrm().QueryTable(TableName("resource"))
	query.OrderBy("-id").All(&list)

	return list, nil
}