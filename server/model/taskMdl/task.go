package taskMdl

import (
	"database/sql"
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

type Task struct {
	global.GVA_MODEL
	TemplateId   int          `json:"templateId" gorm:"type:bigint;not null;default:0;column:template_id"`       // task id
	Status       TaskStatus   `json:"status" gorm:"type:varchar(255);not null;default:'';column:status"`         // 状态
	SystemUserId int          `json:"systemUserId" gorm:"type:bigint;not null;default:0;column:system_user_id" ` // 执行人
	BeginTime    sql.NullTime `json:"beginTime" gorm:"column:begin_time" swaggertype:"string"`                   // 开始时间
	EndTime      sql.NullTime `json:"endTime" gorm:"column:end_time" swaggertype:"string"`                       // 结束时间
	SetTaskId    int          `json:"setTaskId" gorm:"type:bigint;not null;default:0;column:set_task_id"`        // 结束时间
	//FileDownload string       `json:"fileDownload" gorm:"column:file_download"`                // 结束时间
	//SystemId     int          `json:"systemId" gorm:"column:system_id"`                        // 结束时间
	//SshUser      string       `json:"sshUser" gorm:"column:ssh_user"`                          // 结束时间
	FileDownload string `json:"fileDownload" gorm:"-"` // 结束时间
	SystemId     int    `json:"systemId" gorm:"-"`     // 结束时间
	SshUser      string `json:"sshUser" gorm:"-"`      // 结束时间
}

func (m *Task) TableName() string {
	return "application_tasks"
}

type TaskOutput struct {
	global.GVA_MODEL
	TaskId     int       `json:"taskId" gorm:"type:bigint;not null;default:0;column:task_id"`           // task id
	RecordTime time.Time `json:"recordTime" gorm:"column:record_time"`                                  // 记录时间
	Output     string    `json:"output" gorm:"type:text;column:output"`                                 // 输出
	ManageIp   string    `json:"manageIp" gorm:"type:varchar(15);not null;default:'';column:manage_ip"` // 输出
}

func (m *TaskOutput) TableName() string {
	return "application_task_outputs"
}

// TaskWithTpl is the task data with additional fields
//type TaskWithTpl struct {
//	Task          Task    `gorm:"-" json:"task"`
//	TemplateAlias string  `gorm:"-" json:"tplAlias"`
//	UserName      *string `gorm:"-" json:"userName"`
//	BuildTask     *Task   `gorm:"-" json:"buildTask"`
//}
