package ssh

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	SshApi
}

var sshService = service.ServiceGroupApp.SshServiceGroup.SshService
