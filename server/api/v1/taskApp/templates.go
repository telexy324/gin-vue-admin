package taskApp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	templateReq "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	templateRes "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/taskPool"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strconv"
	"strings"
	"sync"
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
// @Summary 批量删除Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteTemplateByIds [post]
func (a *TemplateApi) DeleteTemplateByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := templateService.DeleteTaskTemplateByIds(ids); err != nil {
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
	var pageInfo templateReq.TaskTemplateSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminID := utils.GetUserID(c)
	err, adminSystems := cmdbSystemService.GetAdminSystems(adminID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	if len(pageInfo.SystemIDs) > 0 {
		copedIDs := make([]int, 0, 8)
		for _, adminSystem := range adminSystems {
			for _, id := range pageInfo.SystemIDs {
				if id == int(adminSystem.ID) {
					copedIDs = append(copedIDs, id)
					continue
				}
			}
		}
		pageInfo.SystemIDs = copedIDs
	} else {
		for _, adminSystem := range adminSystems {
			pageInfo.SystemIDs = append(pageInfo.SystemIDs, int(adminSystem.ID))
		}
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
	//err, server := cmdbServerService.GetServerById(info.ServerId)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	failedChan := make(chan string, len(template.TargetServers))
	wg := &sync.WaitGroup{}
	for _, server := range template.TargetServers {
		wg.Add(1)
		go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
			var sshClient common.SSHClient
			defer w.Done()
			defer func() {
				if sshClient.Client != nil {
					sshClient.Client.Close()
				}
			}()
			sshClient, err := common.FillSSHClient(s.ManageIp, template.SysUser, "", s.SshPort)
			err = sshClient.GenerateClient()
			if err != nil {
				global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				f <- s.ManageIp
				return
			}
			if exist, err := templateService.CheckScript(s, &sshClient, template); err != nil || !exist {
				f <- s.ManageIp
				return
			}
			return
		}(wg, server, failedChan)
	}
	wg.Wait()
	close(failedChan)
	failedIPs := make([]string, 0)
	for ip := range failedChan {
		failedIPs = append(failedIPs, ip)
	}
	if len(failedIPs) == 0 && info.Detail {
		sshClient, err := common.FillSSHClient(template.TargetServers[0].ManageIp, template.SysUser, "", template.TargetServers[0].SshPort)
		err = sshClient.GenerateClient()
		if err != nil {
			global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", template.TargetServers[0].ManageIp), zap.Any("err", err))
			response.FailWithMessage(err.Error(), c)
			return
		}
		if output, err := templateService.CheckScriptDetail(&sshClient, template); err != nil {
			sshClient.Client.Close()
			response.FailWithMessage(err.Error(), c)
			return
		} else {
			sshClient.Client.Close()
			response.OkWithDetailed(templateRes.CheckScriptResponse{
				FailedIps: failedIPs,
				Script:    output,
			}, "获取成功", c)
			return
		}
	}
	response.OkWithDetailed(templateRes.CheckScriptResponse{
		FailedIps: failedIPs,
	}, "获取成功", c)
}

// @Tags Template
// @Summary 下载script
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.TemplateScriptRequest true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"下载成功"}"
// @Router /task/template/downloadScript [get]
func (a *TemplateApi) DownloadScript(c *gin.Context) {
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
	fileBytes, err := templateService.DownloadScript(info.ID, server)

	//c.Writer.Header().Add("success", "true")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"serverTemplate.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("success", "true")
	if _, err = c.Writer.Write(fileBytes); err != nil {
		global.GVA_LOG.Error("下载脚本失败!", zap.Any("err", err))
	}
}

