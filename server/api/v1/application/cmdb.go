package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CmdbApi struct {
}

// @Tags Server
// @Summary 新增服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationServer true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addServer [post]
func (a *CmdbApi) AddServer(c *gin.Context) {
	var server application.ApplicationServer
	e := c.ShouldBindJSON(&server)
	global.GVA_LOG.Info("error", zap.Any("err", e))
	if err := utils.Verify(server, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbService.AddServer(server); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Server
// @Summary 删除服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteServer [post]
func (a *CmdbApi) DeleteServer(c *gin.Context) {
	var server request.GetById
	_ = c.ShouldBindJSON(&server)
	if err := utils.Verify(server, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbService.DeleteServer(server.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Server
// @Summary 更新服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationServer true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateServer [post]
func (a *CmdbApi) UpdateServer(c *gin.Context) {
	var server application.ApplicationServer
	_ = c.ShouldBindJSON(&server)
	if err := utils.Verify(server, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbService.UpdateServer(server); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Server
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerById [post]
func (a *CmdbApi) GetServerById(c *gin.Context) {
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
		servers := make([]application.ApplicationServer, 0)
		servers = append(servers, server)
		response.OkWithDetailed(applicationRes.ApplicationServerResponse{
			Servers: servers,
		}, "获取成功", c)
	}
}

// @Tags Server
// @Summary 分页获取基础server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getServerList [post]
func (a *CmdbApi) GetServerList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, serverList, total := cmdbService.GetServerList(); err != nil {
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

// @Tags Server
// @Summary 新增联系
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.SystemRelation true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/system/addRelation [post]
func (a *CmdbApi) AddRelation(c *gin.Context) {
	var relation application.SystemRelation
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

// @Tags Server
// @Summary 获取关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/system/relations [post]
func (a *CmdbApi) SystemRelations(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, relations, nodes := cmdbService.SystemRelations(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		paths := applicationRes.RelationPath{}
		resNodes := make([]applicationRes.Node, 0)
		if err = utils.ConvertStruct(&nodes, &resNodes); err != nil {
			response.FailWithMessage("获取失败", c)
		}
		paths.Nodes = resNodes
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
		paths.Links = links
		response.OkWithDetailed(applicationRes.SystemRelationsResponse{
			Paths: paths,
		}, "获取成功", c)
	}
}
