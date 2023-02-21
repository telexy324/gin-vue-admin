package task

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	TaskApi
}

var taskService = service.ServiceGroupApp.ApplicationServiceGroup.TaskService
