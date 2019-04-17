package jobs

import (
	"github.com/astaxie/beego"
	"sync"
	"webcron/app/cron"

)

var (
	mainCron *cron.Cron
	workPool chan bool
	lock     sync.Mutex
)

func init() {
	if size, _ := beego.AppConfig.Int("jobs.pool"); size > 0 {
		workPool = make(chan bool, size)
	}
	mainCron = cron.New()
	mainCron.Start()
}

func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetEntryById(job.id) != nil {
		return false
	}
	err := mainCron.AddJob(spec, job)
	if err != nil {
		beego.Error("AddJob: ", err.Error())
		return false
	}
	return true
}

func RemoveJob(id int) {
	mainCron.RemoveJob(func(e *cron.Entry) bool {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

func GetEntryById(id int) *cron.Entry {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return e
			}
		}
	}
	return nil
}

func GetEntries(size int) []*cron.Entry {
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}
