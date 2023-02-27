package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/application"
	"github.com/flipped-aurora/gin-vue-admin/server/service/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/scheduleSvr"
	"github.com/flipped-aurora/gin-vue-admin/server/service/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/taskRunnerSvr"
	"github.com/flipped-aurora/gin-vue-admin/server/service/taskSvr"
)

type ServiceGroup struct {
	ExampleServiceGroup     example.ServiceGroup
	SystemServiceGroup      system.ServiceGroup
	AutoCodeServiceGroup    autocode.ServiceGroup
	ApplicationServiceGroup application.ServiceGroup
	SshServiceGroup         ssh.ServiceGroup
	TaskServiceGroup        taskSvr.ServiceGroup
	TaskRunnerServiceGroup  taskRunnerSvr.ServiceGroup
	ScheduleServiceGroup    scheduleSvr.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
