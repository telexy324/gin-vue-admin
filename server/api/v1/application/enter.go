package application

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CmdbServerApi
	CmdbSystemApi
	StaffApi
}

var cmdbServerService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbServerService
var cmdbSystemService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbSystemService
var staffService = service.ServiceGroupApp.ApplicationServiceGroup.StaffService
var taskService = service.ServiceGroupApp.ApplicationServiceGroup.TaskService