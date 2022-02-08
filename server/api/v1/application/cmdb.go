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
	global.GVA_LOG.Info("error",zap.Any("err", e))
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
	if err, server := cmdbService.SystemRelations(idInfo.ID); err != nil {
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

// 企业关系图
message GetEnterpriseRelationChartReq {
int64 company_id = 1;
string cid = 2;
}

message GetEnterpriseRelationChartRsp {
x.common.def.RelationPath path = 1;
bool convert_h5 = 999;
}

message RelationPath {
repeated Node nodes = 1;
repeated Link links = 2;
}

message Node {
// company_id、people_id
int64 id = 1;
x.common.graph.consts.GraphNodeType type = 2;
// company_name 、 people_name
string name = 3;
// value: number of relational nodes
int64 value = 4;
// id to_string
string s_id = 5;
}

message Link {
//GraphVectorType
x.common.graph.consts.GraphVectorType vector_type = 1;
//type string value
string vector_str_value = 2;
// property value
Property property = 3;
// start node id
int64 start_node_id = 4;
// end node id
int64 end_node_id = 5;
// string(start node id)
string s_start_node_id = 6;
// string(end node id)
string s_end_node_id = 7;
}

message Property {
// 职位
string position = 1;
// 持股比例
double shareholding_ratio = 2;
// 持股类型名称（investment、shareholder）
string shareholding_name = 3;
// 股权更新时间
string shareholding_update_date = 4;
}

enum GraphNodeType {
gn_unknown = 0;
gn_company = 1;
gn_people = 2;
}

enum GraphVectorType {
gv_unknown      = 0;
gv_shareholder  = 1;
gv_employee     = 2;
gv_subsidiary   = 3;
}
