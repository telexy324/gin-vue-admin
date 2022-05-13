package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type EnvironmentRouter struct {
}

func (s *EnvironmentRouter) InitEnvironmentRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	environmentRouter := Router.Group("ansible/environment").Use(middleware.OperationRecord())
	environmentRouterWithoutRecord := Router.Group("ansible/environment")
	var authorityEnvironmentApi = v1.ApiGroupApp.AnsibleApiGroup.EnvironmentApi
	{
		environmentRouterWithoutRecord.POST("addEnvironment", authorityEnvironmentApi.AddEnvironment)       // 新增菜单
		environmentRouter.POST("deleteEnvironment", authorityEnvironmentApi.DeleteEnvironment)              // 删除菜单
		environmentRouterWithoutRecord.POST("updateEnvironment", authorityEnvironmentApi.UpdateEnvironment) // 更新菜单
	}
	{
		environmentRouterWithoutRecord.POST("getEnvironmentById", authorityEnvironmentApi.GetEnvironmentById)     // 获取菜单树
		environmentRouterWithoutRecord.POST("getEnvironmentList", authorityEnvironmentApi.GetEnvironmentList) // 分页获取基础menu列表
	}
	return environmentRouter
}
