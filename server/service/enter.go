package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/service/application"
	"github.com/flipped-aurora/gin-vue-admin/server/service/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	ExampleServiceGroup  example.ServiceGroup
	SystemServiceGroup   system.ServiceGroup
	AutoCodeServiceGroup autocode.ServiceGroup
	ApplicationServiceGroup application.ServiceGroup
	AnsibleServiceGroup  ansible.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
