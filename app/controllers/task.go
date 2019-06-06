package controllers

import (
	"strconv"
	"strings"
	"time"
	libcron "webcron/app/cron"
	"github.com/astaxie/beego"
	"webcron/app/models"
	"webcron/app/jobs"
	"webcron/app/libs"
	"net/http"
	"github.com/360EntSecGroup-Skylar/excelize"
	"fmt"
	"os"
	"io"
)

type TaskController struct {
	BaseController
}

// 任务列表
func (this *TaskController) List() {

	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	groupId, _ := this.GetInt("groupid")
	filters := make([]interface{}, 0)
	taskName := strings.TrimSpace(this.GetString("task_name"))
	if taskName != "" {
		filters = append(filters, "task_name__icontains", taskName)
	}

	if groupId > 0 {
		filters = append(filters, "group_id", groupId)
	}

	result, count := models.TaskGetList(page, this.pageSize, this.userId, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.TaskName
		row["cron_spec"] = v.CronSpec
		row["status"] = v.Status
		row["description"] = v.Description

		if v.DieTime > 0 {
			row["die_time"] = beego.Date(time.Unix(v.DieTime, 0), "Y-m-d H:i:s")
		} else {
			row["die_time"] = "_"
		}

		e := jobs.GetEntryById(v.Id)
		if e != nil {
			row["next_time"] = beego.Date(e.Next, "Y-m-d H:i:s")
			row["prev_time"] = "-"
			if e.Prev.Unix() > 0 {
				row["prev_time"] = beego.Date(e.Prev, "Y-m-d H:i:s")
			} else if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			}
			row["running"] = 1
		} else {
			row["next_time"] = "-"
			if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			} else {
				row["prev_time"] = "-"
			}
			row["running"] = 0
		}
		list[k] = row
	}

	excel1, _ := this.GetInt("excel")
	if excel1 == 1 {
		this.WriteFile(result)
		this.WriteExcel(result)
	}

	// 分组列表
	//groups := models.RoleTaskGroupGetByUserId(this.userId)
	var groups interface{}
	if this.userId != 1 {
		groups = models.RoleTaskGroupGetByUserId(this.userId)
	} else {
		groups, _ = models.TaskGroupGetList(1, 100)
	}

	this.Data["pageTitle"] = "任务列表"

	this.Data["list"] = list
	this.Data["task_name"] = taskName
	this.Data["groups"] = groups
	this.Data["groupid"] = groupId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List", "groupid", groupId), true).ToString()
	this.display()
}

// 添加任务
func (this *TaskController) Add() {

	if this.isPost() {
		task := new(models.Task)
		task.UserId = this.userId
		task.GroupId, _ = this.GetInt("group_id")
		task.TaskName = strings.TrimSpace(this.GetString("task_name"))
		task.Description = strings.TrimSpace(this.GetString("description"))
		task.Concurrent, _ = this.GetInt("concurrent")
		task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
		task.Command = strings.TrimSpace(this.GetString("command"))
		task.Notify, _ = this.GetInt("notify")
		task.Timeout, _ = this.GetInt("timeout")

		toBeCharge := strings.TrimSpace(this.GetString("die_time")) //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
		if len(toBeCharge) > 0 {
			timeLayout := "2006-01-02T15:04"                                //转化所需模板
			loc, _ := time.LoadLocation("Local")                            //重要：获取时区
			theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
			task.DieTime = theTime.Unix()
		}                                                           //转化为时间戳 类型是int64

		notifyEmail := strings.TrimSpace(this.GetString("notify_email"))
		if notifyEmail != "" {
			emailList := make([]string, 0)
			tmp := strings.Split(notifyEmail, "\n")
			for _, v := range tmp {
				v = strings.TrimSpace(v)
				if !libs.IsEmail([]byte(v)) {
					this.ajaxMsg("无效的Email地址："+v, MSG_ERR)
				} else {
					emailList = append(emailList, v)
				}
			}
			task.NotifyEmail = strings.Join(emailList, "\n")
		}

		if task.TaskName == "" || task.CronSpec == "" || task.Command == "" {
			this.ajaxMsg("请填写完整信息", MSG_ERR)
		}
		if _, err := libcron.Parse(task.CronSpec); err != nil {
			this.ajaxMsg("cron表达式无效", MSG_ERR)
		}
		if _, err := models.TaskAdd(task); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
	}

	// 分组列表
	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["pageTitle"] = "添加任务"
	this.display()
}

