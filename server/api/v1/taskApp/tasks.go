package taskApp

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	taskReq "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	taskRes "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/taskPool"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskApi struct {
}

//// GetAllTasks returns all tasks for the current project
//func GetAllTasks(w http.ResponseWriter, r *http.Request) {
//	GetTasksList(w, r, 0)
//}

//// GetLastTasks returns the hundred most recent tasks
//func GetLastTasks(w http.ResponseWriter, r *http.Request) {
//	str := r.URL.Query().Get("limit")
//	limit, err := strconv.Atoi(str)
//	if err != nil || limit <= 0 || limit > 200 {
//		limit = 200
//	}
//	GetTasksList(w, r, uint64(limit))
//}

// @Tags Task
// @Summary 新增Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body taskMdl.Task true "Task"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/addTask [post]
func (a *TaskApi) AddTask(c *gin.Context) {
	var task taskMdl.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(task, utils.TaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := int(utils.GetUserID(c))
	if task, err := taskPool.TPool.AddTask(task, userID); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithDetailed(taskRes.TaskResponse{
			Task: task,
		}, "添加成功", c)
	}
}

// @Tags Task
// @Summary 删除Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/deleteTask [post]
func (a *TaskApi) DeleteTask(c *gin.Context) {
	var taskRequest request.GetById
	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(taskRequest, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	activeTask := taskPool.TPool.GetTask(int(taskRequest.ID))
	if activeTask != nil {
		response.FailWithMessage("task正在执行", c)
		return
	}
	//todo 非超级管理员不可删除
	//userID := utils.GetUserID(c)
	//user, err := userService.GetProjectUser(taskRequest.ProjectId, float64(userID))
	//if err != nil {
	//	global.GVA_LOG.Error("获取task管理员失败!", zap.Any("err", err))
	//	response.FailWithMessage("删除失败", c)
	//} else if user.Admin != ansible.IsAdmin {
	//	response.FailWithMessage("非管理员", c)
	//}
	if err := taskService.DeleteTaskWithOutputs(int(taskRequest.ID)); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

//// @Tags Task
//// @Summary 更新Task
//// @Security ApiKeyAuth
//// @accept application/json
//// @Produce application/json
//// @Param data body ansible.Task true "主机名, 架构, 管理ip, 系统, 系统版本"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
//// @Router /ansible/task/updateTask [post]
//func (a *TaskApi) UpdateTask(c *gin.Context) {
//	var task taskMdl.Task
//	if err := c.ShouldBindJSON(&task); err != nil {
//		global.GVA_LOG.Info("error", zap.Any("err", err))
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	if err := utils.Verify(task, utils.TaskVerify); err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	if err := taskService.UpdateTask(task); err != nil {
//		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
//		response.FailWithMessage("更新失败", c)
//	} else {
//		response.OkWithMessage("更新成功", c)
//	}
//}

// @Tags Task
// @Summary 根据id获取Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/getTaskById [post]
func (a *TaskApi) GetTaskById(c *gin.Context) {
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
	if task, err := taskService.GetTask(int(idInfo.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(taskRes.TaskResponse{
			Task: task,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 分页获取基础Task列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetTaskByTemplateId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/getTaskList [post]
func (a *TaskApi) GetTaskList(c *gin.Context) {
	var pageInfo taskReq.GetTaskByTemplateId
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.TemplateId >= 0 {
		if err, tasks, total := taskService.GetTemplateTasks(int(pageInfo.TemplateId), pageInfo.PageInfo); err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
			response.FailWithMessage("获取失败", c)
		} else {
			response.OkWithDetailed(response.PageResult{
				List:     tasks,
				Total:    total,
				Page:     pageInfo.Page,
				PageSize: pageInfo.PageSize,
			}, "获取成功", c)
		}
	} else {
		if err, tasks, total := taskService.GetProjectTasks(pageInfo.PageInfo); err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
			response.FailWithMessage("获取失败", c)
		} else {
			response.OkWithDetailed(response.PageResult{
				List:     tasks,
				Total:    total,
				Page:     pageInfo.Page,
				PageSize: pageInfo.PageSize,
			}, "获取成功", c)
		}
	}
}

// @Tags Task
// @Summary 根据id获取TaskOutputs
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/getTaskOutputs [post]
func (a *TaskApi) GetTaskOutputs(c *gin.Context) {
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
	//task, err := taskService.GetTask(int(idInfo.ID))
	//if err != nil {
	//	global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
	//	response.FailWithMessage("获取失败", c)
	//	return
	//}
	if taskOutput, err := taskService.GetTaskOutputs(int(idInfo.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(taskRes.TaskOutputsResponse{
			TaskOutputs: taskOutput,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 停止Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/task/stopTask [post]
func (a *TaskApi) StopTask(c *gin.Context) {
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
	task, err := taskService.GetTask(int(idInfo.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	err = taskPool.TPool.StopTask(task)
	if err != nil {
		global.GVA_LOG.Error("停止失败!", zap.Any("err", err))
		response.FailWithMessage("停止失败", c)
	} else {
		response.OkWithMessage("停止成功", c)
	}
}
