package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	userRouter := Router.Group("ansible/user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("ansible/user")
	var authorityUserApi = v1.ApiGroupApp.AnsibleApiGroup.UsersApi
	{
		userRouterWithoutRecord.POST("addUser", authorityUserApi.AddUser)       // 新增菜单
		userRouter.POST("deleteUser", authorityUserApi.DeleteUser)              // 删除菜单
		userRouterWithoutRecord.POST("updateUser", authorityUserApi.UpdateUser) // 更新菜单
	}
	{
		userRouterWithoutRecord.POST("getProjectUsers", authorityUserApi.GetProjectUsers)     // 获取菜单树
	}
	return userRouter
}