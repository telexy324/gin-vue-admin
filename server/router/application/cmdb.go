package application

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CmdbRouter struct {
}

func (s *CmdbRouter) InitCmdbRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	cmdbRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	cmdbRouterWithoutRecord := Router.Group("cmdb")
	var authorityServerApi = v1.ApiGroupApp.ApplicationApiGroup.CmdbApi
	{
		cmdbRouterWithoutRecord.POST("addServer", authorityServerApi.AddServer)       // 新增菜单
		cmdbRouter.POST("deleteServer", authorityServerApi.DeleteServer)              // 删除菜单
		cmdbRouterWithoutRecord.POST("updateServer", authorityServerApi.UpdateServer) // 更新菜单
	}
	{
		cmdbRouterWithoutRecord.POST("getServer", authorityServerApi.GetServerById)     // 获取菜单树
		cmdbRouterWithoutRecord.POST("getServerList", authorityServerApi.GetServerList) // 分页获取基础menu列表
	}
	{
		cmdbRouterWithoutRecord.POST("system/addRelation", authorityServerApi.AddRelation)   // 获取菜单树
		cmdbRouterWithoutRecord.POST("system/relations", authorityServerApi.SystemRelations) // 分页获取基础menu列表
	}
	return cmdbRouter
}
