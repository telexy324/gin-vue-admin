package application

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CmdbRouter struct {
}

func (s *CmdbRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	cmdbRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	cmdbRouterWithoutRecord := Router.Group("cmdb")
	var authorityMenuApi = v1.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	{
		cmdbRouter.POST("addAsset", authorityMenuApi.AddBaseMenu)           // 新增菜单
		cmdbRouter.POST("deleteAsset", authorityMenuApi.DeleteBaseMenu)     // 删除菜单
		cmdbRouter.POST("updateAsset", authorityMenuApi.UpdateBaseMenu)     // 更新菜单
	}
	{
		cmdbRouterWithoutRecord.POST("getAsset", authorityMenuApi.GetMenu)                   // 获取菜单树
		cmdbRouterWithoutRecord.POST("getAssetList", authorityMenuApi.GetMenuList)           // 分页获取基础menu列表
	}
	return cmdbRouter
}
