package ansible

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"

	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/context"
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

// GetInventory returns an inventory from the database
func GetInventory(w http.ResponseWriter, r *http.Request) {
	if inventory := context.Get(r, "inventory"); inventory != nil {
		helpers.WriteJSON(w, http.StatusOK, inventory.(db.Inventory))
		return
	}

	project := context.Get(r, "project").(db.Project)

	inventories, err := helpers.Store(r).GetInventories(project.ID, helpers.QueryParams(r.URL))

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, inventories)
}

// AddInventory creates an inventory in the database
func AddInventory(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)

	var inventory db.Inventory

	if !helpers.Bind(w, r, &inventory) {
		return
	}

	if inventory.ProjectID != project.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	switch inventory.Type {
	case db.InventoryStatic, db.InventoryFile:
		break
	default:
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Not supported inventory type",
		})
		return
	}

	newInventory, err := helpers.Store(r).CreateInventory(inventory)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	objType := db.EventInventory
	desc := "Inventory " + inventory.Name + " created"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &project.ID,
		ObjectType:  &objType,
		ObjectID:    &newInventory.ID,
		Description: &desc,
	})

	if err != nil {
		// Write error to log but return ok to user, because inventory created
		log.Error(err)
	}

	helpers.WriteJSON(w, http.StatusCreated, newInventory)
}

// IsValidInventoryPath tests a path to ensure it is below the cwd
func IsValidInventoryPath(path string) bool {

	currentPath, err := os.Getwd()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(currentPath, absPath)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(relPath, "..")
}

// UpdateInventory writes updated values to an existing inventory item in the database
func UpdateInventory(w http.ResponseWriter, r *http.Request) {
	oldInventory := context.Get(r, "inventory").(db.Inventory)

	var inventory db.Inventory

	if !helpers.Bind(w, r, &inventory) {
		return
	}

	if inventory.ID != oldInventory.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Inventory ID in body and URL must be the same",
		})
		return
	}

	if inventory.ProjectID != oldInventory.ProjectID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	switch inventory.Type {
	case db.InventoryStatic:
		break
	case db.InventoryFile:
		if !IsValidInventoryPath(inventory.Inventory) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := helpers.Store(r).UpdateInventory(inventory)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveInventory deletes an inventory from the database
func RemoveInventory(w http.ResponseWriter, r *http.Request) {
	inventory := context.Get(r, "inventory").(db.Inventory)
	var err error

	err = helpers.Store(r).DeleteInventory(inventory.ProjectID, inventory.ID)
	if err == db.ErrInvalidOperation {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Inventory is in use by one or more templates",
			"inUse": true,
		})
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	desc := "Inventory " + inventory.Name + " deleted"

	user := context.Get(r, "user").(*db.User)

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &inventory.ProjectID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Inventory
// @Summary 新增Inventory
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Inventory true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/inventory/addInventory [post]
func (a *EnvironmentApi) AddInventory(c *gin.Context) {
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
// @Param data body request.GetById true "InventoryId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/inventory/deleteInventory [post]
func (a *EnvironmentApi) DeleteInventory(c *gin.Context) {
	var environment request2.GetByProjectId
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := inventoryService.DeleteInventory(environment.ProjectId, environment.ID); err != nil {
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
func (a *EnvironmentApi) UpdateInventory(c *gin.Context) {
	var environment ansible.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.EnvironmentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := environmentService.UpdateEnvironment(environment); err != nil {
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
// @Param data body request2.GetByProjectId true "Inventoryid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/inventory/getInventoryById [post]
func (a *EnvironmentApi) GetInventoryById(c *gin.Context) {
	var idInfo request2.GetByProjectId
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
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/inventory/getInventoryList[post]
func (a *EnvironmentApi) GetInventoryList(c *gin.Context) {
	var pageInfo request2.GetByProjectId
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
