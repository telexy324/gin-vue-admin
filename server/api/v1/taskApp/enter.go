package taskApp

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	TaskApi
	TemplateApi
	ScheduleApi
}

var taskService = service.ServiceGroupApp.TaskServiceGroup.TaskService
var templateService = service.ServiceGroupApp.TaskServiceGroup.TaskTemplatesService
var scheduleService = service.ServiceGroupApp.ScheduleServiceGroup.ScheduleService