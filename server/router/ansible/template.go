package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TemplateRouter struct {
}

func (s *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	templateRouter := Router.Group("ansible/template").Use(middleware.OperationRecord())
	templateRouterWithoutRecord := Router.Group("ansible/template")
	var authorityTemplateApi = v1.ApiGroupApp.AnsibleApiGroup.TemplatesApi
	{
		templateRouterWithoutRecord.POST("addTemplate", authorityTemplateApi.AddTemplate)       // 新增菜单
		templateRouter.POST("deleteTemplate", authorityTemplateApi.DeleteTemplate)              // 删除菜单
		templateRouterWithoutRecord.POST("updateTemplate", authorityTemplateApi.UpdateTemplate) // 更新菜单
	}
	{
		templateRouterWithoutRecord.POST("getTemplateById", authorityTemplateApi.GetTemplateById) // 获取菜单树
		templateRouterWithoutRecord.POST("getTemplateList", authorityTemplateApi.GetTemplateList) // 分页获取基础menu列表
	}
	return templateRouter
}
