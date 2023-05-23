package logUploadApp

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	LogUploadApi
}

var serverService = service.ServiceGroupApp.LogUploadServiceGroup.ServerService
