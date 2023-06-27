package scheduleMdl

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Schedule struct {
	global.GVA_MODEL
	TemplateID int    `gorm:"type:bigint;not null;default:0;column:template_id" json:"templateId"`
	CronFormat string `gorm:"type:varchar(255);not null;default:'';column:cron_format" json:"cronFormat"`
	// LastCommitHash *string `gorm:"column:last_commit_hash" json:"-"`
	Valid int `gorm:"type:tinyint(2);not null;default:0;column:valid" json:"valid"`
}

func (m *Schedule) TableName() string {
	return "application_task_schedules"
}
