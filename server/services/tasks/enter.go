package tasks

import "github.com/flipped-aurora/gin-vue-admin/server/service"

var environmentService = service.ServiceGroupApp.AnsibleServiceGroup.EnvironmentService
var inventoryService = service.ServiceGroupApp.AnsibleServiceGroup.InventoryService
var keyService = service.ServiceGroupApp.AnsibleServiceGroup.KeyService
var projectService = service.ServiceGroupApp.AnsibleServiceGroup.ProjectService
var scheduleService = service.ServiceGroupApp.AnsibleServiceGroup.SchedulesService
var taskService = service.ServiceGroupApp.AnsibleServiceGroup.TaskService
var templateService = service.ServiceGroupApp.AnsibleServiceGroup.TemplatesService
var userService = service.ServiceGroupApp.AnsibleServiceGroup.UserService

var systemUserService = service.ServiceGroupApp.SystemServiceGroup.UserService