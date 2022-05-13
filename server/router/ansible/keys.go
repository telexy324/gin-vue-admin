package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type KeyRouter struct {
}

func (s *KeyRouter) InitKeyRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	keyRouter := Router.Group("ansible/key").Use(middleware.OperationRecord())
	keyRouterWithoutRecord := Router.Group("ansible/key")
	var authorityKeyApi = v1.ApiGroupApp.AnsibleApiGroup.KeysApi
	{
		keyRouterWithoutRecord.POST("addKey", authorityKeyApi.AddKey)       // 新增菜单
		keyRouter.POST("deleteKey", authorityKeyApi.DeleteKey)              // 删除菜单
		keyRouterWithoutRecord.POST("updateKey", authorityKeyApi.UpdateKey) // 更新菜单
	}
	{
		keyRouterWithoutRecord.POST("getKeyById", authorityKeyApi.GetKeyById)     // 获取菜单树
		keyRouterWithoutRecord.POST("getKeyList", authorityKeyApi.GetKeyList) // 分页获取基础menu列表
	}
	return keyRouter
}