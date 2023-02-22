package taskMdl

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type TaskTemplate struct {
	global.GVA_MODEL
	Name            string `json:"name" gorm:"column:name"`                         // task名称
	Description     string `json:"description" gorm:"column:description"`           // task描述
	TargetServerIds string `json:"targetServerIds" gorm:"column:target_server_ids"` // 关联服务器id
	Mode            int    `json:"mode" gorm:"column:mode"`                         // 执行方式 1 命令 2 脚本
	Command         string `json:"command" gorm:"column:command"`                   // 命令
	ScriptPath      string `json:"scriptPath" gorm:"column:script_path"`            // 脚本位置
	Cron            string `json:"cron" gorm:"column:cron"`                         // 定时任务
	LastTaskId      int    `json:"lastTaskId" gorm:"column:last_task_id"`           // 最后一次task id
	TargetServers   []int  `json:"targetServers" gorm:"-"`
}

func (m *TaskTemplate) TableName() string {
	return "application_task_templates"
}