package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	applicationReq "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ApplicationRecordApi struct {
}

// @Tags ApplicationRecord
// @Summary 创建ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationRecord true "创建ApplicationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/createApplicationRecord [post]
func (s *ApplicationRecordApi) CreateApplicationRecord(c *gin.Context) {
	var applicationRecord application.ApplicationRecord
	_ = c.ShouldBindJSON(&applicationRecord)
	if err := applicationRecordService.CreateApplicationRecord(applicationRecord); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags ApplicationRecord
// @Summary 删除ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationRecord true "ApplicationRecord模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteApplicationRecord [delete]
func (s *ApplicationRecordApi) DeleteApplicationRecord(c *gin.Context) {
	var applicationRecord application.ApplicationRecord
	_ = c.ShouldBindJSON(&applicationRecord)
	if err := applicationRecordService.DeleteApplicationRecord(applicationRecord); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags ApplicationRecord
// @Summary 批量删除ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ApplicationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /cmdb/deleteApplicationRecordByIds [delete]
func (s *ApplicationRecordApi) DeleteApplicationRecordByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := applicationRecordService.DeleteApplicationRecordByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags ApplicationRecord
// @Summary 用id查询ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query application.ApplicationRecord true "Id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /cmdb/findApplicationRecord [get]
func (s *ApplicationRecordApi) FindApplicationRecord(c *gin.Context) {
	var applicationRecord application.ApplicationRecord
	_ = c.ShouldBindQuery(&applicationRecord)
	if err := utils.Verify(applicationRecord, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, reapplicationRecordRecord := applicationRecordService.GetApplicationRecord(applicationRecord.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithDetailed(gin.H{"reapplicationRecordRecord": reapplicationRecordRecord}, "查询成功", c)
	}
}

// @Tags ApplicationRecord
// @Summary 分页获取ApplicationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.ApplicationRecordSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getApplicationRecordList [get]
func (s *ApplicationRecordApi) GetApplicationRecordList(c *gin.Context) {
	var pageInfo applicationReq.ApplicationRecordSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := applicationRecordService.GetApplicationRecordInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags ApplicationRecord
// @Summary 导出Excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body request.IdsReq true "导出Excel文件信息"
// @Success 200
// @Router /cmdb/exportApplicationRecord [post]
func (s *ApplicationRecordApi) ExportApplicationRecord(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	//filePath := global.GVA_CONFIG.Excel.Dir + excelInfo.FileName
	buf, err := applicationRecordService.ParseInfoList2Excel(IDS.Ids)
	if err != nil {
		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
		response.FailWithMessage("转换Excel失败", c)
		return
	}
	fileName := "logRecord" + strconv.Itoa(int(time.Now().UnixNano()))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("success", "true")

	if _, err = c.Writer.Write(buf.Bytes()); err != nil {
		global.GVA_LOG.Error("下载文件失败!", zap.Any("err", err))
		response.FailWithMessage("download file failed", c)
	}
}
