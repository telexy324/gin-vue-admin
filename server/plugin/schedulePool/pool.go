package schedulePool

import (
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/taskPool"
	"github.com/robfig/cron/v3"
	"sync"
)

var SPool SchedulePool

type ScheduleRunner struct {
	scheduleID int
	pool       *SchedulePool
}

func (r ScheduleRunner) Run() {
	schedule, err := scheduleService.GetSchedule(float64(r.scheduleID))
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}

	_, err = r.pool.taskPool.AddTask(taskMdl.Task{
		TemplateId:  schedule.TemplateID,
		CommandVars: schedule.CommandVars,
	}, 999999, 0)

	if err != nil {
		global.GVA_LOG.Error(err.Error())
	}
}

type SchedulePool struct {
	cron     *cron.Cron
	locker   sync.Locker
	taskPool *taskPool.TaskPool
}

func (p *SchedulePool) init() {
	p.cron = cron.New()
	p.locker = &sync.Mutex{}
}

func (p *SchedulePool) Refresh() {
	schedules, err := scheduleService.GetSchedules()

	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}

	p.locker.Lock()
	defer p.locker.Unlock()
	p.clear()
	for _, schedule := range schedules {
		if schedule.Valid != consts.ScheduleValid {
			continue
		}
		_, err = p.addRunner(ScheduleRunner{
			scheduleID: int(schedule.ID),
			pool:       p,
		}, schedule.CronFormat)
		if err != nil {
			global.GVA_LOG.Error(err.Error())
		}
	}
}

func (p *SchedulePool) addRunner(runner ScheduleRunner, cronFormat string) (int, error) {
	id, err := p.cron.AddJob(cronFormat, runner)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (p *SchedulePool) Run() {
	p.cron.Run()
}

func (p *SchedulePool) clear() {
	runners := p.cron.Entries()
	for _, r := range runners {
		p.cron.Remove(r.ID)
	}
}

func (p *SchedulePool) Destroy() {
	defer p.locker.Unlock()
	p.locker.Lock()
	p.cron.Stop()
	p.clear()
	p.cron = nil
}

func CreateSchedulePool(taskPool *taskPool.TaskPool) {
	SPool = SchedulePool{
		taskPool: taskPool,
	}
	SPool.init()
	if global.GVA_DB != nil {
		SPool.Refresh()
	}
	go SPool.Run()
}

func ValidateCronFormat(cronFormat string) error {
	_, err := cron.ParseStandard(cronFormat)
	return err
}
