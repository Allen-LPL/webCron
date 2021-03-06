package jobs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"io"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
	"webcron/app/mail"
	"webcron/app/models"
)

var mailTpl *template.Template

func init() {
	mailTpl, _ = template.New("mail_tpl").Parse(`
	你好 {{.username}}，<br/>

<p>以下是任务执行结果：</p>

<p>
任务ID：{{.task_id}}<br/>
任务名称：{{.task_name}}<br/>       
执行时间：{{.start_time}}<br />
执行耗时：{{.process_time}}秒<br />
执行状态：{{.status}}
</p>
<p>-------------以下是任务执行输出-------------</p>
<p>{{.output}}</p>
<p>
--------------------------------------------<br />
本邮件由系统自动发出，请勿回复<br />
如果要取消邮件通知，请登录到系统进行设置<br />
</p>
`)

}

type Job struct {
	id         int                                               // 任务ID
	logId      int64                                             // 日志记录ID
	name       string                                            // 任务名称
	task       *models.Task                                      // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	Concurrent bool                                              // 同一个任务是否允许并行执行
}

func NewJobFromTask(task *models.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}

	if task.DieTime > 0 && task.DieTime < time.Now().Unix() {
		return nil, fmt.Errorf("进程终止！")
	}

	job := NewCommandJob(task.Id, task.TaskName, task.Command)
	job.task = task
	job.Concurrent = task.Concurrent == 1
	return job, nil
}

func NewCommandJob(id int, name string, command string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {

		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		//is := exec.Command("/bin/bash", "-c", "ps aux |grep 'ping 103.94.185.192'")
		//is.Start()
		//fmt.Println("Exit Code", is.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())

		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	return job
}

func (j *Job) Status() int {
	return j.status
}

func (j *Job) GetName() string {
	return j.name
}

func (j *Job) GetId() int {
	return j.id
}

func (j *Job) GetLogId() int64 {
	return j.logId
}

func (j *Job) Run() {
	//beego.Error("发送邮件超时：", j.id)
	if j.task.DieTime > 0 && j.task.DieTime < time.Now().Unix() {
		RemoveJob(j.id)
		if task, err := models.TaskGetById(j.id); err == nil {
			task.Status = 0
			task.Update()
		}
		beego.Warn(fmt.Sprintf("任务[%d]进程中止！", j.id))
		return
	}

	if !j.Concurrent && j.status > 0 {
		println(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id))
		beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			beego.Error(err, "\n", string(debug.Stack()))
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	beego.Debug(fmt.Sprintf("开始执行任务: %d", j.id))

	j.status++
	defer func() {
		j.status--
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	// 执行命令之前先审查是否有相同的命令正在执行
	errBool := isHaving(j.task.Command)
	println(fmt.Sprintf("isHaving: %t", errBool))
	if !errBool {
		println(fmt.Sprintf("任务[%d]上一次执行的进程尚未结束，本次被忽略。", j.id))
		beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id))
		return
	}

	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)
	//if err == nil && cmdOut == "" && cmdErr == "" {
	//	println(fmt.Sprintf("任务[%d]上一次执行的进程尚未结束，本次被忽略。", j.id))
	//	beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id))
	//	return
	//}
	//实时循环读取输出流中的一行内容

	title := strings.Split(cmdOut, "\n")
	//println(fmt.Sprintf("cmdOUt: %s", title[0]))
	//reader := bufio.NewReader(cmdOut)
	//for {
	//	line, err2 := reader.ReadString('\n')
	//	if err2 != nil || io.EOF == err2 {
	//		break
	//	}
	//	fmt.Println(line)
	//}

	ut := time.Now().Sub(t) / time.Millisecond

	// 插入日志
	log := new(models.TaskLog)
	log.TaskId = j.id
	log.Output = cmdOut
	log.Error = cmdErr
	log.ProcessTime = int(ut)
	log.CreateTime = t.Unix()
	log.Title = title[0]

	if isTimeout {
		log.Status = models.TASK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), cmdErr)
	} else if err != nil {
		log.Status = models.TASK_ERROR
		log.Error = err.Error() + ":" + cmdErr
	}
	j.logId, _ = models.TaskLogAdd(log)

	// 更新上次执行时间
	j.task.PrevTime = t.Unix()
	j.task.ExecuteTimes++
	j.task.Update("PrevTime", "ExecuteTimes")

	// 发送邮件通知
	if (j.task.Notify == 1 && err != nil) || j.task.Notify == 2 {
		user, uerr := models.UserGetByIdOld(j.task.UserId)
		if uerr != nil {
			beego.Warn(fmt.Sprintf("任务[%d]用户信息获取失败！", j.id))
			return
		}

		var title string

		data := make(map[string]interface{})
		data["task_id"] = j.task.Id
		data["username"] = user.UserName
		data["task_name"] = j.task.TaskName
		data["start_time"] = beego.Date(t, "Y-m-d H:i:s")
		data["process_time"] = float64(ut) / 1000
		data["output"] = cmdOut

		if isTimeout {
			title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "超时")
			data["status"] = fmt.Sprintf("超时（%d秒）", int(timeout/time.Second))
		} else if err != nil {
			title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "失败")
			data["status"] = "失败（" + err.Error() + "）"
		} else {
			title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "成功")
			data["status"] = "成功"
		}

		content := new(bytes.Buffer)
		mailTpl.Execute(content, data)
		ccList := make([]string, 0)
		if j.task.NotifyEmail != "" {
			ccList = strings.Split(j.task.NotifyEmail, "\n")
		}
		if !mail.SendMail(user.Email, user.UserName, title, content.String(), ccList) {
			beego.Error("发送邮件超时：", user.Email)
		}
	}

}

func isHaving(command string) bool {
	// 判断是否正在执行
	isHasCmd := exec.Command("/bin/sh", "-c", `ps -ef | grep -v "grep" | grep "`+command+`"`)
	fmt.Println(`ps -ef | grep -v "grep" | grep "`+command+`"`)
	out, err := isHasCmd.StdoutPipe()

	if err != nil {
		fmt.Println("StdoutPipe: " + err.Error())
		return false
	}

	_, err = isHasCmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		return false
	}

	if err := isHasCmd.Start(); err != nil {
		fmt.Println("Start: ", err.Error())
		return false
	}

	if err := isHasCmd.Wait(); err == nil {
		fmt.Println("Wait: 已执行!")
		return false
	}

	reader := bufio.NewReader(out)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}


	return true
}
