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
	// cmdbRouterWithoutRecord := Router.Group("cmdb")
	var authorityServerApi = v1.ApiGroupApp.ApplicationApiGroup.CmdbServerApi
	{
		cmdbRouter.POST("addServer", authorityServerApi.AddServer)                 // 新增菜单
		cmdbRouter.POST("deleteServer", authorityServerApi.DeleteServer)           // 删除菜单
		cmdbRouter.POST("deleteServerByIds", authorityServerApi.DeleteServerByIds) // 删除菜单
		cmdbRouter.POST("updateServer", authorityServerApi.UpdateServer)           // 更新菜单
	}
	{
		cmdbRouter.POST("getServerById", authorityServerApi.GetServerById)           // 获取菜单树
		cmdbRouter.POST("getServerList", authorityServerApi.GetServerList)           // 分页获取基础menu列表
		cmdbRouter.POST("getSystemServers", authorityServerApi.GetSystemServers)     // 分页获取基础menu列表
		cmdbRouter.GET("getAllServerIds", authorityServerApi.GetAllServerIds)        // 分页获取基础menu列表
		cmdbRouter.POST("getSystemServerIds", authorityServerApi.GetSystemServerIds) // 分页获取基础menu列表
	}
	{
		cmdbRouter.POST("server/addRelation", authorityServerApi.AddRelation)   // 获取菜单树
		cmdbRouter.POST("server/relations", authorityServerApi.ServerRelations) // 分页获取基础menu列表
	}
	{
		cmdbRouter.POST("importExcel", authorityServerApi.ImportExcel)          // 导入Excel
		cmdbRouter.POST("exportExcel", authorityServerApi.ExportExcel)          // 导出Excel
		cmdbRouter.GET("downloadTemplate", authorityServerApi.DownloadTemplate) // 下载模板文件
	}
	{
		cmdbRouter.POST("uploadFile", authorityServerApi.UploadFile) // 获取应用
	}
	{
		cmdbRouter.POST("addApp", authorityServerApi.AddApp)       // 新增应用
		cmdbRouter.POST("deleteApp", authorityServerApi.DeleteApp) // 删除应用
		cmdbRouter.POST("updateApp", authorityServerApi.UpdateApp) // 更新应用
	}
	{
		cmdbRouter.POST("getAppById", authorityServerApi.GetAppById) // 获取应用
		cmdbRouter.POST("getAppList", authorityServerApi.GetAppList) // 分页获取基础应用列表
	}

	var authoritySystemApi = v1.ApiGroupApp.ApplicationApiGroup.CmdbSystemApi
	{
		cmdbRouter.POST("addSystem", authoritySystemApi.AddSystem)                 // 新增菜单
		cmdbRouter.POST("deleteSystem", authoritySystemApi.DeleteSystem)           // 删除菜单
		cmdbRouter.POST("deleteSystemByIds", authoritySystemApi.DeleteSystemByIds) // 删除菜单
		cmdbRouter.POST("updateSystem", authoritySystemApi.UpdateSystem)           // 更新菜单
	}
	{
		cmdbRouter.POST("getSystemById", authoritySystemApi.GetSystemById)     // 获取菜单树
		cmdbRouter.POST("getSystemList", authoritySystemApi.GetSystemList)     // 分页获取基础menu列表
		cmdbRouter.POST("getAdminSystems", authoritySystemApi.GetAdminSystems) // 分页获取基础menu列表
	}
	{
		cmdbRouter.POST("system/addRelation", authoritySystemApi.AddRelation)   // 获取菜单树
		cmdbRouter.POST("system/relations", authoritySystemApi.SystemRelations) // 分页获取基础menu列表
	}
	{
		cmdbRouter.POST("system/addEditRelation", authoritySystemApi.AddEditRelation)              // 获取菜单树
		cmdbRouter.POST("system/deleteEditRelation", authoritySystemApi.DeleteEditRelation)        // 分页获取基础menu列表
		cmdbRouter.POST("system/updateEditRelation", authoritySystemApi.UpdateEditRelation)        // 获取菜单树
		cmdbRouter.POST("system/getSystemEditRelation", authoritySystemApi.GetSystemEditRelations) // 分页获取基础menu列表
	}

	var authorityStaffApi = v1.ApiGroupApp.ApplicationApiGroup.StaffApi
	{
		cmdbRouter.POST("addAdmin", authorityStaffApi.AddAdmin)       // 新增管理员
		cmdbRouter.POST("deleteAdmin", authorityStaffApi.DeleteAdmin) // 删除管理员
		cmdbRouter.POST("updateAdmin", authorityStaffApi.UpdateAdmin) // 更新管理员
	}
	{
		cmdbRouter.POST("getAdminById", authorityStaffApi.GetAdminById)         // 获取管理员
		cmdbRouter.POST("getAdminList", authorityStaffApi.GetAdminList)         // 分页获取管理员列表
		cmdbRouter.POST("getDepartmentAll", authorityStaffApi.GetDepartmentAll) // 获取所有部门
	}

	var authorityApplicationRecordApi = v1.ApiGroupApp.ApplicationApiGroup.ApplicationRecordApi
	{
		cmdbRouter.POST("createApplicationRecord", authorityApplicationRecordApi.CreateApplicationRecord)           // 新增操作记录
		cmdbRouter.POST("deleteApplicationRecord", authorityApplicationRecordApi.DeleteApplicationRecord)           // 删除操作记录
		cmdbRouter.POST("deleteApplicationRecordByIds", authorityApplicationRecordApi.DeleteApplicationRecordByIds) // 删除操作记录
	}
	{
		cmdbRouter.GET("findApplicationRecord", authorityApplicationRecordApi.FindApplicationRecord)       // 获取操作记录树
		cmdbRouter.GET("getApplicationRecordList", authorityApplicationRecordApi.GetApplicationRecordList) // 分页获取操作记录
		cmdbRouter.POST("exportApplicationRecord", authorityApplicationRecordApi.ExportApplicationRecord)  // 获取管理员
	}
	return cmdbRouter
}
