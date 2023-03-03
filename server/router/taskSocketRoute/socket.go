package taskSocketRoute

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TaskSocketRouter struct {
}

func (s *TaskSocketRouter) InitTaskSocketRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	taskSocketRouter := Router.Group("task").Use(middleware.OperationRecord())
	taskSocketRouterWithoutRecord := Router.Group("task")
	var authorityTaskSocketApi = v1.ApiGroupApp.TaskSocketApiGroup.SocketApi
	{
		taskSocketRouterWithoutRecord.GET("ws", authorityTaskSocketApi.Handler) // 新增管理员
	}
	return taskSocketRouter
}
