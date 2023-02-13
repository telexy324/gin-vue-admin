package application

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	StaffApi
}

var staffService = service.ServiceGroupApp.SshServiceGroup.SshService
