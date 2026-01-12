package jumpServer

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JumpServerApi struct {
}

// @Tags JumpServer
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /jumpServer/getToken [post]
func (a *JumpServerApi) GetToken(c *gin.Context) {
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
	userID := int(utils.GetUserID(c))
	if token, err := jumpServerService.GenerateJumpToken(userID, int(idInfo.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		c.JSON(200, gin.H{
			"url": fmt.Sprintf("myjump://%s", token),
		})
	}
}

// @Tags JumpServer
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "服务器id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /jumpServer/getServer [post]
func (a *JumpServerApi) GetServer(c *gin.Context) {
	var token request2.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(token, utils.TokenVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, err := jumpServerService.ParseJumpToken(token.Token)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	err, server := applicationServerService.GetServerById(float64(claims.TargetID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	////从数据库读取真实目标服务器信息
	//srv := GetServerByID(claims.TargetID)

	//// 生成一次性 Jumpserver 登录用户名
	//jumpUser := "token-" + req.Token[:12] // 或专门签名过的临时用户标识
	//
	//// 记录审计日志 / 创建会话
	//sessionID := CreateAuditSession(claims.UserID, srv.IP)
	//
	//c.JSON(200, gin.H{
	//	"target_ip":   srv.IP,
	//	"target_port": srv.Port,
	//	"user":        srv.User,
	//	"jump_user":   jumpUser,
	//	"session_id":  sessionID,
	//})
	response.OkWithDetailed(jumpServerMdl.ConnInfo{
		JumpHost: server.ManageIp,
		Port:     server.SshPort,
		User:     server.SshUser,
		//Protocol: "",
		//Client:   "",
		//Password: "",
	}, "获取成功", c)
}
