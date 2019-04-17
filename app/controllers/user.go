package controllers

import (
	"strconv"
	"strings"
	"webcron/app/models"
	"webcron/app/libs"
	"github.com/astaxie/beego"
)

type UserController struct {
	BaseController
}

func (this *UserController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	list, count:= models.UserGetListSql(page, this.pageSize)

	this.Data["pageTitle"] = "用户列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("UserController.List"), true).ToString()
	this.display()
}

func (this *UserController) Add() {
	// 获取角色列表
	roleList, errRole := models.RoleList()
	if errRole != nil {
		this.showMsg(errRole.Error())
	}

	if this.isPost() {
		user := new(models.User)
		user.UserName = strings.TrimSpace(this.GetString("user_name"))
		user.Email = strings.TrimSpace(this.GetString("email"))
		user.Password = strings.TrimSpace(this.GetString("password"))
		user.Password = libs.Md5([]byte(user.Password+user.Salt))
		roleId, _ := this.GetInt("role_id")
		//roleInfo, _ := models.RoleGetById(roleId)
		//user.RoleName = roleInfo.RoleName
		// 写入用户信息数据
		id, err := models.UserAdd(user)
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		// 记录用户的角色ID
		userRole := new(models.UserRole)
		userRole.RoleId = roleId
		userRole.Id = int(id)

		//userRoleInfo, err := models.UserRoleGetByUserId(userRole.Id)
		//if userRoleInfo != nil {
			models.UserRoleAdd(userRole)
		//} else {
		//	models.UserRoleUpdate(userRole)
		//}

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["roleList"] = roleList
	this.Data["pageTitle"] = "添加用户"
	this.display()
}

func (this *UserController) Edit() {
	id, _ := this.GetInt("id")

	user, err := models.UserGetByIdOld(id)
	userRole, err := models.UserRoleGetByUserId(id)
	roleList, err := models.RoleList()
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		user.UserName = strings.TrimSpace(this.GetString("user_name"))
		user.Email = strings.TrimSpace(this.GetString("email"))
		user.Status, _ = this.GetInt("status")
		userRole.Id, _ = this.GetInt("id")
		roleId, _ := this.GetInt("role_id")
		userRole.RoleId = roleId

		// 修改用户表信息
		err := user.Update()
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		//errExec := user.UserExecForRoleId(roleId, id)
		//if errExec != nil {
		//	this.ajaxMsg(err.Error(), MSG_ERR)
		//}

		// 修改用户-角色表信息
		errRole := userRole.Update()
		if errRole != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "编辑用户"
	this.Data["user"] = user
	this.Data["userRole"] = userRole
	this.Data["roleList"] = roleList
	this.display()
}

func (this *UserController) Batch() {
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
			models.UserDelById(id)
			//models.TaskResetGroupId(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}
