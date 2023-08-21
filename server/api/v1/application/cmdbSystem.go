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
	if err := c.ShouldBindJSON(&addSystemRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(addSystemRequest, utils.SystemVerify); err != nil {
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
	if err := c.ShouldBindJSON(&system); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
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
// @Summary 批量删除系统
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteSystemByIds [post]
func (a *CmdbSystemApi) DeleteSystemByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := cmdbSystemService.DeleteSystemByIds(ids); err != nil {
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
	if err := c.ShouldBindJSON(&addSystemRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(addSystemRequest, utils.SystemVerify); err != nil {
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
// @Summary 根据id获取系统
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemById [post]
func (a *CmdbSystemApi) GetSystemById(c *gin.Context) {
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
	err, system, admins, adminIds := cmdbSystemService.GetSystemById(idInfo.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	}
	adminInfos := make([]application.Admin, 0, len(admins))
	for _, admin := range admins {
		if err, admin := staffService.GetAdminById(float64(admin.ID)); err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
			response.FailWithMessage("获取失败", c)
		} else {
			adminInfos = append(adminInfos, admin)
		}
	}
	response.OkWithDetailed(applicationRes.ApplicationSystemResponse{
		System:   system,
		Admins:   adminInfos,
		AdminIds: adminIds,
	}, "获取成功", c)
}

// @Tags CmdbSystem
// @Summary 分页获取基础server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemList [post]
func (a *CmdbSystemApi) GetSystemList(c *gin.Context) {
	var pageInfo request2.SystemSearch
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
	if err, systemList, total := cmdbSystemService.GetSystemList(pageInfo, adminID); err != nil {
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

// @Tags CmdbSystem
// @Summary 获取管理员所有系统
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Empty true "空"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAdminSystems [post]
func (a *CmdbSystemApi) GetAdminSystems(c *gin.Context) {
	adminID := utils.GetUserID(c)
	if err, systemList := cmdbSystemService.GetAdminSystems(adminID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.ApplicationSystemsResponse{
			Systems: systemList,
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
	var relation application.SystemRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(relation, utils.SystemRelationVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbSystemService.AddRelation(relation); err != nil {
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
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/system/relations [post]
func (a *CmdbSystemApi) SystemRelations(c *gin.Context) {
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
	if err, relations, nodes := cmdbSystemService.SystemRelations(idInfo.ID); err != nil {
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
						Url:              relation.EndSystemUrl,
						ServerUpdateDate: relation.UpdatedAt.Format("2006-01-02 15:04:05"),
					},
					StartNodeId: relation.StartSystemId,
					EndNodeId:   relation.EndSystemId,
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
// @Summary 新增编辑器关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationSystemEditRelation true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/system/addEditRelation [post]
func (a *CmdbSystemApi) AddEditRelation(c *gin.Context) {
	var relation application.ApplicationSystemEditRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(relation, utils.SystemEditRelationVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cmdbSystemService.AddEditRelations(relation); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 删除编辑器关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/system/deleteEditRelation [post]
func (a *CmdbSystemApi) DeleteEditRelation(c *gin.Context) {
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
	if err := cmdbSystemService.DeleteEditRelations(idInfo.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags CmdbSystem
// @Summary 更新编辑器关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.ApplicationSystemEditRelation true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/system/updateEditRelation [post]
func (a *CmdbSystemApi) UpdateEditRelation(c *gin.Context) {
	var relation application.ApplicationSystemEditRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(relation, utils.SystemEditRelationVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if relation.ID > 0 {
		if err := cmdbSystemService.UpdateEditRelations(relation); err != nil {
			global.GVA_LOG.Error("更新失败!", zap.Any("err", err))

			response.FailWithMessage("更新失败", c)
		} else {
			response.OkWithMessage("更新成功", c)
		}
	} else {
		if err := cmdbSystemService.AddEditRelations(relation); err != nil {
			global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

			response.FailWithMessage("添加失败", c)
		} else {
			response.OkWithMessage("添加成功", c)
		}
	}
}

// @Tags CmdbSystem
// @Summary 获取编辑器关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/system/getSystemEditRelation [post]
func (a *CmdbSystemApi) GetSystemEditRelations(c *gin.Context) {
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
	if err, relations := cmdbSystemService.GetSystemEditRelations(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(relations, "获取成功", c)
	}
}
