package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CmdbServerApi struct {
}

// @Tags CmdbServer
// @Summary 新增服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationServer true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addServer [post]
func (a *CmdbServerApi) AddServer(c *gin.Context) {
	var server request2.AddServer
	if err := c.ShouldBindJSON(&server); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.AddServer(server); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbServer
// @Summary 删除服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteServer [post]
func (a *CmdbServerApi) DeleteServer(c *gin.Context) {
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
	if err := cmdbServerService.DeleteServer(server.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags CmdbServer
// @Summary 更新服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationServer true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateServer [post]
func (a *CmdbServerApi) UpdateServer(c *gin.Context) {
	var server request2.UpdateServer
	if err := c.ShouldBindJSON(&server); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.UpdateServer(server); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags CmdbServer
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerById [post]
func (a *CmdbServerApi) GetServerById(c *gin.Context) {
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
	if err, server := cmdbServerService.GetServerById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationServerResponse{
			Server: server,
		}, "获取成功", c)
	}
}

// @Tags CmdbServer
// @Summary 分页获取基础server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerList [post]
func (a *CmdbServerApi) GetServerList(c *gin.Context) {
	var pageInfo request2.ServerSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, serverList, total := cmdbServerService.GetServerList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     serverList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags CmdbServer
// @Summary 根据系统id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemServers [post]
func (a *CmdbServerApi) GetSystemServers(c *gin.Context) {
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
	if err, serverList := cmdbServerService.GetSystemServers(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationServersResponse{
			Servers: serverList,
		}, "获取成功", c)
	}
}

// @Tags CmdbServer
// @Summary 获取所有服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAllServerIds [get]
func (a *CmdbServerApi) GetAllServerIds(c *gin.Context) {
	err, systemList, _ := cmdbSystemService.GetSystemList(request2.SystemSearch{
		PageInfo: request.PageInfo{
			Page:     1,
			PageSize: 99999,
		},
	})
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	res := make([]applicationRes.AllServersResponse, 0)
	for _, system := range systemList.([]applicationRes.ApplicationSystemResponse) {
		if err, serverList := cmdbServerService.GetSystemServers(float64(system.System.ID)); err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
			response.FailWithMessage("获取失败", c)
		} else {
			resServers := make([]applicationRes.Children, 0)
			for _, server := range serverList {
				resServers = append(resServers, applicationRes.Children{
					ID:   server.ID,
					Name: server.ManageIp,
				})
			}
			res = append(res, applicationRes.AllServersResponse{
				ID:       system.System.ID,
				Name:     system.System.Name,
				Children: resServers,
			})
		}
	}
	response.OkWithDetailed(res, "获取成功", c)
}

// @Tags CmdbServer
// @Summary 新增联系
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.SystemRelation true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/server/addRelation [post]
func (a *CmdbServerApi) AddRelation(c *gin.Context) {
	var relation application.ServerRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(relation, utils.ServerRelationVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.AddRelation(relation); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbServer
// @Summary 获取关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/server/relations [post]
func (a *CmdbServerApi) ServerRelations(c *gin.Context) {
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
	if err, relations, nodes := cmdbServerService.ServerRelations(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		path := applicationRes.RelationPath{}
		resNodes := make([]applicationRes.Node, 0)
		if err = utils.ConvertStruct(&nodes, &resNodes); err != nil {
			response.FailWithMessage("获取失败", c)
		}
		path.Nodes = resNodes
		links := make([]applicationRes.Link, 0)
		mapLinks := make(map[int]bool)
		for _, relation := range relations {
			if mapLinks[int(relation.ID)] == false {
				links = append(links, applicationRes.Link{
					VectorType:     0,
					VectorStrValue: relation.Relation,
					Property: applicationRes.Property{
						Relation:         relation.Relation,
						Url:              relation.EndServerUrl,
						ServerUpdateDate: relation.UpdatedAt.Format("2006-01-02 15:04:05"),
					},
					StartNodeId: relation.StartServerId,
					EndNodeId:   relation.EndServerId,
				})
				mapLinks[int(relation.ID)] = true
			}
		}
		path.Links = links
		response.OkWithDetailed(applicationRes.SystemRelationsResponse{
			Path: path,
		}, "获取成功", c)
	}
}

// @Tags CmdbServer
// @Summary 导出Excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body request2.ExcelInfo true "导出Excel文件信息"
// @Success 200
// @Router /cmdb/exportExcel [post]
func (e *CmdbServerApi) ExportExcel(c *gin.Context) {
	var excelInfo request2.ExcelInfo
	if err := c.ShouldBindJSON(&excelInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	filePath := global.GVA_CONFIG.Excel.Dir + excelInfo.FileName
	err := cmdbServerService.ParseInfoList2Excel(excelInfo.InfoList, excelInfo.Header, filePath)
	if err != nil {
		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
		response.FailWithMessage("转换Excel失败", c)
		return
	}
	c.Writer.Header().Add("success", "true")
	c.File(filePath)
}

// @Tags CmdbServer
// @Summary 导入Excel文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导入成功"}"
// @Router /cmdb/importExcel [post]
func (e *CmdbServerApi) ImportExcel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err = cmdbServerService.ImportExcel2db(file, header)
	if err != nil {
		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
		response.FailWithMessage("转换Excel失败", c)
		return
	}
	response.OkWithMessage("导入成功", c)
}

// @Tags CmdbServer
// @Summary 下载模板
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Success 200
// @Router /cmdb/downloadTemplate [get]
func (e *CmdbServerApi) DownloadTemplate(c *gin.Context) {
	excel, err := cmdbServerService.ExportTemplate()
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
	if err = excel.Write(c.Writer); err != nil {
		global.GVA_LOG.Error("下载模板失败!", zap.Any("err", err))
	}
}

// @Tags CmdbServer
// @Summary 新增应用
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.App true "类型, 名称, 版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addApp [post]
func (a *CmdbServerApi) AddApp(c *gin.Context) {
	var app application.App
	if err := c.ShouldBindJSON(&app); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(app, utils.AppVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.AddApp(app); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbServer
// @Summary 删除应用
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "应用id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteApp [post]
func (a *CmdbServerApi) DeleteApp(c *gin.Context) {
	var app request.GetById
	if err := c.ShouldBindJSON(&app); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(app, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.DeleteApp(app.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags CmdbServer
// @Summary 更新应用
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.App true "类型, 名称, 版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateApp [post]
func (a *CmdbServerApi) UpdateApp(c *gin.Context) {
	var app application.App
	if err := c.ShouldBindJSON(&app); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(app, utils.AppVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbServerService.UpdateApp(app); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags CmdbServer
// @Summary 根据id获取应用
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAppById [post]
func (a *CmdbServerApi) GetAppById(c *gin.Context) {
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
	if err, app := cmdbServerService.GetAppById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationAppResponse{
			App: app,
		}, "获取成功", c)
	}
}

// @Tags CmdbServer
// @Summary 分页获取基础app列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAppList [post]
func (a *CmdbServerApi) GetAppList(c *gin.Context) {
	var pageInfo request2.AppSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, appList, total := cmdbServerService.GetAppList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     appList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
