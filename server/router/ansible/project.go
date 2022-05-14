package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProjectRouter struct {
}

func (s *ProjectRouter) InitProjectRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	projectRouter := Router.Group("ansible/project").Use(middleware.OperationRecord())
	projectRouterWithoutRecord := Router.Group("ansible/project")
	var authorityProjectApi = v1.ApiGroupApp.AnsibleApiGroup.ProjectApi
	{
		projectRouterWithoutRecord.POST("addProject", authorityProjectApi.AddProject)       // 新增菜单
		projectRouter.POST("deleteProject", authorityProjectApi.DeleteProject)              // 删除菜单
		projectRouterWithoutRecord.POST("updateProject", authorityProjectApi.UpdateProject) // 更新菜单
	}
	{
		projectRouterWithoutRecord.POST("getProjectById", authorityProjectApi.GetProjectById) // 获取菜单树
		projectRouterWithoutRecord.POST("getProjectList", authorityProjectApi.GetProjectList) // 分页获取基础menu列表
	}
	{
		projectRouterWithoutRecord.POST("mustBeAdmin", authorityProjectApi.MustBeAdmin) // 获取菜单树
	}
	return projectRouter
}
