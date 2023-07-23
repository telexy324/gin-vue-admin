package taskApp

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
	scheduleReq "github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl/request"
	scheduleRes "github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl/response"
	schedules "github.com/flipped-aurora/gin-vue-admin/server/plugin/schedulePool"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ScheduleApi struct {
}

func refreshSchedulePool() {
	schedules.SPool.Refresh()
}

func validateCronFormat(cronFormat string) bool {
	err := schedules.ValidateCronFormat(cronFormat)
	if err == nil {
		return true
	}
	return false
}

// @Tags Schedule
// @Summary 确认Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body scheduleMdl.Schedule true "Schedule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/schedule/validateScheduleFormat [post]
func (a *ScheduleApi) ValidateScheduleCronFormat(c *gin.Context) {
	var schedule scheduleMdl.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.ScheduleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !validateCronFormat(schedule.CronFormat) {
		response.FailWithMessage("验证失败, 请检查schedule格式", c)
	} else {
		response.OkWithMessage("验证成功", c)
	}
}

// @Tags Schedule
// @Summary 新增Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body scheduleMdl.Schedule true "Schedule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/schedule/addSchedule [post]
func (a *ScheduleApi) AddSchedule(c *gin.Context) {
	var schedule scheduleMdl.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.ScheduleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !validateCronFormat(schedule.CronFormat) {
		response.FailWithMessage("验证失败", c)
		return
	}
	if _, err := scheduleService.CreateSchedule(schedule); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		refreshSchedulePool()
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Schedule
// @Summary 删除Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "ScheduleId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/schedule/deleteSchedule [post]
func (a *ScheduleApi) DeleteSchedule(c *gin.Context) {
	var schedule request.GetById
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := scheduleService.DeleteSchedule(schedule.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		refreshSchedulePool()
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Schedule
// @Summary 批量删除Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/schedule/deleteScheduleByIds [post]
func (a *ScheduleApi) DeleteScheduleByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := scheduleService.DeleteScheduleByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		refreshSchedulePool()
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Schedule
// @Summary 更新Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body scheduleMdl.Schedule true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/schedule/updateSchedule [post]
func (a *ScheduleApi) UpdateSchedule(c *gin.Context) {
	var schedule scheduleMdl.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.ScheduleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !validateCronFormat(schedule.CronFormat) {
		response.FailWithMessage("验证失败", c)
		return
	}
	if err := scheduleService.UpdateSchedule(schedule); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		refreshSchedulePool()
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Schedule
// @Summary 根据id获取Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "EnvironmentId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/schedule/getScheduleById [post]
func (a *ScheduleApi) GetScheduleById(c *gin.Context) {
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
	if schedule, err := scheduleService.GetSchedule(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(scheduleRes.ScheduleResponse{
			Schedule: schedule,
		}, "获取成功", c)
	}
}

// @Tags Schedule
// @Summary 分页获取基础Schedule列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetScheduleByTemplateId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/schedule/getTemplateScheduleList [post]
func (a *ScheduleApi) GetTemplateScheduleList(c *gin.Context) {
	var pageInfo scheduleReq.GetScheduleByTemplateId
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, total := scheduleService.GetScheduleList(int(pageInfo.TemplateId), pageInfo); err != nil {
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
