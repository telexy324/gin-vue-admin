package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ScheduleRouter struct {
}

func (s *ScheduleRouter) InitScheduleRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	scheduleRouter := Router.Group("ansible/schedule").Use(middleware.OperationRecord())
	scheduleRouterWithoutRecord := Router.Group("ansible/schedule")
	var authorityScheduleApi = v1.ApiGroupApp.AnsibleApiGroup.SchedulesApi
	{
		scheduleRouterWithoutRecord.POST("addSchedule", authorityScheduleApi.AddSchedule)       // 新增菜单
		scheduleRouter.POST("deleteSchedule", authorityScheduleApi.DeleteSchedule)              // 删除菜单
		scheduleRouterWithoutRecord.POST("updateSchedule", authorityScheduleApi.UpdateSchedule) // 更新菜单
	}
	{
		scheduleRouterWithoutRecord.POST("getScheduleById", authorityScheduleApi.GetScheduleById)     // 获取菜单树
		scheduleRouterWithoutRecord.POST("getTemplateScheduleList", authorityScheduleApi.GetTemplateScheduleList) // 分页获取基础menu列表
	}
	return scheduleRouter
}