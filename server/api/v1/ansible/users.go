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

type UsersApi struct {
}

// @Tags User
// @Summary 新增User
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.User true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/user/addUser [post]
func (a *UsersApi) AddUser(c *gin.Context) {
	var userRequest request2.AddUserByProjectId
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(userRequest, utils.ProjectIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := userService.CreateProjectUser(ansible.ProjectUser{ProjectId: int(userRequest.ProjectId), UserId: int(userRequest.UserId), Admin: userRequest.Admin}); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags User
// @Summary 删除User
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "EnvronmentId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/user/deleteUser [post]
func (a *UsersApi) DeleteUser(c *gin.Context) {
	var userRequest request2.DeleteUserByProjectId
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(userRequest, utils.ProjectIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := userService.DeleteProjectUser(int(userRequest.ProjectId), int(userRequest.UserId)); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags User
// @Summary 更新User
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.User true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/user/updateUser [post]
func (a *UsersApi) UpdateUser(c *gin.Context) {
	var userRequest request2.AddUserByProjectId
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(userRequest, utils.ProjectIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := userService.UpdateProjectUser(ansible.ProjectUser{ProjectId: int(userRequest.ProjectId), UserId: int(userRequest.UserId), Admin: userRequest.Admin}); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags User
// @Summary 分页获取基础User列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/user/getProjectUsers[post]
func (a *UsersApi) GetProjectUsers(c *gin.Context) {
	var userRequest request2.GetByProjectId
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(userRequest.ProjectId, utils.ProjectIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if users, err := userService.GetProjectUsers(int(userRequest.ProjectId)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.UsersResponse{
			Users: users,
		}, "获取成功", c)
	}
}
