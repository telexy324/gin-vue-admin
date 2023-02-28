package schedules

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ServiceGroup struct {
	ScheduleRunnerService
}

var scheduleService = service.ServiceGroupApp.ScheduleServiceGroup
