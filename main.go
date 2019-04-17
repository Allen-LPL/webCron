package main

import (
	"html/template"
	"net/http"
	"github.com/astaxie/beego"
	"webcron/app/models"
	"webcron/app/jobs"
	"webcron/app/controllers"
	"time"
)

const VERSION = "1.0.0"

// 模版中转换时间
func convertT(in int64) (out string) {
	tm := time.Unix(in, 0)
	out = tm.Format("2006-01-02T15:04")
	return
}

func main() {
	models.Init()
	jobs.InitJobs()

	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})

	// 生产环境不输出debug日志
	if beego.AppConfig.String("runmode") == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}
	beego.SetLevel(beego.LevelEmergency)

	beego.AppConfig.Set("version", VERSION)

	// 路由设置
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/gettime", &controllers.MainController{}, "*:GetTime")
	beego.Router("/help", &controllers.HelpController{}, "*:Index")
	beego.Router("/noPermission", &controllers.ErrorController{}, "*:NoPermission")

	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.GroupController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.ResourceController{})
	beego.AutoRouter(&controllers.UserController{})

	//beego.ErrorController(&controllers.ErrorController{})

	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.AddFuncMap("convertt", convertT)
	beego.Run()
}
