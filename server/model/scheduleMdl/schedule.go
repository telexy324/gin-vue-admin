package scheduleMdl

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Schedule struct {
	global.GVA_MODEL
	TemplateID int    `gorm:"type:bigint;not null;default:0;column:template_id" json:"templateId"`
	CronFormat string `gorm:"type:varchar(255);not null;default:'';column:cron_format" json:"cronFormat"`
	// LastCommitHash *string `gorm:"column:last_commit_hash" json:"-"`
	Valid       int      `gorm:"type:tinyint(2);not null;default:0;column:valid" json:"valid"`
	CommandVar  string   `json:"commandVar" gorm:"type:text;column:command_var"`
	TargetId    string   `json:"targetId" gorm:"type:text;column:target_id"`
	CommandVars []string `json:"commandVars" gorm:"-"`
	TargetIds   []int    `json:"targetIds" gorm:"-"` // 结束时间
}

func (m *Schedule) TableName() string {
	return "application_task_schedules"
}

func (m *Schedule) AfterFind(tx *gorm.DB) (err error) {
	commandVars := make([]string, 0)
	if m.CommandVar != "" {
		if err = json.Unmarshal([]byte(m.CommandVar), &commandVars); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
		m.CommandVars = commandVars
	}
	targetIds := make([]int, 0)
	if m.TargetId != "" {
		if err = json.Unmarshal([]byte(m.TargetId), &targetIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
		m.TargetIds = targetIds
	}
	return nil
}
