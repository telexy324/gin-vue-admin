package ansible

import (
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

	"github.com/gorilla/context"
)

type ProjectApi struct {
}

// ProjectMiddleware ensures a project exists and loads it to the context
func ProjectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*db.User)

		projectID, err := helpers.GetIntParam("project_id", w, r)

		if err != nil {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid project ID",
			})
			return
		}

		// check if user it project's team
		_, err = helpers.Store(r).GetProjectUser(projectID, user.ID)

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		project, err := helpers.Store(r).GetProject(projectID)

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		context.Set(r, "project", project)
		next.ServeHTTP(w, r)
	})
}

// MustBeAdmin ensures that the user has administrator rights
func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := context.Get(r, "project").(db.Project)
		user := context.Get(r, "user").(*db.User)

		projectUser, err := helpers.Store(r).GetProjectUser(project.ID, user.ID)

		if err == db.ErrNotFound {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		if !projectUser.Admin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//GetProject returns a project details
func GetProject(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, context.Get(r, "project"))
}

// UpdateProject saves updated project details to the database
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	var body db.Project

	if !helpers.Bind(w, r, &body) {
		return
	}

	if body.ID != project.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	err := helpers.Store(r).UpdateProject(body)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteProject removes a project from the database
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)

	err := helpers.Store(r).DeleteProject(project.ID)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Project
// @Summary 新增Project
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Project true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/project/addProject [post]
func (a *ProjectApi) AddProject(c *gin.Context) {
	var project ansible.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(project, utils.ProjectVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := projectService.CreateProject(project); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Project
// @Summary 删除Project
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "EnvronmentId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/project/deleteProject [post]
func (a *ProjectApi) DeleteProject(c *gin.Context) {
	var project request2.GetByProjectId
	if err := c.ShouldBindJSON(&project); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(project, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.DeleteProject(project.ProjectId, project.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Project
// @Summary 更新Project
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Project true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/project/updateProject [post]
func (a *ProjectApi) UpdateProject(c *gin.Context) {
	var project ansible.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(project, utils.ProjectVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.UpdateProject(project); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Project
// @Summary 根据id获取Project
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "ProjectId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/project/getProjectById [post]
func (a *ProjectApi) GetProjectById(c *gin.Context) {
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
	if project, err := projectService.GetProject(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.ProjectResponse{
			Project: project,
		}, "获取成功", c)
	}
}

// @Tags Project
// @Summary 分页获取基础Project列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/project/getProjectList[post]
func (a *ProjectApi) GetProjectList(c *gin.Context) {
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
	if err, projects, total := projectService.GetProjects(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     projects,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
