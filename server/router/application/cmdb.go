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
	var authorityServerApi = v1.ApiGroupApp.ApplicationApiGroup.CmdbServerApi
	{
		cmdbRouterWithoutRecord.POST("addServer", authorityServerApi.AddServer)       // 新增菜单
		cmdbRouter.POST("deleteServer", authorityServerApi.DeleteServer)              // 删除菜单
		cmdbRouterWithoutRecord.POST("updateServer", authorityServerApi.UpdateServer) // 更新菜单
	}
	{
		cmdbRouterWithoutRecord.POST("getServerById", authorityServerApi.GetServerById)     // 获取菜单树
		cmdbRouterWithoutRecord.POST("getServerList", authorityServerApi.GetServerList) // 分页获取基础menu列表
	}
	{
		cmdbRouterWithoutRecord.POST("server/addRelation", authorityServerApi.AddRelation)   // 获取菜单树
		cmdbRouterWithoutRecord.POST("server/relations", authorityServerApi.ServerRelations) // 分页获取基础menu列表
	}
	{
		cmdbRouterWithoutRecord.POST("importExcel", authorityServerApi.ImportExcel)          // 导入Excel
		cmdbRouterWithoutRecord.POST("exportExcel", authorityServerApi.ExportExcel)          // 导出Excel
		cmdbRouterWithoutRecord.GET("downloadTemplate", authorityServerApi.DownloadTemplate) // 下载模板文件
	}

	var authoritySystemApi = v1.ApiGroupApp.ApplicationApiGroup.CmdbSystemApi
	{
		cmdbRouterWithoutRecord.POST("addSystem", authoritySystemApi.AddSystem)       // 新增菜单
		cmdbRouter.POST("deleteSystem", authoritySystemApi.DeleteSystem)              // 删除菜单
		cmdbRouterWithoutRecord.POST("updateSystem", authoritySystemApi.UpdateSystem) // 更新菜单
	}
	{
		cmdbRouterWithoutRecord.POST("getSystemById", authoritySystemApi.GetSystemById)     // 获取菜单树
		cmdbRouterWithoutRecord.POST("getSystemList", authoritySystemApi.GetSystemList) // 分页获取基础menu列表
	}
	{
		cmdbRouterWithoutRecord.POST("system/addRelation", authoritySystemApi.AddRelation)   // 获取菜单树
		cmdbRouterWithoutRecord.POST("system/relations", authoritySystemApi.SystemRelations) // 分页获取基础menu列表
	}
	return cmdbRouter
}
