package logUploadRoute

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type LogUploadRouter struct {
}

func (s *LogUploadRouter) InitLogUploadRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	logUploadRouter := Router.Group("logUpload").Use(middleware.OperationRecord())
	var authorityServerApi = v1.ApiGroupApp.LogUploadApiGroup.LogUploadApi
	{
		logUploadRouter.POST("addServer", authorityServerApi.AddServer)                 // 新增菜单
		logUploadRouter.POST("deleteServer", authorityServerApi.DeleteServer)           // 删除菜单
		logUploadRouter.POST("deleteServerByIds", authorityServerApi.DeleteServerByIds) // 删除菜单
		logUploadRouter.POST("updateServer", authorityServerApi.UpdateServer)           // 更新菜单
	}
	{
		logUploadRouter.POST("getServerById", authorityServerApi.GetServerById) // 获取菜单树
		logUploadRouter.POST("getServerList", authorityServerApi.GetServerList) // 分页获取基础menu列表
	}
	return logUploadRouter
}
