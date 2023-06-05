package taskPool

import "github.com/flipped-aurora/gin-vue-admin/server/service"

var taskService = service.ServiceGroupApp.TaskServiceGroup
var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
var applicationService = service.ServiceGroupApp.ApplicationServiceGroup
