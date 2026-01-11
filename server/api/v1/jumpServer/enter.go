package jumpServer

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	JumpServerApi
}

var jumpServerService = service.ServiceGroupApp.JumpServerServiceGroup.JumpServerService
