package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"reflect"
)

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

type Task struct {
	Id           int
	UserId       int
	GroupId      int
	TaskName     string
	TaskType     int
	Description  string
	CronSpec     string
	Concurrent   int
	Command      string
	Status       int
	Notify       int
	NotifyEmail  string
	Timeout      int
	ExecuteTimes int
	PrevTime     int64
	CreateTime   int64
	DieTime 	 int64
}

func (t *Task) TableName() string {
	return TableName("task")
}

func (t *Task) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskAdd(task *Task) (int64, error) {
	if task.TaskName == "" {
		return 0, fmt.Errorf("TaskName字段不能为空")
	}
	if task.CronSpec == "" {
		return 0, fmt.Errorf("CronSpec字段不能为空")
	}
	if task.Command == "" {
		return 0, fmt.Errorf("Command字段不能为空")
	}
	if task.CreateTime == 0 {
		task.CreateTime = time.Now().Unix()
	}
	return orm.NewOrm().Insert(task)
}

func TaskGetList(page, pageSize int, userId int, filters ...interface{}) ([]*Task, int64) {
	offset := (page - 1) * pageSize

	tasks := make([]*Task, 0)

	query := orm.NewOrm().QueryTable(TableName("task"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

	// 获取用户的分组
	if userId != 0 {
		groups := RoleTaskGroupGetByUserId(userId)
		if len(groups) != 0 {
			searchGroups := make([]int, 0)
			for _, v := range groups {
				searchGroups = append(searchGroups, v.Id)
			}
			query = query.Filter("group_id__in", searchGroups)
		}
	}

	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&tasks)

	return tasks, total
}

func TaskResetGroupId(groupId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task")).Filter("group_id", groupId).Update(orm.Params{
		"group_id": 0,
	})
}

func TaskGetById(id int) (*Task, error) {
	task := &Task{
		Id: id,
	}

	err := orm.NewOrm().Read(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func TaskDel(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task")).Filter("id", id).Delete()
	return err
}

func SliceColumn(structSlice []interface{}, key string) []interface{} {
	rt := reflect.TypeOf(structSlice)
	rv := reflect.ValueOf(structSlice)
	if rt.Kind() == reflect.Slice { //切片类型
		var sliceColumn []interface{}
		elemt := rt.Elem() //获取切片元素类型
		for i := 0; i < rv.Len(); i++ {
			inxv := rv.Index(i)
			if elemt.Kind() == reflect.Struct {
				for i := 0; i < elemt.NumField(); i++ {
					if elemt.Field(i).Name == key {
						strf := inxv.Field(i)
						switch strf.Kind() {
						case reflect.String:
							sliceColumn = append(sliceColumn, strf.String())
						case reflect.Float64:
							sliceColumn = append(sliceColumn, strf.Float())
						case reflect.Int, reflect.Int64:
							sliceColumn = append(sliceColumn, strf.Int())
						default:
							//do nothing
						}
					}
				}
			}
		}
		return sliceColumn
	}
	return nil
}
