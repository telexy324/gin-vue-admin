package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/application"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/logUploadApp"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/taskApp"
)

type ApiGroup struct {
	ExampleApiGroup     example.ApiGroup
	SystemApiGroup      system.ApiGroup
	AutoCodeApiGroup    autocode.ApiGroup
	ApplicationApiGroup application.ApiGroup
	SshApiGroup         ssh.ApiGroup
	TaskApiGroup        taskApp.ApiGroup
	TaskSocketApiGroup  socket.ApiGroup
	LogUploadApiGroup   logUploadApp.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
