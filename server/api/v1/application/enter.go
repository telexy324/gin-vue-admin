package application

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CmdbApi
}

var cmdbService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbService
