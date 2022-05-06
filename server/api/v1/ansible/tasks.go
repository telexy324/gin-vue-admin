package ansible

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/util"
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

type TasksApi struct {
}

// AddTask inserts a task into the database and returns a header or returns error
func AddTask(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	user := context.Get(r, "user").(*db.User)

	var taskObj db.Task

	if !helpers.Bind(w, r, &taskObj) {
		return
	}

	newTask, err := helpers.TaskPool(r).AddTask(taskObj, &user.ID, project.ID)

	if err != nil {
		util.LogErrorWithFields(err, log.Fields{"error": "Cannot write new event to database"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, newTask)
}

// GetTasksList returns a list of tasks for the current project in desc order to limit or error
func GetTasksList(w http.ResponseWriter, r *http.Request, limit uint64) {
	project := context.Get(r, "project").(db.Project)
	tpl := context.Get(r, "template")

	var err error
	var tasks []db.TaskWithTpl

	if tpl != nil {
		tasks, err = helpers.Store(r).GetTemplateTasks(tpl.(db.Template).ProjectID, tpl.(db.Template).ID, db.RetrieveQueryParams{
			Count: int(limit),
		})
	} else {
		tasks, err = helpers.Store(r).GetProjectTasks(project.ID, db.RetrieveQueryParams{
			Count: int(limit),
		})
	}

	if err != nil {
		util.LogErrorWithFields(err, log.Fields{"error": "Bad request. Cannot get tasks list from database"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, tasks)
}

// GetAllTasks returns all tasks for the current project
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	GetTasksList(w, r, 0)
}

// GetLastTasks returns the hundred most recent tasks
func GetLastTasks(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(str)
	if err != nil || limit <= 0 || limit > 200 {
		limit = 200
	}
	GetTasksList(w, r, uint64(limit))
}

// GetTask returns a task based on its id
func GetTask(w http.ResponseWriter, r *http.Request) {
	task := context.Get(r, "task").(db.Task)
	helpers.WriteJSON(w, http.StatusOK, task)
}

// GetTaskMiddleware is middleware that gets a task by id and sets the context to it or panics
func GetTaskMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := context.Get(r, "project").(db.Project)
		taskID, err := helpers.GetIntParam("task_id", w, r)

		if err != nil {
			util.LogErrorWithFields(err, log.Fields{"error": "Bad request. Cannot get task_id from request"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task, err := helpers.Store(r).GetTask(project.ID, taskID)
		if err != nil {
			util.LogErrorWithFields(err, log.Fields{"error": "Bad request. Cannot get task from database"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		context.Set(r, "task", task)
		next.ServeHTTP(w, r)
	})
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

// @Tags Task
// @Summary 新增Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Task true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/task/addTask [post]
func (a *TasksApi) AddTask(c *gin.Context) {
	var task ansible.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(task, utils.TaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := taskService.CreateTask(task); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Task
// @Summary 删除Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/task/deleteTask [post]
func (a *TasksApi) DeleteTask(c *gin.Context) {
	var task request2.GetByProjectId
	if err := c.ShouldBindJSON(&task); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(task, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.DeleteTaskWithOutputs(task.ProjectId, task.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Task
// @Summary 更新Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Task true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/task/updateTask [post]
func (a *TasksApi) UpdateTask(c *gin.Context) {
	var task ansible.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(task, utils.TaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.UpdateTask(task); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Task
// @Summary 根据id获取Task
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "TaskId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/task/getTaskById [post]
func (a *TasksApi) GetTaskById(c *gin.Context) {
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
	if task, err := taskService.GetTask(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.TaskResponse{
			Task: task,
		}, "获取成功", c)
	}
}

// @Tags Task
// @Summary 分页获取基础Task列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/task/getTaskList[post]
func (a *TasksApi) GetTaskList(c *gin.Context) {
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
	if err, tasks, total := taskService.GetTemplateTasks(pageInfo); err != nil {
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
