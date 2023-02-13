package application

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SshRouter struct {
}

func (s *SshRouter) InitSshRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	cmdbRouter := Router.Group("ssh").Use(middleware.OperationRecord())
	cmdbRouterWithoutRecord := Router.Group("ssh")
	var authorityStaffApi = v1.ApiGroupApp.ApplicationApiGroup.StaffApi
	{
		cmdbRouterWithoutRecord.POST("addAdmin", authorityStaffApi.AddAdmin)       // 新增管理员
		cmdbRouter.POST("deleteAdmin", authorityStaffApi.DeleteAdmin)              // 删除管理员
		cmdbRouterWithoutRecord.POST("updateAdmin", authorityStaffApi.UpdateAdmin) // 更新管理员
	}
	{
		cmdbRouterWithoutRecord.POST("getAdminById", authorityStaffApi.GetAdminById)     // 获取管理员
		cmdbRouterWithoutRecord.POST("getAdminList", authorityStaffApi.GetAdminList) // 分页获取管理员列表
		cmdbRouterWithoutRecord.POST("getDepartmentAll", authorityStaffApi.GetDepartmentAll) // 获取所有部门
	}
	return cmdbRouter
}
