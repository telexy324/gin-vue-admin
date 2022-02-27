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

type StaffApi struct {
}

// @Tags Staff
// @Summary 新增管理员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.Admin true "姓名, 电话, 邮件, 部门"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /cmdb/addAdmin [post]
func (a *StaffApi) AddAdmin(c *gin.Context) {
	var admin application.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(admin, utils.ServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := staffService.AddAdmin(admin); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Staff
// @Summary 删除管理员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteAdmin [post]
func (a *StaffApi) DeleteAdmin(c *gin.Context) {
	var admin request.GetById
	if err := c.ShouldBindJSON(&admin); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(admin, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := staffService.DeleteAdmin(admin.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Staff
// @Summary 更新管理员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body application.Admin true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cmdb/updateAdmin [post]
func (a *StaffApi) UpdateAdmin(c *gin.Context) {
	var admin application.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(admin, utils.AdminVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := staffService.UpdateAdmin(admin); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Staff
// @Summary 根据id获取管理员
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "管理员id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAdminById [post]
func (a *StaffApi) GetAdminById(c *gin.Context) {
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
	if err, admin := staffService.GetAdminById(idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.AdminResponse{
			Admin: admin,
		}, "获取成功", c)
	}
}

// @Tags Staff
// @Summary 分页获取基础admin列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAdminList [post]
func (a *StaffApi) GetAdminList(c *gin.Context) {
	var pageInfo request2.AdminSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, adminList, total := staffService.GetAdminList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     adminList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags Staff
// @Summary 获取所有部门
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getDepartmentAll [post]
func (a *StaffApi) GetDepartmentAll(c *gin.Context) {
	if err, departmentList := staffService.GetDepartmentAll(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(applicationRes.DepartmentsResponse{
			Departments: departmentList,
		}, "获取成功", c)
	}
}
