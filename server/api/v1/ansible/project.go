package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProjectApi struct {
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
	var project request2.GetById
	if err := c.ShouldBindJSON(&project); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(project, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.DeleteProject(int(project.ID)); err != nil {
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
	var idInfo request2.GetById
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if project, err := projectService.GetProject(int(idInfo.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.ProjectResponse{
			Project: project,
		}, "获取成功", c)
	}
}

// @Tags Project
// @Summary 检验是否为管理员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Project true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/project/mustBeAdmin [post]
func (a *ProjectApi) MustBeAdmin(c *gin.Context) {
	var project request2.GetByProjectId
	if err := c.ShouldBindJSON(&project); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(project, utils.ProjectVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := int(utils.GetUserID(c))
	user, err := userService.GetProjectUser(project.ProjectId, float64(userID))
	if err != nil {
		global.GVA_LOG.Error("获取project管理员失败!", zap.Any("err", err))
		response.FailWithMessage("验证失败", c)
	} else if user.Admin != ansible.IsAdmin {
		response.FailWithMessage("非管理员", c)
	}
}
