package taskPool

import "github.com/flipped-aurora/gin-vue-admin/server/service"

var taskService = service.ServiceGroupApp.TaskServiceGroup
var sshService = service.ServiceGroupApp.SshServiceGroup
var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
