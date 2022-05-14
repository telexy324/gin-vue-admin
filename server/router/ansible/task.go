package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
}

func (s *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	taskRouter := Router.Group("ansible/task").Use(middleware.OperationRecord())
	taskRouterWithoutRecord := Router.Group("ansible/task")
	var authorityTaskApi = v1.ApiGroupApp.AnsibleApiGroup.TasksApi
	{
		taskRouterWithoutRecord.POST("addTask", authorityTaskApi.AddTask)       // 新增菜单
		taskRouter.POST("deleteTask", authorityTaskApi.DeleteTask)              // 删除菜单
		taskRouterWithoutRecord.POST("updateTask", authorityTaskApi.UpdateTask) // 更新菜单
	}
	{
		taskRouterWithoutRecord.POST("getTaskById", authorityTaskApi.GetTaskById) // 获取菜单树
		taskRouterWithoutRecord.POST("getTaskList", authorityTaskApi.GetTaskList) // 分页获取基础menu列表
	}
	{
		taskRouterWithoutRecord.POST("getTaskOutputs", authorityTaskApi.GetTaskOutputs) // 获取菜单树
		taskRouterWithoutRecord.POST("stopTask", authorityTaskApi.StopTask)             // 分页获取基础menu列表
	}
	return taskRouter
}
