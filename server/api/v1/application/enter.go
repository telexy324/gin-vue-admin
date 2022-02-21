package application

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CmdbApi
	StaffApi
}

var cmdbService = service.ServiceGroupApp.ApplicationServiceGroup.CmdbService
var staffService = service.ServiceGroupApp.ApplicationServiceGroup.StaffService
