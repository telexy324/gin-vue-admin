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

type LogUploadSecretApi struct {
}

// @Tags LogUpload
// @Summary 新增密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body logUploadMdl.Secret true "主机名, 管理ip"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /logUpload/addSecret [post]
func (a *LogUploadSecretApi) AddSecret(c *gin.Context) {
	var secret logUploadMdl.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(secret, utils.LogSecretVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := secretService.AddSecret(secret); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags LogUpload
// @Summary 删除密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "密钥id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteSecret [post]
func (a *LogUploadSecretApi) DeleteSecret(c *gin.Context) {
	var secret request.GetById
	if err := c.ShouldBindJSON(&secret); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(secret, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := secretService.DeleteSecret(secret.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags LogUpload
// @Summary 批量删除密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteSecretByIds [post]
func (a *LogUploadSecretApi) DeleteSecretByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := secretService.DeleteSecretByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags LogUpload
// @Summary 更新密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body logUploadMdl.Secret true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /logUpload/updateSecret [post]
func (a *LogUploadSecretApi) UpdateSecret(c *gin.Context) {
	var secret logUploadMdl.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(secret, utils.LogSecretVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := secretService.UpdateSecret(secret); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags LogUpload
// @Summary 根据id获取密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "密钥id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getSecretById [post]
func (a *LogUploadSecretApi) GetSecretById(c *gin.Context) {
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
	if err, secret := secretService.GetSecretById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(logUploadRes.SecretResponse{
			Secret: secret,
		}, "获取成功", c)
	}
}

// @Tags LogUpload
// @Summary 分页获取基础secret列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getSecretList [post]
func (a *LogUploadSecretApi) GetSecretList(c *gin.Context) {
	var pageInfo request2.SecretSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, secretList, total := secretService.GetSecretList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     secretList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

