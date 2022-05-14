package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InventoryApi struct {
}

//func GetInventoryRefs(w http.ResponseWriter, r *http.Request) {
//	inventory := context.Get(r, "inventory").(db.Inventory)
//	refs, err := helpers.Store(r).GetInventoryRefs(inventory.ProjectID, inventory.ID)
//	if err != nil {
//		helpers.WriteError(w, err)
//		return
//	}
//
//	helpers.WriteJSON(w, http.StatusOK, refs)
//}

// @Tags Inventory
// @Summary 新增Inventory
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Inventory true "Inventory"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/inventory/addInventory [post]
func (a *InventoryApi) AddInventory(c *gin.Context) {
	var inventory ansible.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(inventory, utils.InventoryVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := inventoryService.CreateInventory(inventory); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Inventory
// @Summary 删除Inventory
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "InventoryId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/inventory/deleteInventory [post]
func (a *InventoryApi) DeleteInventory(c *gin.Context) {
	var inventory request.GetByProjectId
	if err := c.ShouldBindJSON(&inventory); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(inventory, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := inventoryService.DeleteInventory(inventory.ProjectId, inventory.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Inventory
// @Summary 更新Inventory
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Inventory true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/inventory/updateInventory [post]
func (a *InventoryApi) UpdateInventory(c *gin.Context) {
	var inventory ansible.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(inventory, utils.InventoryVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := inventoryService.UpdateInventory(inventory); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Inventory
// @Summary 根据id获取Inventory
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "InventoryId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/inventory/getInventoryById [post]
func (a *InventoryApi) GetInventoryById(c *gin.Context) {
	var idInfo request.GetByProjectId
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if inventory, err := inventoryService.GetInventory(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.InventoryResponse{
			Inventory: inventory,
		}, "获取成功", c)
	}
}

// @Tags Inventory
// @Summary 分页获取基础Inventory列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/inventory/getInventoryList [post]
func (a *InventoryApi) GetInventoryList(c *gin.Context) {
	var pageInfo request.GetByProjectId
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, inventories, total := inventoryService.GetInventories(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     inventories,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
