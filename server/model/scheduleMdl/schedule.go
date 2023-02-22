package scheduleMdl

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Schedule struct {
	global.GVA_MODEL
	TemplateID     int     `gorm:"column:template_id" json:"templateId"`
	CronFormat     string  `gorm:"column:cron_format" json:"cronFormat"`
	// LastCommitHash *string `gorm:"column:last_commit_hash" json:"-"`
}

func (m *Schedule) TableName() string {
	return "ansible_task_schedules"
}

