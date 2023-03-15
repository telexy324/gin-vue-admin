package taskRoute

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
}

func (s *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	taskRouter := Router.Group("task").Use(middleware.OperationRecord())
	// taskRouterWithoutRecord := Router.Group("task")
	var authorityServerApi = v1.ApiGroupApp.TaskApiGroup.TaskApi
	{
		taskRouter.POST("addTask", authorityServerApi.AddTask)       // 新增Task
		taskRouter.POST("deleteTask", authorityServerApi.DeleteTask) // 删除Task
	}
	{
		taskRouter.POST("getTaskById", authorityServerApi.GetTaskById) // 获取Task
		taskRouter.POST("getTaskList", authorityServerApi.GetTaskList) // 分页获取Task
	}
	{
		taskRouter.POST("getTaskOutputs", authorityServerApi.GetTaskOutputs) // 停止Task
		taskRouter.POST("stopTask", authorityServerApi.StopTask)             // 停止Task
	}

	var authorityTemplateApi = v1.ApiGroupApp.TaskApiGroup.TemplateApi
	{
		taskRouter.POST("template/addTemplate", authorityTemplateApi.AddTemplate)       // 新增管理员
		taskRouter.POST("template/deleteTemplate", authorityTemplateApi.DeleteTemplate) // 删除管理员
		taskRouter.POST("template/updateTemplate", authorityTemplateApi.UpdateTemplate) // 更新管理员
	}
	{
		taskRouter.POST("template/getTemplateById", authorityTemplateApi.GetTemplateById) // 获取管理员
		taskRouter.POST("template/getTemplateList", authorityTemplateApi.GetTemplateList) // 分页获取管理员列表
		taskRouter.POST("template/checkScript", authorityTemplateApi.CheckScript)         // 分页获取管理员列表
		taskRouter.POST("template/uploadScript", authorityTemplateApi.UploadScript)       // 分页获取管理员列表
		taskRouter.POST("template/downloadScript", authorityTemplateApi.DownloadScript)   // 分页获取管理员列表
	}

	var authorityScheduleApi = v1.ApiGroupApp.TaskApiGroup.ScheduleApi
	{
		taskRouter.POST("schedule/addSchedule", authorityScheduleApi.AddSchedule)       // 新增管理员
		taskRouter.POST("schedule/deleteSchedule", authorityScheduleApi.DeleteSchedule) // 删除管理员
		taskRouter.POST("schedule/updateSchedule", authorityScheduleApi.UpdateSchedule) // 更新管理员
	}
	{
		taskRouter.POST("schedule/getScheduleById", authorityScheduleApi.GetScheduleById)                       // 获取管理员
		taskRouter.POST("schedule/getTemplateScheduleList", authorityScheduleApi.GetTemplateScheduleList)       // 分页获取管理员列表
		taskRouter.POST("schedule/validateScheduleCronFormat", authorityScheduleApi.ValidateScheduleCronFormat) // 获取所有部门
	}
	return taskRouter
}
