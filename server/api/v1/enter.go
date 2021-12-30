package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/application"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

type ApiGroup struct {
	ExampleApiGroup     example.ApiGroup
	SystemApiGroup      system.ApiGroup
	AutoCodeApiGroup    autocode.ApiGroup
	ApplicationApiGroup application.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
