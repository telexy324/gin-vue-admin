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

type EnvironmentApi struct {
}

//func GetEnvironmentRefs(w http.ResponseWriter, r *http.Request) {
//	env := context.Get(r, "environment").(db.Environment)
//	refs, err := helpers.Store(r).GetEnvironmentRefs(env.ProjectID, env.ID)
//	if err != nil {
//		helpers.WriteError(w, err)
//		return
//	}
//
//	helpers.WriteJSON(w, http.StatusOK, refs)
//}

// @Tags Environment
// @Summary 新增Environment
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Environment true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/environment/addEnvironment [post]
func (a *EnvironmentApi) AddEnvironment(c *gin.Context) {
	var environment ansible.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.EnvironmentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := environmentService.CreateEnvironment(environment); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Environment
// @Summary 删除Environment
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "EnvronmentId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/environment/deleteEnvironment [post]
func (a *EnvironmentApi) DeleteEnvironment(c *gin.Context) {
	var environment request2.GetByProjectId
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := environmentService.DeleteEnvironment(environment.ProjectId, environment.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Environment
// @Summary 更新Environment
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Environment true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/environment/updateEnvironment [post]
func (a *EnvironmentApi) UpdateEnvironment(c *gin.Context) {
	var environment ansible.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.EnvironmentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := environmentService.UpdateEnvironment(environment); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Environment
// @Summary 根据id获取Environment
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "Environmentid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/environment/getEnvironmentById [post]
func (a *EnvironmentApi) GetEnvironmentById(c *gin.Context) {
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
	if environment, err := environmentService.GetEnvironment(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.EnvironmentResponse{
			Environment: environment,
		}, "获取成功", c)
	}
}

// @Tags Environment
// @Summary 分页获取基础Environment列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/environment/getEnvironmentList[post]
func (a *EnvironmentApi) GetEnvironmentList(c *gin.Context) {
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
	if err, environments, total := environmentService.GetEnvironments(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     environments,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
