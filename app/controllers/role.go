package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"webcron/app/models"
	"webcron/app/libs"
	//"log"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	list, count := models.RoleGetList(page, this.pageSize)

	this.Data["pageTitle"] = "角色列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("RoleController.List"), true).ToString()
	this.display()
}

func (this *RoleController) Add() {
	// 角色规则列表
	resourceList, errResource := models.ResourceList()
	if errResource != nil {
		this.showMsg(errResource.Error())
	}

	if this.isPost() {
		role := new(models.Roles)
		role.RoleName = strings.TrimSpace(this.GetString("role_name"))
		role.UserId = this.userId
		role.Description = strings.TrimSpace(this.GetString("description"))
		id, err := models.RoleAdd(role)
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		// 写入角色的角色规则记录
		resourceModel := new(models.RoleResource)
		resourceIds := this.GetStrings("resourceIds")
		for _, v := range resourceIds {
			resourceModel.ResourceId, err = strconv.Atoi(v)
			resourceModel.Id = int(id)

			_, err := models.RoleResourceAdd(resourceModel)
			if err != nil {
				this.ajaxMsg(err.Error(), MSG_ERR)
			}
		}

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["resourceList"] = resourceList
	this.Data["pageTitle"] = "添加角色"
	this.display()
}

func (this *RoleController) Edit() {
	id, _ := this.GetInt("id")

	// 角色信息数据
	role, err := models.RoleGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	// 角色规则列表
	resourceList, errResource := models.ResourceList()
	if errResource != nil {
		this.showMsg(errResource.Error())
	}

	// 角色 - 角色规则
	roleResource, errRoleResource := models.RoleResourceGetByRoleId(id)
	if errRoleResource != nil {
		this.showMsg(errRoleResource.Error())
	}

	// 分组列表
	groups, _ := models.TaskGroupGetAll(1, 100)

	// 角色 - 分组列表
	roleTaskGroup, errRoleTaskGroup := models.RoleTaskGroupGetByRoleId(id)
	if errRoleTaskGroup != nil {
		this.showMsg(errRoleTaskGroup.Error())
	}

	//this.jsonResult(roleTaskGroup)

	//var out []string
	//
	//for _, i := range resourceList {
	//	for _, o := range roleResource {
	//		if i.Id == o {
	//			i["isCon"] = true
	//		} else {
	//			i["isCon"] = false
	//		}
	//		out[] = i
	//
	//	}
	//}

	if this.isPost() {
		// 修改角色信息
		role.RoleName = strings.TrimSpace(this.GetString("role_name"))
		role.Description = strings.TrimSpace(this.GetString("description"))
		err := role.Update()
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		/*{{{*/
		// 修改角色-角色规则信息
		models.RoleResourceDelById(id)	// 删除角色的规则信息

		// 写入角色的角色规则记录
		resourceModel := new(models.RoleResource)
		resourceIds := this.GetStrings("resourceIds")
		for _, v := range resourceIds {
			resourceModel.ResourceId, err = strconv.Atoi(v)
			resourceModel.Id = id

			_, err := models.RoleResourceAdd(resourceModel)
			if err != nil {
				this.ajaxMsg(err.Error(), MSG_ERR)
			}
		}
		/*}}}*/

		/*{{{*/
		// 修改角色-角色规则信息
		models.RoleTaskGroupDelById(id)	// 删除角色的规则信息

		// 写入角色的角色规则记录
		roleTaskGroupModel := new(models.RoleTaskGroup)
		roleTaskGroupIds := this.GetStrings("roleTaskGroupIds")
		for _, v := range roleTaskGroupIds {
			roleTaskGroupModel.TaskGroupId, err = strconv.Atoi(v)
			roleTaskGroupModel.Id = id

			_, err := models.RoleTaskGroupAdd(roleTaskGroupModel)
			if err != nil {
				this.ajaxMsg(err.Error(), MSG_ERR)
			}
		}
		/*}}}*/

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "编辑角色"
	this.Data["role"] = role
	this.Data["resourceList"] = resourceList
	this.Data["roleResource"] = roleResource
	this.Data["groups"] = groups
	this.Data["roleTaskGroup"] = roleTaskGroup
	//this.jsonResult(groups)
	//log.Print(groups)

	beego.Info()
	beego.Debug()

	this.display()
}

func (this *RoleController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			models.RoleDelById(id)
			//models.TaskResetGroupId(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}
