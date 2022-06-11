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
	TemplateID int `gorm:"column:template_id" json:"templateId"`
	ProjectID  int `gorm:"column:project_id" json:"projectId"`

	Status TaskStatus `gorm:"column:status" json:"status"`
	Debug  bool       `gorm:"column:debug" json:"debug"`

	DryRun bool `gorm:"column:dry_run" json:"dryRun"`

	// override variables
	Playbook    string `gorm:"column:playbook" json:"playbook"`
	Environment string `gorm:"column:environment" json:"environment"`

	UserID *int `gorm:"column:user_id" json:"userId"`

	Created time.Time  `gorm:"column:created" json:"created"`
	Start   *time.Time `gorm:"column:start" json:"start"`
	End     *time.Time `gorm:"column:end" json:"end"`

	Message string `gorm:"column:message" json:"message"`

	// CommitMessage is a git commit hash of playbook repository which
	// was active when task was created.
	CommitHash *string `gorm:"column:commit_hash" json:"commitHash"`
	// CommitMessage contains message retrieved from git repository after checkout to CommitHash.
	// It is readonly by API.
	CommitMessage string `gorm:"column:commit_message" json:"commitMessage"`

	BuildTaskID *int `gorm:"column:build_task_id" json:"buildTaskId"`

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
	TemplatePlaybook string       `gorm:"-" json:"tplPlaybook"`
	TemplateAlias    string       `gorm:"-" json:"tplAlias"`
	TemplateType     TemplateType `gorm:"-" json:"tplType"`
	UserName         *string      `gorm:"-" json:"userName"`
	BuildTask        *Task        `gorm:"-" json:"buildTask"`
}

// TaskOutput is the ansible log output from the task
type TaskOutput struct {
	global.GVA_MODEL
	TaskID int       `gorm:"column:task_id" json:"taskId"`
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
