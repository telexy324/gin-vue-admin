package taskMdl

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskTemplate struct {
	global.GVA_MODEL
	Name            string                          `json:"name" gorm:"column:name"`                         // task名称
	Description     string                          `json:"description" gorm:"column:description"`           // task描述
	TargetServerIds string                          `json:"targetServerIds" gorm:"column:target_server_ids"` // 关联服务器id
	Mode            int                             `json:"mode" gorm:"column:mode"`                         // 执行方式 1 命令 2 脚本
	Command         string                          `json:"command" gorm:"column:command"`                   // 命令
	ScriptPath      string                          `json:"scriptPath" gorm:"column:script_path"`            // 脚本位置
	LastTaskId      int                             `json:"lastTaskId" gorm:"column:last_task_id"`           // 最后一次task id
	SysUser         string                          `json:"sysUser" gorm:"column:sys_user"`                  // 执行用户
	TargetIds       []int                           `json:"targetIds" gorm:"-"`
	TargetServers   []application.ApplicationServer `json:"targetServers" gorm:"-"`
	LastTask        Task                            `json:"lastTask" gorm:"-"`
}

func (m *TaskTemplate) TableName() string {
	return "application_task_templates"
}

func (m *TaskTemplate) AfterFind(tx *gorm.DB) (err error) {
	serverIds := make([]int, 0)
	if m.TargetServerIds != "" {
		if err = json.Unmarshal([]byte(m.TargetServerIds), &serverIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(serverIds) > 0 {
		for _, id := range serverIds {
			server := application.ApplicationServer{}
			if err = tx.Model(&application.ApplicationServer{}).Where("id = ?", id).Find(&server).Error; err != nil {
				global.GVA_LOG.Error("转换失败", zap.Any("err", err))
				return
			}
			m.TargetServers = append(m.TargetServers, server)
		}
	}
	if m.LastTaskId > 0 {
		if err = tx.Model(&Task{}).Where("id = ?", serverIds).Find(&m.LastTask).Error; err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	return nil
}