// 编辑任务
func (this *TaskController) Edit() {
	id, _ := this.GetInt("id")
	groupid, _ := this.GetInt("groupid")
	taskName := strings.TrimSpace(this.GetString("task_name"))

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		task.TaskName = strings.TrimSpace(this.GetString("taskName"))
		task.Description = strings.TrimSpace(this.GetString("description"))
		task.GroupId, _ = this.GetInt("group_id")
		task.Concurrent, _ = this.GetInt("concurrent")
		task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
		task.Command = strings.TrimSpace(this.GetString("command"))
		task.Notify, _ = this.GetInt("notify")
		task.Timeout, _ = this.GetInt("timeout")

		//loc, _ := time.LoadLocation("Local")
		//tm, _ := time.ParseInLocation("2018-01-01", "2018-06-06", time.Local)
		//task.DieTime = tm.Unix()

		toBeCharge := strings.TrimSpace(this.GetString("die_time")) //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
		if len(toBeCharge) > 0 {
			timeLayout := "2006-01-02T15:04"                                //转化所需模板
			loc, _ := time.LoadLocation("Local")                            //重要：获取时区
			theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
			task.DieTime = theTime.Unix()
		} else {
			task.DieTime = 0
		}

		notifyEmail := strings.TrimSpace(this.GetString("notify_email"))
		if notifyEmail != "" {
			tmp := strings.Split(notifyEmail, "\n")
			emailList := make([]string, 0, len(tmp))
			for _, v := range tmp {
				v = strings.TrimSpace(v)
				if !libs.IsEmail([]byte(v)) {
					this.ajaxMsg("无效的Email地址："+v, MSG_ERR)
				} else {
					emailList = append(emailList, v)
				}
			}
			task.NotifyEmail = strings.Join(emailList, "\n")
		} else {
			task.NotifyEmail = ""
		}

		if task.TaskName == "" || task.CronSpec == "" || task.Command == "" {
			this.ajaxMsg("请填写完整信息", MSG_ERR)
		}
		if _, err := libcron.Parse(task.CronSpec); err != nil {
			this.ajaxMsg("cron表达式无效", MSG_ERR)
		}
		if err := task.Update(); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
		//this.redirect(beego.URLFor("TaskController.List", "groupid", task.GroupId))
	}

	// 分组列表
	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groupid"] = groupid
	this.Data["task_name"] = taskName
	this.Data["groups"] = groups
	this.Data["task"] = task
	this.Data["pageTitle"] = "编辑任务"
	this.display()
}

// 任务执行日志列表
func (this *TaskController) Logs() {
	taskId, _ := this.GetInt("id")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	task, err := models.TaskGetById(taskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	result, count := models.TaskLogGetList(page, this.pageSize, "task_id", task.Id)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["status"] = v.Status
		row["title"] = v.Title
		list[k] = row
	}

	this.Data["pageTitle"] = "任务执行日志"
	this.Data["list"] = list
	this.Data["task"] = task
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.Logs", "id", taskId), true).ToString()
	this.display()
}

// 查看日志详情
func (this *TaskController) ViewLog() {
	id, _ := this.GetInt("id")

	taskLog, err := models.TaskLogGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	task, err := models.TaskGetById(taskLog.TaskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	data := make(map[string]interface{})
	data["id"] = taskLog.Id
	data["output"] = taskLog.Output
	data["error"] = taskLog.Error
	data["start_time"] = beego.Date(time.Unix(taskLog.CreateTime, 0), "Y-m-d H:i:s")
	data["process_time"] = float64(taskLog.ProcessTime) / 1000
	data["ouput_size"] = libs.SizeFormat(float64(len(taskLog.Output)))
	data["status"] = taskLog.Status

	this.Data["task"] = task
	this.Data["data"] = data
	this.Data["pageTitle"] = "查看日志"
	this.display()
}

// 批量操作日志
func (this *TaskController) LogBatch() {
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
			models.TaskLogDelById(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 批量操作
func (this *TaskController) Batch() {
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
		case "active":
			if task, err := models.TaskGetById(id); err == nil {
				job, err := jobs.NewJobFromTask(task)
				if err == nil {
					jobs.AddJob(task.CronSpec, job)
					task.Status = 1
					task.Update()
				}
			}
		case "pause":
			jobs.RemoveJob(id)
			if task, err := models.TaskGetById(id); err == nil {
				task.Status = 0
				task.Update()
			}
		case "delete":
			models.TaskDel(id)
			models.TaskLogDelByTaskId(id)
			jobs.RemoveJob(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 启动任务
func (this *TaskController) Start() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	job, err := jobs.NewJobFromTask(task)
	if err != nil {
		this.showMsg(err.Error())
	}

	if jobs.AddJob(task.CronSpec, job) {
		task.Status = 1
		task.Update()
	}

	refer := this.Ctx.Request.Referer()
	if refer == "" {
		refer = beego.URLFor("TaskController.List")
	}

	this.redirect(refer)
}

// 暂停任务
func (this *TaskController) Pause() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	jobs.RemoveJob(id)
	task.Status = 0
	task.Update()

	refer := this.Ctx.Request.Referer()
	if refer == "" {
		refer = beego.URLFor("TaskController.List")
	}
	this.redirect(refer)
}

// 立即执行
func (this *TaskController) Run() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	job, err := jobs.NewJobFromTask(task)
	if err != nil {
		this.showMsg(err.Error())
	}

	job.Run()

	this.redirect(beego.URLFor("TaskController.ViewLog", "id", job.GetLogId()))
}

//func ConvertToString(src string, srcCode string, tagCode string) string {
//    srcCoder := mahonia.NewDecoder(srcCode)
//    srcResult := srcCoder.ConvertString(src)
//    tagCoder := mahonia.NewDecoder(tagCode)
//    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
//    result := string(cdata)
//    return result
//}

// 导出文本模式
func (this *TaskController) WriteFile(list []*models.Task) {
	fileName := "./crontab.txt"
	var f *os.File
	var err1 error
	var info string
	if checkFileIsExist(fileName) {
		//f, err1 = os.OpenFile(fileName, os.O_APPEND, 0666)
		err1 = os.Remove(fileName)
	}

	f, err1 = os.Create(fileName)
	if err1 != nil {
		panic(err1)
	}

	for _, v := range list {
		if v.DieTime < time.Now().Unix() {
			continue
		}

		var cronInfo string
		v.CronSpec = strings.Replace(v.CronSpec, "?", "*", 1)
		ccList := strings.Split(v.CronSpec, " ")
		for k, j := range ccList {
			if k == 0 {
				continue
			}

			cronInfo = cronInfo + " " + j + " "
		}
		//v.TaskName = ConvertToString(v.TaskName, "gbk", "utf8")
		info = cronInfo + " " + v.Command + " // " + v.TaskName + "\r\n"

		println(info)

		_, err1 = io.WriteString(f, info)
		if err1 != nil {
			panic(err1)
		}
	}

	//f, err1 = os.Create(fileName)
	// 从服务器读取
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "crontable")) //文件名
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, fileName)
	//this.jsonResult(list)
}

