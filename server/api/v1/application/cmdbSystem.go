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

type CmdbSystemApi struct {
}

// @Tags CmdbSystem
// @Summary 新增系统
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.AddSystem true "系统名, 位置, 管理员id, 主管"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addSystem [post]
func (a *CmdbSystemApi) AddSystem(c *gin.Context) {
	var addSystemRequest request2.AddSystem
	e := c.ShouldBindJSON(&addSystemRequest)
	global.GVA_LOG.Info("error", zap.Any("err", e))
	if err := utils.Verify(addSystemRequest.System, utils.SystemVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbSystemService.AddSystem(addSystemRequest); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 删除系统
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteSystem [post]
func (a *CmdbSystemApi) DeleteSystem(c *gin.Context) {
	var system request.GetById
	_ = c.ShouldBindJSON(&system)
	if err := utils.Verify(system, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbSystemService.DeleteSystem(system.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 更新服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.AddSystem true "系统名, 位置, 管理员id, 主管"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateSystem [post]
func (a *CmdbSystemApi) UpdateSystem(c *gin.Context) {
	var addSystemRequest request2.AddSystem
	e := c.ShouldBindJSON(&addSystemRequest)
	global.GVA_LOG.Info("error", zap.Any("err", e))
	if err := utils.Verify(addSystemRequest.System, utils.SystemVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbSystemService.UpdateSystem(addSystemRequest); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerById [post]
func (a *CmdbSystemApi) GetServerById(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, server := cmdbService.GetServerById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationServerResponse{
			Server: server,
		}, "获取成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 分页获取基础server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerList [post]
func (a *CmdbSystemApi) GetServerList(c *gin.Context) {
	var pageInfo request2.ServerSearch
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, serverList, total := cmdbService.GetServerList(pageInfo); err != nil {
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

// @Tags CmdbSystem
// @Summary 根据系统id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemServers [post]
func (a *CmdbSystemApi) GetSystemServers(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, serverList := cmdbService.GetSystemServers(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationServersResponse{
			Servers: serverList,
		}, "获取成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 新增联系
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.SystemRelation true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/system/addRelation [post]
func (a *CmdbSystemApi) AddRelation(c *gin.Context) {
	var relation application.ServerRelation
	e := c.ShouldBindJSON(&relation)
	global.GVA_LOG.Info("error", zap.Any("err", e))
	if err := utils.Verify(relation, utils.SystemRelationVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbService.AddRelation(relation); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 获取关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/system/relations [post]
func (a *CmdbSystemApi) SystemRelations(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, relations, nodes := cmdbService.ServerRelations(idInfo.ID); err != nil {
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

// @Tags CmdbSystem
// @Summary 导出Excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body request2.ExcelInfo true "导出Excel文件信息"
// @Success 200
// @Router /cmdb/exportExcel [post]
func (e *CmdbSystemApi) ExportExcel(c *gin.Context) {
	var excelInfo request2.ExcelInfo
	_ = c.ShouldBindJSON(&excelInfo)
	filePath := global.GVA_CONFIG.Excel.Dir + excelInfo.FileName
	err := cmdbService.ParseInfoList2Excel(excelInfo.InfoList, excelInfo.Header, filePath)
	if err != nil {
		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
		response.FailWithMessage("转换Excel失败", c)
		return
	}
	c.Writer.Header().Add("success", "true")
	c.File(filePath)
}

// @Tags CmdbSystem
// @Summary 导入Excel文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导入成功"}"
// @Router /cmdb/importExcel [post]
func (e *CmdbSystemApi) ImportExcel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err = cmdbService.ImportExcel2db(file, header)
	if err != nil {
		global.GVA_LOG.Error("转换Excel失败!", zap.Any("err", err))
		response.FailWithMessage("转换Excel失败", c)
		return
	}
	response.OkWithMessage("导入成功", c)
}

// @Tags CmdbSystem
// @Summary 下载模板
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Success 200
// @Router /cmdb/downloadTemplate [get]
func (e *CmdbSystemApi) DownloadTemplate(c *gin.Context) {
	excel, err := cmdbService.ExportTemplate()
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
