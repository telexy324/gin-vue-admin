package ansible

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/services/schedules"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type SchedulesApi struct {
}

func refreshSchedulePool(r *http.Request) {
	pool := context.Get(r, "schedule_pool").(schedules.SchedulePool)
	pool.Refresh()
}

// GetSchedule returns single template by ID
func GetSchedule(w http.ResponseWriter, r *http.Request) {
	schedule := context.Get(r, "schedule").(db.Schedule)
	helpers.WriteJSON(w, http.StatusOK, schedule)
}

func GetTemplateSchedules(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	templateID, err := helpers.GetIntParam("template_id", w, r)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "template_id must be provided",
		})
		return
	}

	tplSchedules, err := helpers.Store(r).GetTemplateSchedules(project.ID, templateID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, tplSchedules)
}

func validateCronFormat(cronFormat string, w http.ResponseWriter) bool {
	err := schedules.ValidateCronFormat(cronFormat)
	if err == nil {
		return true
	}
	helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
		"error": "Cron: " + err.Error(),
	})
	return false
}

func ValidateScheduleCronFormat(w http.ResponseWriter, r *http.Request) {
	var schedule db.Schedule
	if !helpers.Bind(w, r, &schedule) {
		return
	}

	_ = validateCronFormat(schedule.CronFormat, w)
}

// AddSchedule adds a template to the database
func AddSchedule(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)

	var schedule db.Schedule
	if !helpers.Bind(w, r, &schedule) {
		return
	}

	if !validateCronFormat(schedule.CronFormat, w) {
		return
	}

	schedule.ProjectID = project.ID
	schedule, err := helpers.Store(r).CreateSchedule(schedule)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)
	objType := db.EventSchedule
	desc := "Schedule ID " + strconv.Itoa(schedule.ID) + " created"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &project.ID,
		ObjectType:  &objType,
		ObjectID:    &schedule.ID,
		Description: &desc,
	})
	if err != nil {
		log.Error(err)
	}

	refreshSchedulePool(r)

	helpers.WriteJSON(w, http.StatusCreated, schedule)
}

// UpdateSchedule writes a schedule to an existing key in the database
func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	oldSchedule := context.Get(r, "schedule").(db.Schedule)

	var schedule db.Schedule
	if !helpers.Bind(w, r, &schedule) {
		return
	}

	// project ID and schedule ID in the body and the path must be the same

	if schedule.ID != oldSchedule.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "schedule id in URL and in body must be the same",
		})
		return
	}

	if schedule.ProjectID != oldSchedule.ProjectID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "You can not move schedule to other project",
		})
		return
	}

	if !validateCronFormat(schedule.CronFormat, w) {
		return
	}

	err := helpers.Store(r).UpdateSchedule(schedule)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	desc := "Schedule ID " + strconv.Itoa(schedule.ID) + " updated"
	objType := db.EventSchedule

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &schedule.ProjectID,
		Description: &desc,
		ObjectID:    &schedule.ID,
		ObjectType:  &objType,
	})

	if err != nil {
		log.Error(err)
	}

	refreshSchedulePool(r)

	w.WriteHeader(http.StatusNoContent)
}

// RemoveSchedule deletes a schedule from the database
func RemoveSchedule(w http.ResponseWriter, r *http.Request) {
	schedule := context.Get(r, "schedule").(db.Schedule)

	err := helpers.Store(r).DeleteSchedule(schedule.ProjectID, schedule.ID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)
	desc := "Schedule ID " + strconv.Itoa(schedule.ID) + " deleted"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &schedule.ProjectID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	refreshSchedulePool(r)

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Schedule
// @Summary 新增Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Schedule true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/schedule/addSchedule [post]
func (a *SchedulesApi) AddSchedule(c *gin.Context) {
	var schedule ansible.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.ScheduleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := scheduleService.CreateSchedule(schedule); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
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
// @Router /ansible/schedule/deleteSchedule [post]
func (a *SchedulesApi) DeleteSchedule(c *gin.Context) {
	var schedule request2.GetByProjectId
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := scheduleService.DeleteSchedule(schedule.ProjectId, schedule.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Schedule
// @Summary 更新Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Schedule true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/schedule/updateSchedule [post]
func (a *SchedulesApi) UpdateSchedule(c *gin.Context) {
	var schedule ansible.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(schedule, utils.ScheduleVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := scheduleService.UpdateSchedule(schedule); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Schedule
// @Summary 根据id获取Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "EnvironmentId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/schedule/getScheduleById [post]
func (a *SchedulesApi) GetScheduleById(c *gin.Context) {
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
	if schedule, err := scheduleService.GetSchedule(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.ScheduleResponse{
			Schedule: schedule,
		}, "获取成功", c)
	}
}

// @Tags Schedule
// @Summary 分页获取基础Schedule列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/schedule/getScheduleList[post]
func (a *SchedulesApi) GetScheduleList(c *gin.Context) {
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
	if err, schedules, total := scheduleService.GetSchedules(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     schedules,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func validateCronFormat(cronFormat string, w http.ResponseWriter) bool {
	err := schedules.ValidateCronFormat(cronFormat)
	if err == nil {
		return true
	}
	helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
		"error": "Cron: " + err.Error(),
	})
	return false
}