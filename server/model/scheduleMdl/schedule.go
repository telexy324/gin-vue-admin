package scheduleMdl

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Schedule struct {
	global.GVA_MODEL
	TemplateID int    `gorm:"column:template_id" json:"templateId"`
	CronFormat string `gorm:"column:cron_format" json:"cronFormat"`
	// LastCommitHash *string `gorm:"column:last_commit_hash" json:"-"`
	Valid int `gorm:"column:valid" json:"valid"`
}

func (m *Schedule) TableName() string {
	return "application_task_schedules"
}
