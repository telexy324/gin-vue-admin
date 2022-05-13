package ansible

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InventoryRouter struct {
}

func (s *InventoryRouter) InitInventoryRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	inventoryRouter := Router.Group("ansible/inventory").Use(middleware.OperationRecord())
	inventoryRouterWithoutRecord := Router.Group("ansible/inventory")
	var authorityInventoryApi = v1.ApiGroupApp.AnsibleApiGroup.InventoryApi
	{
		inventoryRouterWithoutRecord.POST("addInventory", authorityInventoryApi.AddInventory)       // 新增菜单
		inventoryRouter.POST("deleteInventory", authorityInventoryApi.DeleteInventory)              // 删除菜单
		inventoryRouterWithoutRecord.POST("updateInventory", authorityInventoryApi.UpdateInventory) // 更新菜单
	}
	{
		inventoryRouterWithoutRecord.POST("getInventoryById", authorityInventoryApi.GetInventoryById)     // 获取菜单树
		inventoryRouterWithoutRecord.POST("getInventoryList", authorityInventoryApi.GetInventoryList) // 分页获取基础menu列表
	}
	return inventoryRouter
}