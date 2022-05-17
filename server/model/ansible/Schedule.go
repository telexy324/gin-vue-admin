package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Schedule struct {
	global.GVA_MODEL
	ProjectID      int     `gorm:"project_id" json:"project_id"`
	TemplateID     int     `gorm:"template_id" json:"template_id"`
	CronFormat     string  `gorm:"cron_format" json:"cron_format"`
	LastCommitHash *string `gorm:"last_commit_hash" json:"-"`
}