// 查询文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func (this *TaskController) WriteExcel(list []*models.Task) {
	//this.jsonResult(list)
	xlsx := excelize.NewFile()
	// Create a new sheet.
	index := xlsx.NewSheet("Sheet1")

	//    style, err := xlsx.NewStyle(`{"border":[{"type":"left","color":"0000FF","style":3},{"type":"top","color":"00FF00","style":4},{"type":"bottom","color":"FFFF00","style":5},{"type":"right","color":"FF0000","style":6},{"type":"diagonalDown","color":"A020F0","style":7},{"type":"diagonalUp","color":"A020F0","style":8}]}`)
	//    if err != nil {
	//        fmt.Println(err)
	//    }
	//    xlsx.SetCellStyle("Sheet1", "H9", "H9", style)
	//xlsx.SetCellStyle("Sheet1", "A1", "Id", style)
	xlsx.SetCellValue("Sheet1", "A1", "Id")
	xlsx.SetCellValue("Sheet1", "B1", "UserId")
	xlsx.SetCellValue("Sheet1", "C1", "GroupId")
	xlsx.SetCellValue("Sheet1", "D1", "TaskName")
	xlsx.SetCellValue("Sheet1", "E1", "TaskType")
	xlsx.SetCellValue("Sheet1", "F1", "Description")
	xlsx.SetCellValue("Sheet1", "G1", "CronSpec")
	xlsx.SetCellValue("Sheet1", "H1", "Concurrent")
	xlsx.SetCellValue("Sheet1", "I1", "Command")
	xlsx.SetCellValue("Sheet1", "J1", "Status")
	xlsx.SetCellValue("Sheet1", "K1", "Notify")
	xlsx.SetCellValue("Sheet1", "L1", "NotifyEmail")
	xlsx.SetCellValue("Sheet1", "M1", "Timeout")
	xlsx.SetCellValue("Sheet1", "N1", "ExecuteTimes")
	xlsx.SetCellValue("Sheet1", "O1", "PrevTime")
	xlsx.SetCellValue("Sheet1", "P1", "CreateTime")
	xlsx.SetCellValue("Sheet1", "Q1", "DieTime")
	for k, v := range list {
		num := k + 2
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", num), v.Id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", num), v.UserId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", num), v.GroupId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", num), v.TaskName)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", num), v.TaskType)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", num), v.Description)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", num), v.CronSpec)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", num), v.Concurrent)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", num), v.Command)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("J%d", num), v.Status)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("K%d", num), v.Notify)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("L%d", num), v.NotifyEmail)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("M%d", num), v.Timeout)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("N%d", num), v.ExecuteTimes)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("O%d", num), v.PrevTime)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("P%d", num), v.CreateTime)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("Q%d", num), v.DieTime)
	}
	// Set value of a cell.
	//xlsx.SetCellValue("Sheet1", "A1", "Hello world.")
	//xlsx.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	errSave := xlsx.SaveAs("./listData.xlsx")
	if errSave != nil {
		fmt.Println(errSave)
	}

	// 从服务器读取
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "file.xlsx")) //文件名
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, "./listData.xlsx")
}
