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
var cmdbServerService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbServerService
var cmdbSystemService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbSystemService
var logUploadServerService = service.ServiceGroupApp.LogUploadServiceGroup.ServerService
var logUploadSecretService = service.ServiceGroupApp.LogUploadServiceGroup.SecretService
