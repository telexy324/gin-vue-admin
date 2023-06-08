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
	var authorityTaskApi = v1.ApiGroupApp.TaskApiGroup.TaskApi
	{
		taskRouter.POST("addTask", authorityTaskApi.AddTask)       // 新增Task
		taskRouter.POST("deleteTask", authorityTaskApi.DeleteTask) // 删除Task
	}
	{
		taskRouter.POST("getTaskById", authorityTaskApi.GetTaskById) // 获取Task
		taskRouter.POST("getTaskList", authorityTaskApi.GetTaskList) // 分页获取Task
	}
	{
		taskRouter.POST("getTaskOutputs", authorityTaskApi.GetTaskOutputs)             // 获取Task输出
		taskRouter.POST("stopTask", authorityTaskApi.StopTask)                         // 停止Task
		taskRouter.POST("getTaskDashboardInfo", authorityTaskApi.GetTaskDashboardInfo) // 获取Task面板信息
	}

	var authorityTemplateApi = v1.ApiGroupApp.TaskApiGroup.TemplateApi
	{
		taskRouter.POST("template/addTemplate", authorityTemplateApi.AddTemplate)                 // 新增管理员
		taskRouter.POST("template/deleteTemplate", authorityTemplateApi.DeleteTemplate)           // 删除管理员
		taskRouter.POST("template/deleteTemplateByIds", authorityTemplateApi.DeleteTemplateByIds) // 删除管理员
		taskRouter.POST("template/updateTemplate", authorityTemplateApi.UpdateTemplate)           // 更新管理员
	}
	{
		taskRouter.POST("template/getTemplateById", authorityTemplateApi.GetTemplateById) // 获取管理员
		taskRouter.POST("template/getTemplateList", authorityTemplateApi.GetTemplateList) // 分页获取管理员列表
		taskRouter.POST("template/checkScript", authorityTemplateApi.CheckScript)         // 分页获取管理员列表
		taskRouter.POST("template/uploadScript", authorityTemplateApi.UploadScript)       // 分页获取管理员列表
		taskRouter.POST("template/downloadScript", authorityTemplateApi.DownloadScript)   // 分页获取管理员列表
		taskRouter.POST("template/getFileList", authorityTemplateApi.GetFileList)         // 获取管理员
		taskRouter.GET("template/downloadFile", authorityTemplateApi.DownloadFile)        // 获取管理员
		taskRouter.POST("template/uploadLogServer", authorityTemplateApi.UploadLogServer)
		taskRouter.POST("template/deployServer", authorityTemplateApi.DeployServer)
	}
	{
		taskRouter.POST("template/addSet", authorityTemplateApi.AddSet)                 // 新增管理员
		taskRouter.POST("template/deleteSet", authorityTemplateApi.DeleteSet)           // 删除管理员
		taskRouter.POST("template/deleteSetByIds", authorityTemplateApi.DeleteSetByIds) // 删除管理员
		taskRouter.POST("template/updateSet", authorityTemplateApi.UpdateSet)           // 更新管理员
		taskRouter.POST("template/getSetById", authorityTemplateApi.GetSetById)         // 获取管理员
		taskRouter.POST("template/getSetList", authorityTemplateApi.GetSetList)         // 分页获取管理员列表
		taskRouter.POST("template/addSetTask", authorityTemplateApi.AddSetTask)         // 分页获取管理员列表
		taskRouter.POST("template/processSetTask", authorityTemplateApi.ProcessSetTask) // 分页获取管理员列表
		taskRouter.POST("template/getSetTaskById", authorityTemplateApi.GetSetTaskById) // 获取管理员
		taskRouter.POST("template/getSetTaskList", authorityTemplateApi.GetSetTaskList) // 获取管理员
	}

	var authorityScheduleApi = v1.ApiGroupApp.TaskApiGroup.ScheduleApi
	{
		taskRouter.POST("schedule/addSchedule", authorityScheduleApi.AddSchedule)                 // 新增管理员
		taskRouter.POST("schedule/deleteSchedule", authorityScheduleApi.DeleteSchedule)           // 删除管理员
		taskRouter.POST("schedule/deleteScheduleByIds", authorityScheduleApi.DeleteScheduleByIds) // 删除管理员
		taskRouter.POST("schedule/updateSchedule", authorityScheduleApi.UpdateSchedule)           // 更新管理员
	}
	{
		taskRouter.POST("schedule/getScheduleById", authorityScheduleApi.GetScheduleById)                       // 获取管理员
		taskRouter.POST("schedule/getTemplateScheduleList", authorityScheduleApi.GetTemplateScheduleList)       // 分页获取管理员列表
		taskRouter.POST("schedule/validateScheduleCronFormat", authorityScheduleApi.ValidateScheduleCronFormat) // 获取所有部门
	}
	return taskRouter
}
