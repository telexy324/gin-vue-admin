package taskApp

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	templateReq "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	templateRes "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TemplateApi struct {
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
// @Param data body taskMdl.TaskTemplate true "Template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/template/addTemplate [post]
func (a *TemplateApi) AddTemplate(c *gin.Context) {
	var template taskMdl.TaskTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := templateService.CreateTaskTemplate(template); err != nil {
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
// @Param data body request.GetById true "TaskTemplateId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteTemplate [post]
func (a *TemplateApi) DeleteTemplate(c *gin.Context) {
	var template request.GetById
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.DeleteTaskTemplate(template.ID); err != nil {
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
// @Param data body taskMdl.TaskTemplate true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/updateTemplate [post]
func (a *TemplateApi) UpdateTemplate(c *gin.Context) {
	var template taskMdl.TaskTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(template, utils.TaskTemplateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.UpdateTaskTemplate(template); err != nil {
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
// @Param data body request.GetById true "TemplateId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getTemplateById [post]
func (a *TemplateApi) GetTemplateById(c *gin.Context) {
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
	if template, err := templateService.GetTaskTemplate(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(templateRes.TaskTemplateResponse{
			Template: template,
		}, "获取成功", c)
	}
}

// @Tags Template
// @Summary 分页获取基础Template列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getTemplateList [post]
func (a *TemplateApi) GetTemplateList(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, templates, total := templateService.GetTaskTemplates(pageInfo); err != nil {
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

//// @Tags Template
//// @Summary 上传脚本
//// @Security ApiKeyAuth
//// @accept multipart/form-data
//// @Produce  application/json
//// @Param file formData file true "上传脚本"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"导入成功"}"
//// @Router /task/template/uploadScript [post]
//func (a *TemplateApi) UploadScript(c *gin.Context) {
//	file, header, err := c.Request.FormFile("file")
//	if err != nil {
//		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
//		response.FailWithMessage("接收文件失败", c)
//		return
//	}
//	err = templateService.UploadScript(file, header)
//	if err != nil {
//		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
//		response.FailWithMessage("转换Excel失败", c)
//		return
//	}
//	response.OkWithMessage("导入成功", c)
//}

// @Tags Template
// @Summary 检查script
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.TemplateScriptRequest true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/checkScript [post]
func (a *TemplateApi) CheckScript(c *gin.Context) {
	var info templateReq.TemplateScriptRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	//if err := utils.Verify(info.ID, utils.IdVerify); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	err, server := cmdbServerService.GetServerById(info.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sshClient, err := sshService.FillSSHClient(server.ManageIp, template.SysUser, "123456", server.SshPort)
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		return
	}
	exist, output, err := templateService.CheckScript(server, info.Detail, &sshClient, template)
	if err != nil {
		global.GVA_LOG.Error("check script failed", zap.Any("err", err))
		response.FailWithMessage("check script failed", c)
	}
	response.OkWithDetailed(templateRes.TemplateScriptResponse{
		Exist:  exist,
		Script: output,
	}, "获取成功", c)
}

// @Tags Template
// @Summary 检查script
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.TemplateScriptRequest true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"下载成功"}"
// @Router /task/template/downloadTemplate [get]
func (a *TemplateApi) DownloadTemplate(c *gin.Context) {
	var info templateReq.TemplateScriptRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, server := cmdbServerService.GetServerById(info.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sshClient, err := sshService.FillSSHClient(server.ManageIp, template.SysUser, "123456", server.SshPort)
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		return
	}
	file, err := sshClient.Download(template.ScriptPath)
	if err != nil {
		global.GVA_LOG.Error("下载模板失败!", zap.Any("err", err))
		response.FailWithMessage("下载模板失败", c)
		return
	}
	//c.Writer.Header().Add("success", "true")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"serverTemplate.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("success", "true")
	if _, err = file.WriteTo(c.Writer); err != nil {
		global.GVA_LOG.Error("下载脚本失败!", zap.Any("err", err))
	}
}
