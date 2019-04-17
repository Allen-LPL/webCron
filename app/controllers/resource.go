package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"webcron/app/models"
	"webcron/app/libs"
)

type ResourceController struct {
	BaseController
}

func (this *ResourceController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	list, count := models.ResourceGetList(page, this.pageSize)

	this.Data["pageTitle"] = "用户列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ResourceController.List"), true).ToString()
	this.display()
}

func (this *ResourceController) Add() {
	if this.isPost() {
		resource := new(models.Resource)
		resource.Name = strings.TrimSpace(this.GetString("resource_name"))
		resource.Url = strings.TrimSpace(this.GetString("resource_url"))

		_, err := models.ResourceAdd(resource)
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "添加角色规则"
	this.display()
}

func (this *ResourceController) Edit() {
	id, _ := this.GetInt("id")

	resource, err := models.ResourceGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		resource.Url= strings.TrimSpace(this.GetString("resource_url"))
		resource.Name = strings.TrimSpace(this.GetString("resource_name"))

		err := resource.Update()
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "编辑角色规则"
	this.Data["resource"] = resource
	this.display()
}

func (this *ResourceController) Batch() {
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
			models.ResourceDelById(id)
			//models.TaskResetGroupId(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}