// @Tags Template
// @Summary 导入脚本
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Script文件"
// @Param int query int false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导入成功"}"
// @Router /task/template/uploadScript [post]
func (a *TemplateApi) UploadScript(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	idStr := c.Request.FormValue("ID")
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	userID := utils.GetUserID(c)
	scriptPath := c.Request.FormValue("scriptPath")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	tmpl, err := templateService.GetTaskTemplate(float64(ID))
	if err != nil {
		global.GVA_LOG.Error("获取模板失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		global.GVA_LOG.Error("计算文件md5失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	tmpl.ScriptHash = hex.EncodeToString(hash.Sum(nil))
	if err = templateService.UpdateTaskTemplate(tmpl); err != nil {
		global.GVA_LOG.Error("更新模板脚本md5失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file.Seek(0, io.SeekStart)
	failedIps, err := templateService.UploadScript(ID, file, scriptPath, userID)
	if err != nil {
		global.GVA_LOG.Error("上传脚本失败!", zap.Any("err", err))
		response.FailWithMessage("上传脚本失败", c)
		return
	}
	response.OkWithDetailed(templateRes.UploadScriptResponse{
		FailedIps: failedIps,
	}, "获取成功", c)
}

// @Tags Template
// @Summary 新增模板集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.AddSet true "系统名, 位置, 管理员id, 主管"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/template/addSet [post]
func (a *TemplateApi) AddSet(c *gin.Context) {
	var addSetRequest templateReq.AddSet
	if err := c.ShouldBindJSON(&addSetRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(addSetRequest, utils.TaskTemplateSetVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.AddSet(addSetRequest); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Template
// @Summary 删除模板集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteSet [post]
func (a *TemplateApi) DeleteSet(c *gin.Context) {
	var system request.GetById
	if err := c.ShouldBindJSON(&system); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(system, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.DeleteSet(system.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Template
// @Summary 批量删除模板集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteSetByIds [post]
func (a *TemplateApi) DeleteSetByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := templateService.DeleteSetByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Template
// @Summary 更新模板集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.AddSet true "系统名, 位置, 管理员id, 主管"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/updateSet [post]
func (a *TemplateApi) UpdateSet(c *gin.Context) {
	var addSetRequest templateReq.AddSet
	if err := c.ShouldBindJSON(&addSetRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(addSetRequest, utils.TaskTemplateSetVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := templateService.UpdateSet(addSetRequest); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Template
// @Summary 根据id获取模板集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getSetById [post]
func (a *TemplateApi) GetSetById(c *gin.Context) {
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
	err, set, templates := templateService.GetSetById(idInfo.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	response.OkWithDetailed(templateRes.TaskTemplateSetResponse{
		TaskTemplateSet: set,
		Templates:       templates,
	}, "获取成功", c)
}

// @Tags Template
// @Summary 分页获取基础模板集列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getSetList [post]
func (a *TemplateApi) GetSetList(c *gin.Context) {
	var pageInfo templateReq.TaskTemplateSetSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminID := utils.GetUserID(c)
	err, adminSystems := cmdbSystemService.GetAdminSystems(adminID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	if len(pageInfo.SystemIDs) > 0 {
		copedIDs := make([]int, 0, 8)
		for _, adminSystem := range adminSystems {
			for _, id := range pageInfo.SystemIDs {
				if id == int(adminSystem.ID) {
					copedIDs = append(copedIDs, id)
					continue
				}
			}
		}
		pageInfo.SystemIDs = copedIDs
	} else {
		for _, adminSystem := range adminSystems {
			pageInfo.SystemIDs = append(pageInfo.SystemIDs, int(adminSystem.ID))
		}
	}
	if err, systemList, total := templateService.GetSetList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     systemList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags Template
// @Summary 新增模板集任务集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.AddSet true "系统名, 位置, 管理员id, 主管"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /task/template/addSetTask [post]
func (a *TemplateApi) AddSetTask(c *gin.Context) {
	var addSetTaskRequest taskMdl.SetTask
	if err := c.ShouldBindJSON(&addSetTaskRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(addSetTaskRequest, utils.TaskTemplateSetTaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	addSetTaskRequest.SystemUserId = int(utils.GetUserID(c))
	if err := templateService.AddSetTask(addSetTaskRequest); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Template
// @Summary 更新模板集任务集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.ProcessTaskRequest true "模板集id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/processSetTask [post]
func (a *TemplateApi) ProcessSetTask(c *gin.Context) {
	var processTaskRequest templateReq.ProcessTaskRequest
	if err := c.ShouldBindJSON(&processTaskRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(processTaskRequest, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, setTask := templateService.GetSetTaskById(processTaskRequest.ID)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
		return
	}
	if setTask.Tasks != nil && len(setTask.Tasks) > 0 && setTask.Tasks[len(setTask.Tasks)-1].Status != taskMdl.TaskSuccessStatus {
		global.GVA_LOG.Error("任务未结束或存在异常!", zap.Any("err", err))
		response.FailWithMessage("任务未结束或存在异常!", c)
		return
	}
	if setTask.TotalSteps == setTask.CurrentStep {
		global.GVA_LOG.Error("任务已结束!", zap.Any("err", err))
		response.FailWithMessage("任务已结束!", c)
		return
	}
	userID := int(utils.GetUserID(c))
	if setTask.SystemUserId != userID {
		global.GVA_LOG.Error("非创建人!", zap.Any("err", err))
		response.FailWithMessage("非创建人!", c)
		return
	}
	var task taskMdl.Task
	task.TemplateId = int(setTask.Templates[setTask.CurrentStep].ID)
	task.CommandVars = processTaskRequest.CommandVars
	newTask, err := taskPool.TPool.AddTask(task, userID, int(setTask.ID))
	if err != nil {
		global.GVA_LOG.Error("添加任务失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
		return
	}
	setTask.CurrentStep += 1
	setTask.Tasks = append(setTask.Tasks, newTask)
	taskIds := make([]int, 0, len(setTask.Tasks))
	for _, t := range setTask.Tasks {
		taskIds = append(taskIds, int(t.ID))
	}
	s, err := json.Marshal(taskIds)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
		return
	}
	setTask.TasksString = string(s)
	setTask.CurrentTaskId = int(newTask.ID)
	if err = templateService.UpdateSetTask(setTask); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithDetailed(templateRes.TaskResponse{
			Task: newTask,
		}, "添加成功", c)
	}
}

// @Tags Template
// @Summary 根据id获取模板集任务集
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "模板集id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getSetTaskById [post]
func (a *TemplateApi) GetSetTaskById(c *gin.Context) {
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
	if err, setTask := templateService.GetSetTaskById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithDetailed(setTask, "获取成功", c)
	}
}

// @Tags Template
// @Summary 分页获取模板集任务集列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getSetTaskList [post]
func (a *TemplateApi) GetSetTaskList(c *gin.Context) {
	var pageInfo templateReq.SetTaskSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo, utils.TaskTemplateSetTaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, setTaskList, total := templateService.GetSetTaskList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     setTaskList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags Template
// @Summary 获取文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "模板id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/getFileList [post]
func (a *TemplateApi) GetFileList(c *gin.Context) {
	var idInfo templateReq.FileListRequest
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.TaskFileListVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(idInfo.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template.TargetIds == nil || len(template.TargetIds) == 0 {
		response.FailWithMessage("template target ids is null", c)
	}
	if template.ExecuteType != consts.ExecuteTypeDownload {
		response.FailWithMessage("template type is not download", c)
	}
	err, server := cmdbServerService.GetServerById(float64(template.TargetIds[0]))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sshClient, err := common.FillSSHClient(server.ManageIp, template.SysUser, "", server.SshPort)
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		return
	}
	fileInfos, isTop, err := templateService.GetFileList(&sshClient, template, idInfo.Directory)
	if err != nil {
		global.GVA_LOG.Error("get file list failed", zap.Any("err", err))
		response.FailWithMessage("get file list failed", c)
	} else {
		response.OkWithDetailed(templateRes.TemplateFileListResponse{
			FileInfos:        fileInfos,
			IsTop:            isTop,
			CurrentDirectory: strings.TrimRight(idInfo.Directory, "/") + "/",
		}, "获取成功", c)
	}
}

// @Tags Template
// @Summary 下载文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.DownLoadFileRequest true "id,文件路径"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/downloadFile [get]
func (a *TemplateApi) DownloadFile(c *gin.Context) {
	var info templateReq.DownLoadFileRequest
	if err := c.ShouldBindQuery(&info); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(info, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template.TargetIds == nil || len(template.TargetIds) == 0 {
		response.FailWithMessage("template target ids is null", c)
		return
	}
	if template.ExecuteType != consts.ExecuteTypeDownload {
		response.FailWithMessage("template type is not download", c)
		return
	}
	if !strings.Contains(info.File, template.LogPath) {
		response.FailWithMessage("file not in log path", c)
		return
	}
	err, server := cmdbServerService.GetServerById(float64(template.TargetIds[0]))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sshClient, err := common.FillSSHClient(server.ManageIp, template.SysUser, "", server.SshPort)
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		response.FailWithMessage("create ssh client failed", c)
		return
	}
	fileBytes, err := sshClient.Download(info.File)
	if err != nil {
		global.GVA_LOG.Error("download file failed", zap.Any("err", err))
		response.FailWithMessage("download file failed", c)
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+info.File[strings.LastIndex(info.File, "/")+1:])
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("success", "true")
	//fileBytes,err:=io.ReadAll(file)
	//if err != nil {
	//	global.GVA_LOG.Error("read file failed", zap.Any("err", err))
	//	response.FailWithMessage("read file failed", c)
	//	return
	//}
	//if _, err = c.Writer.Write(fileBytes); err != nil {
	//	global.GVA_LOG.Error("下载文件失败!", zap.Any("err", err))
	//}
	if _, err = c.Writer.Write(fileBytes); err != nil {
		global.GVA_LOG.Error("下载文件失败!", zap.Any("err", err))
		response.FailWithMessage("download file failed", c)
	}
}

// @Tags Template
// @Summary 上传日志服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.DownLoadFileRequest true "id,文件路径"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/uploadLogServer [post]
func (a *TemplateApi) UploadLogServer(c *gin.Context) {
	var info templateReq.DownLoadFileRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(info, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template.TargetIds == nil || len(template.TargetIds) == 0 {
		response.FailWithMessage("template target ids is null", c)
		return
	}
	if template.ExecuteType != consts.ExecuteTypeDownload {
		response.FailWithMessage("template type is not download", c)
		return
	}
	if template.LogOutput != consts.LogOutputTypeUpload {
		response.FailWithMessage("template type is not upload to log server", c)
		return
	}
	if !strings.Contains(info.File, template.LogPath) {
		response.FailWithMessage("file not in log path", c)
		return
	}
	userID := int(utils.GetUserID(c))
	task := taskMdl.Task{
		TemplateId:   int(info.ID),
		FileDownload: info.File,
	}
	if taskNew, err := taskPool.TPool.AddTask(task, userID, 0); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))
		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithDetailed(templateRes.TaskResponse{
			Task: taskNew,
		}, "添加成功", c)
	}
	//err, server := cmdbServerService.GetServerById(float64(template.TargetIds[0]))
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//sshClient, err := common.FillSSHClient(server.ManageIp, template.SysUser, "", server.SshPort)
	//err = sshClient.GenerateClient()
	//if err != nil {
	//	global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//filePath := "/" + strings.Trim(template.LogPath, "/") + "/"
	//fileReader, err := sshClient.DownloadReader(filePath + info.File)
	//if err != nil {
	//	global.GVA_LOG.Error("download file failed", zap.Any("err", err))
	//	response.FailWithMessage("download file failed", c)
	//	return
	//}
	//err, logServer := logUploadServerService.GetServerById(float64(template.DstServerId))
	//if err != nil {
	//	global.GVA_LOG.Error("get log server failed", zap.Any("err", err))
	//	response.FailWithMessage("get log server failed", c)
	//	return
	//}
	//err, secret := logUploadSecretService.GetSecretById(float64(template.SecretId))
	//if err != nil {
	//	global.GVA_LOG.Error("get log server secret failed", zap.Any("err", err))
	//	response.FailWithMessage("get log server secret failed", c)
	//	return
	//}
	//filePathUpload := "/" + strings.Trim(template.LogDst, "/") + "/"
	//if logServer.Mode == consts.LogServerModeFtp {
	//	ftpClient, err := common.NewFtpClient(logServer.ManageIp, logServer.Port, secret.Name, secret.Password)
	//	if err != nil {
	//		global.GVA_LOG.Error("create ftp client failed", zap.Any("err", err))
	//		response.FailWithMessage("create ftp client failed", c)
	//		return
	//	}
	//	defer ftpClient.Conn.Quit()
	//	if err = ftpClient.Upload(filePathUpload+info.File, fileReader); err != nil {
	//		global.GVA_LOG.Error("upload via ftp failed", zap.Any("err", err))
	//		response.FailWithMessage("upload via ftp failed", c)
	//		return
	//	} else {
	//		response.OkWithMessage("上传成功", c)
	//		return
	//	}
	//} else if logServer.Mode == consts.LogServerModeSSH {
	//	sshClientUpload, err := common.FillSSHClient(logServer.ManageIp, secret.Name, secret.Password, logServer.Port)
	//	err = sshClientUpload.GenerateClient()
	//	if err != nil {
	//		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
	//		response.FailWithMessage(err.Error(), c)
	//		return
	//	}
	//	if err = sshClientUpload.Upload(fileReader, filePathUpload+info.File); err != nil {
	//		global.GVA_LOG.Error("upload via sftp failed", zap.Any("err", err))
	//		response.FailWithMessage("upload via sftp failed", c)
	//		return
	//	} else {
	//		response.OkWithMessage("上传成功", c)
	//		return
	//	}
	//} else {
	//	response.FailWithMessage("上传失败", c)
	//	return
	//}
}

// @Tags Template
// @Summary 上传日志服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body templateReq.DownLoadFileRequest true "id,文件路径"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/template/deployServer [post]
func (a *TemplateApi) DeployServer(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(info, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTaskTemplate(info.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if template.TargetIds == nil || len(template.TargetIds) == 0 {
		response.FailWithMessage("template target ids is null", c)
		return
	}
	if template.ExecuteType != consts.ExecuteTypeDeploy {
		response.FailWithMessage("template type is not download", c)
		return
	}
	userID := int(utils.GetUserID(c))
	task := taskMdl.Task{
		TemplateId: int(info.ID),
	}
	if taskNew, err := taskPool.TPool.AddTask(task, userID, 0); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))
		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithDetailed(templateRes.TaskResponse{
			Task: taskNew,
		}, "添加成功", c)
	}
}
