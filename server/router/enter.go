package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/application"
	"github.com/flipped-aurora/gin-vue-admin/server/router/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/taskRoute"
)

type RouterGroup struct {
	System      system.RouterGroup
	Example     example.RouterGroup
	Autocode    autocode.RouterGroup
	Application application.RouterGroup
	Ssh         ssh.RouterGroup
	Task        taskRoute.TaskRouter
}

var RouterGroupApp = new(RouterGroup)
