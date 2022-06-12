package middleware

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ansibleUserService = service.ServiceGroupApp.AnsibleServiceGroup.UserService

func MustBeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project ansible.Project
		if err := c.ShouldBindJSON(&project); err != nil {
			global.GVA_LOG.Info("error", zap.Any("err", err))
			response.FailWithMessage(err.Error(), c)
			return
		}
		if err := utils.Verify(project, utils.ProjectVerify); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		userID := int(utils.GetUserID(c))
		user, err := ansibleUserService.GetProjectUser(float64(project.ID), float64(userID))
		if err != nil {
			global.GVA_LOG.Error("获取project管理员失败!", zap.Any("err", err))
			response.FailWithMessage("验证失败", c)
		} else if user.Admin != ansible.IsAdmin {
			response.FailWithMessage("非管理员", c)
		}
	}
}
