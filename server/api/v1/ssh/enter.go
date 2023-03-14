package ssh

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	SshApi
}

var cmdbServerService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbServerService
