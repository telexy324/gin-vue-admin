package jumpServerRoute

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/jumpServer"
	"github.com/gin-gonic/gin"
)

type JumpServerRouter struct {
}

func (s *JumpServerRouter) InitJumpServerRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	jumpServerRouter := Router.Group("jumpServer").Use(middleware.OperationRecord())
	var jumpServerApi = v1.ApiGroupApp.JumpServerApiGroup.JumpServerApi
	{
		jumpServerRouter.POST("getToken", jumpServerApi.GetToken) // 获得token // 获取服务器
	}
	return jumpServerRouter
}

func (s *JumpServerRouter) InitJumpServerPubRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	jumpServerRouter := Router.Group("jumpServer")
	{
		jumpServerRouter.GET("ws", jumpServer.HandleWebSocket)
	}
	return jumpServerRouter
}
