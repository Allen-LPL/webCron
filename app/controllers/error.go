package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) NoPermission() {
	c.Data["content"] = "page not found"
	c.TplName = "error/noPermission.html"
}
