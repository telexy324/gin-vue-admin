package ssh

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SshRouter struct {
}

func (s *SshRouter) InitSshRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	sshRouter := Router.Group("ssh").Use(middleware.OperationRecord())
	sshRouterWithoutRecord := Router.Group("ssh")
	var authoritySshApi = v1.ApiGroupApp.SshApiGroup.SshApi
	{
		sshRouterWithoutRecord.GET("addAdmin", authoritySshApi.ShellWeb)       // 新增管理员
	}
	return sshRouter
}
