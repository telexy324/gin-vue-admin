package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Schedule struct {
	global.GVA_MODEL
	ProjectID      int     `gorm:"column:project_id" json:"projectId"`
	TemplateID     int     `gorm:"column:template_id" json:"templateId"`
	CronFormat     string  `gorm:"column:cron_format" json:"cronFormat"`
	LastCommitHash *string `gorm:"column:last_commit_hash" json:"-"`
}

func (m *Schedule) TableName() string {
	return "ansible_schedules"
}
