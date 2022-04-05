package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CmdbServerApi
	CmdbSystemApi
	StaffApi
}

var environmentService = service.ServiceGroupApp.AnsibleServiceGroup.EnvironmentService
var inventoryService = service.ServiceGroupApp.AnsibleServiceGroup.InventoryService
var keyService = service.ServiceGroupApp.AnsibleServiceGroup.KeyService
var scheduleService = service.ServiceGroupApp.AnsibleServiceGroup.SchedulesService
var taskService = service.ServiceGroupApp.AnsibleServiceGroup.TaskService
var templateService = service.ServiceGroupApp.AnsibleServiceGroup.TemplatesService
var userService = service.ServiceGroupApp.AnsibleServiceGroup.UserService
