package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/songzhibin97/gkit/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskApi struct {
}

// @Tags Task
// @Summary 新增Task模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.TaskTemplate true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addTaskTemplate [post]
func (a *TaskApi) AddTaskTemplate(c *gin.Context) {
	var taskTemplate request2.AddTaskTemplate
	if err := c.ShouldBindJSON(&taskTemplate); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(taskTemplate, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.AddTaskTemplate(*taskTemplate.TaskTemplate); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Task
// @Summary 删除TaskTemplate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "Task id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteTaskTemplate [post]
func (a *TaskApi) DeleteTaskTemplate(c *gin.Context) {
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
	if err := taskService.DeleteTaskTemplate(server.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Task
// @Summary 更新Task模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.TaskTemplate true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateTaskTemplate [post]
func (a *TaskApi) UpdateTaskTemplate(c *gin.Context) {
	var taskTemplate request2.UpdateTaskTemplate
	if err := c.ShouldBindJSON(&taskTemplate); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(taskTemplate, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.UpdateTaskTemplate(*taskTemplate.TaskTemplate); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Task
// @Summary 根据id获取task模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "task id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getTaskTemplateById [post]
func (a *TaskApi) GetTaskTemplateById(c *gin.Context) {
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
	if err, taskTemplate := taskService.GetTaskTemplateById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.TaskTemplateResponse{
			TaskTemplate: taskTemplate,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 分页获取Task模板列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getTaskTemplateList [post]
func (a *TaskApi) GetTaskTemplateList(c *gin.Context) {
	var pageInfo request2.TaskTemplateSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, taskTemplateList, total := taskService.GetTaskTemplateList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     taskTemplateList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 新增Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.Task true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addTaskTemplate [post]
func (a *TaskApi) AddTask(c *gin.Context) {
	var taskTemplate request2.AddTask
	if err := c.ShouldBindJSON(&taskTemplate); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(taskTemplate, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.AddTaskTemplate(*taskTemplate.TaskTemplate); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Task
// @Summary 删除TaskTemplate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "Task id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteTaskTemplate [post]
func (a *TaskApi) DeleteTaskTemplate(c *gin.Context) {
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
	if err := taskService.DeleteTaskTemplate(server.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Task
// @Summary 更新Task模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.TaskTemplate true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateTaskTemplate [post]
func (a *TaskApi) UpdateTaskTemplate(c *gin.Context) {
	var taskTemplate request2.UpdateTaskTemplate
	if err := c.ShouldBindJSON(&taskTemplate); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(taskTemplate, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.UpdateTaskTemplate(*taskTemplate.TaskTemplate); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Task
// @Summary 根据id获取task模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "task id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getTaskTemplateById [post]
func (a *TaskApi) GetTaskTemplateById(c *gin.Context) {
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
	if err, taskTemplate := taskService.GetTaskTemplateById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.TaskTemplateResponse{
			TaskTemplate: taskTemplate,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 分页获取Task模板列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getTaskTemplateList [post]
func (a *TaskApi) GetTaskTemplateList(c *gin.Context) {
	var pageInfo request2.TaskTemplateSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, taskTemplateList, total := taskService.GetTaskTemplateList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     taskTemplateList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetTaskOutput returns the logged task output by id and writes it as json or returns error
func GetTaskOutput(w http.ResponseWriter, r *http.Request) {
	task := context.Get(r, "task").(db.Task)
	project := context.Get(r, "project").(db.Project)

	var output []db.TaskOutput
	output, err := helpers.Store(r).GetTaskOutputs(project.ID, task.ID)

	if err != nil {
		util.LogErrorWithFields(err, log.Fields{"error": "Bad request. Cannot get task output from database"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, output)
}

func StopTask(w http.ResponseWriter, r *http.Request) {
	targetTask := context.Get(r, "task").(db.Task)
	project := context.Get(r, "project").(db.Project)

	if targetTask.ProjectID != project.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := helpers.TaskPool(r).StopTask(targetTask)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveTask removes a task from the database
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	targetTask := context.Get(r, "task").(db.Task)
	editor := context.Get(r, "user").(*db.User)
	project := context.Get(r, "project").(db.Project)

	activeTask := helpers.TaskPool(r).GetTask(targetTask.ID)

	if activeTask != nil {
		// can't delete task in queue or running
		// task must be stopped firstly
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !editor.Admin {
		log.Warn(editor.Username + " is not permitted to delete task logs")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := helpers.Store(r).DeleteTaskWithOutputs(project.ID, targetTask.ID)
	if err != nil {
		util.LogErrorWithFields(err, log.Fields{"error": "Bad request. Cannot delete task from database"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}