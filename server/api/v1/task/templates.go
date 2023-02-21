package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TemplatesApi struct {
}

//func GetTemplateRefs(w http.ResponseWriter, r *http.Request) {
//	tpl := context.Get(r, "template").(db.Template)
//	refs, err := helpers.Store(r).GetTemplateRefs(tpl.ProjectID, tpl.ID)
//	if err != nil {
//		helpers.WriteError(w, err)
//		return
//	}
//
//	helpers.WriteJSON(w, http.StatusOK, refs)
//}

// @Tags Template
// @Summary 新增Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Template true "Template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/template/addTemplate [post]
func (a *TemplatesApi) AddTemplate(c *gin.Context) {
	var template ansible.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.TemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := templateService.CreateTemplate(template); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Template
// @Summary 删除Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "TemplateId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/template/deleteTemplate [post]
func (a *TemplatesApi) DeleteTemplate(c *gin.Context) {
	var template request.GetByProjectId
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.DeleteTemplate(template.ProjectId, template.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Template
// @Summary 更新Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Template true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/template/updateTemplate [post]
func (a *TemplatesApi) UpdateTemplate(c *gin.Context) {
	var template ansible.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.TemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template.Arguments != nil && *template.Arguments == "" {
		template.Arguments = nil
	}
	if err := templateService.UpdateTemplate(template); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Template
// @Summary 根据id获取Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "TemplateId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/template/getTemplateById [post]
func (a *TemplatesApi) GetTemplateById(c *gin.Context) {
	var idInfo request.GetByProjectId
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template, err := templateService.GetTemplate(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.TemplateResponse{
			Template: template,
		}, "获取成功", c)
	}
}

// @Tags Template
// @Summary 分页获取基础Template列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/template/getTemplateList [post]
func (a *TemplatesApi) GetTemplateList(c *gin.Context) {
	var pageInfo request.GetByProjectId
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, templates, total := templateService.GetTemplates(pageInfo, ansible.TemplateFilter{}); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     templates,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
