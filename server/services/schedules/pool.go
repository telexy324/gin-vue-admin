package schedules

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/services/tasks"
	"github.com/robfig/cron/v3"
	"sync"
)

var AnsibleSchedulePool *SchedulePool

type ScheduleRunner struct {
	projectID  int
	scheduleID int
	pool       *SchedulePool
}

func (r ScheduleRunner) Run() {
	schedule, err := scheduleService.GetSchedule(float64(r.projectID), float64(r.scheduleID))
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}

	_, err = r.pool.taskPool.AddTask(ansible.Task{
		TemplateID: schedule.TemplateID,
		ProjectID:  schedule.ProjectID,
	}, nil, schedule.ProjectID)

	if err != nil {
		global.GVA_LOG.Error(err.Error())
	}
}

type SchedulePool struct {
	cron     *cron.Cron
	locker   sync.Locker
	taskPool *tasks.TaskPool
}

func (p *SchedulePool) init() {
	p.cron = cron.New()
	p.locker = &sync.Mutex{}
}

func (p *SchedulePool) Refresh() {
	defer p.locker.Unlock()

	schedules, err := scheduleService.GetSchedules()

	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}

	p.locker.Lock()
	p.clear()
	for _, schedule := range schedules {
		_, err := p.addRunner(ScheduleRunner{
			projectID:  schedule.ProjectID,
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

func CreateSchedulePool(taskPool *tasks.TaskPool) SchedulePool {
	pool := SchedulePool{
		taskPool: taskPool,
	}
	pool.init()
	pool.Refresh()
	return pool
}

func ValidateCronFormat(cronFormat string) error {
	_, err := cron.ParseStandard(cronFormat)
	return err
}