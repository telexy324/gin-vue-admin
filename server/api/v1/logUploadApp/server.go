package logUploadApp

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl/request"
	logUploadRes "github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogUploadApi struct {
}

// @Tags LogUpload
// @Summary 新增服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body logUploadMdl.Server true "主机名, 管理ip"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /logUpload/addServer [post]
func (a *LogUploadApi) AddServer(c *gin.Context) {
	var server logUploadMdl.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server, utils.LogServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.AddServer(server); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags LogUpload
// @Summary 删除服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteServer [post]
func (a *LogUploadApi) DeleteServer(c *gin.Context) {
	var server request.GetById
	if err := c.ShouldBindJSON(&server); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.DeleteServer(server.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags LogUpload
// @Summary 批量删除服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteServerByIds [post]
func (a *LogUploadApi) DeleteServerByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := serverService.DeleteServerByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags LogUpload
// @Summary 更新服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body logUploadMdl.Server true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /logUpload/updateServer [post]
func (a *LogUploadApi) UpdateServer(c *gin.Context) {
	var server logUploadMdl.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.UpdateServer(server); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags LogUpload
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getServerById [post]
func (a *LogUploadApi) GetServerById(c *gin.Context) {
	var idInfo request.GetById
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, server := serverService.GetServerById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(logUploadRes.ServerResponse{
			Server: server,
		}, "获取成功", c)
	}
}

// @Tags LogUpload
// @Summary 分页获取基础server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getServerList [post]
func (a *LogUploadApi) GetServerList(c *gin.Context) {
	var pageInfo request2.ServerSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, serverList, total := serverService.GetServerList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     serverList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
