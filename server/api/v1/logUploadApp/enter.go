package logUploadApp

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	LogUploadServerApi
	LogUploadSecretApi
}

var serverService = service.ServiceGroupApp.LogUploadServiceGroup.ServerService
var secretService = service.ServiceGroupApp.LogUploadServiceGroup.SecretService
