package ansible

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

	Arguments *string  `gorm:"column:arguments" json:"arguments"`
	Template  Template `json:"template"`
	Project   Project  `json:"project"`
}

// TaskWithTpl is the task data with additional fields
type TaskWithTpl struct {
	global.GVA_MODEL
	Task
	TemplatePlaybook string       `gorm:"tpl_playbook" json:"tpl_playbook"`
	TemplateAlias    string       `gorm:"tpl_alias" json:"tpl_alias"`
	TemplateType     TemplateType `gorm:"tpl_type" json:"tpl_type"`
	UserName         *string      `gorm:"user_name" json:"user_name"`
	BuildTask        *Task        `gorm:"-" json:"build_task"`
}

// TaskOutput is the ansible log output from the task
type TaskOutput struct {
	global.GVA_MODEL
	TaskID int       `gorm:"task_id" json:"task_id"`
	Task   string    `gorm:"task" json:"task"`
	Time   time.Time `gorm:"time" json:"time"`
	Output string    `gorm:"output" json:"output"`
}

func (m *Task) GetTask(projectID int, taskID int) (task Task, err error) {
	err = global.GVA_DB.Preload("Template","project_id=?",projectID).
		Where("id = ?", taskID).First(&task).Error
	return
}
