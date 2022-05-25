package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
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

//Task is a model of a task which will be executed by the runner
type Task struct {
	global.GVA_MODEL
	TemplateID int `gorm:"column:template_id" json:"template_id"`
	ProjectID  int `gorm:"column:project_id" json:"project_id"`

	Status TaskStatus `gorm:"column:status" json:"status"`
	Debug  bool       `gorm:"column:debug" json:"debug"`

	DryRun bool `gorm:"column:dry_run" json:"dry_run"`

	// override variables
	Playbook    string `gorm:"column:playbook" json:"playbook"`
	Environment string `gorm:"column:environment" json:"environment"`

	UserID *int `gorm:"column:user_id" json:"user_id"`

	Created time.Time  `gorm:"column:created" json:"created"`
	Start   *time.Time `gorm:"column:start" json:"start"`
	End     *time.Time `gorm:"column:end" json:"end"`

	Message string `gorm:"column:message" json:"message"`

	// CommitMessage is a git commit hash of playbook repository which
	// was active when task was created.
	CommitHash *string `gorm:"column:commit_hash" json:"commit_hash"`
	// CommitMessage contains message retrieved from git repository after checkout to CommitHash.
	// It is readonly by API.
	CommitMessage string `gorm:"column:commit_message" json:"commit_message"`

	BuildTaskID *int `gorm:"column:build_task_id" json:"build_task_id"`

	// Version is a build version.
	// This field available only for Build tasks.
	Version *string `gorm:"column:version" json:"version"`

	Arguments *string        `gorm:"column:arguments" json:"arguments"`
	Template  Template       `json:"template"`
	Project   Project        `json:"project"`
	User      system.SysUser `json:"user"`
}

// TaskWithTpl is the task data with additional fields
type TaskWithTpl struct {
	Task             Task         `gorm:"-" json:"task"`
	TemplatePlaybook string       `gorm:"-" json:"tpl_playbook"`
	TemplateAlias    string       `gorm:"-" json:"tpl_alias"`
	TemplateType     TemplateType `gorm:"-" json:"tpl_type"`
	UserName         *string      `gorm:"-" json:"user_name"`
	BuildTask        *Task        `gorm:"-" json:"build_task"`
}

// TaskOutput is the ansible log output from the task
type TaskOutput struct {
	global.GVA_MODEL
	TaskID int       `gorm:"column:task_id" json:"task_id"`
	Task   string    `gorm:"column:task" json:"task"`
	Time   time.Time `gorm:"column:time" json:"time"`
	Output string    `gorm:"column:output" json:"output"`
}

func (m *Task) TableName() string {
	return "ansible_tasks"
}

func (m *TaskOutput) TableName() string {
	return "ansible_task_outputs"
}
