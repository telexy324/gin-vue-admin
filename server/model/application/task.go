package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

type TaskStatus string

const (
	TaskRunningStatus  TaskStatus = "running"
	TaskWaitingStatus  TaskStatus = "waiting"
	TaskStoppingStatus TaskStatus = "stopping"
	TaskStoppedStatus  TaskStatus = "stopped"
	TaskSuccessStatus  TaskStatus = "success"
	TaskFailStatus     TaskStatus = "error"
)

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

type Task struct {
	global.GVA_MODEL
	TemplateId   int        `json:"templateId" gorm:"column:template_id"`      // task id
	Status       TaskStatus `json:"status" gorm:"column:status"`               // 状态
	SystemUserId int        `json:"systemUserId" gorm:"column:system_user_id"` // 执行人
	BeginTime    time.Time  `json:"beginTime" gorm:"column:begin_time"`        // 开始时间
	EndTime      time.Time  `json:"endTime" gorm:"column:end_time"`            // 结束时间
}

func (m *Task) TableName() string {
	return "application_tasks"
}

type TaskOutput struct {
	global.GVA_MODEL
	TaskId     int       `json:"taskId" gorm:"column:task_id"`         // task id
	RecordTime time.Time `json:"recordTime" gorm:"column:record_time"` // 记录时间
	Output     string    `json:"output" gorm:"column:output"`          // 输出
}

func (m *TaskOutput) TableName() string {
	return "application_task_outputs"
}