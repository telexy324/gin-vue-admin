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
	var authorityServerApi = v1.ApiGroupApp.LogUploadApiGroup.LogUploadServerApi
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
	var authoritySecretApi = v1.ApiGroupApp.LogUploadApiGroup.LogUploadSecretApi
	{
		logUploadRouter.POST("addSecret", authoritySecretApi.AddSecret)                 // 新增菜单
		logUploadRouter.POST("deleteSecret", authoritySecretApi.DeleteSecret)           // 删除菜单
		logUploadRouter.POST("deleteSecretByIds", authoritySecretApi.DeleteSecretByIds) // 删除菜单
		logUploadRouter.POST("updateSecret", authoritySecretApi.UpdateSecret)           // 更新菜单
	}
	{
		logUploadRouter.POST("getSecretById", authoritySecretApi.GetSecretById) // 获取菜单树
		logUploadRouter.POST("getSecretList", authoritySecretApi.GetSecretList) // 分页获取基础menu列表
	}
	return logUploadRouter
}
