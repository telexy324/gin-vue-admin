package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
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
	Task
	TemplatePlaybook string       `gorm:"-" json:"tpl_playbook"`
	TemplateAlias    string       `gorm:"-" json:"tpl_alias"`
	TemplateType     TemplateType `gorm:"-" json:"tpl_type"`
	UserName         *string      `gorm:"-" json:"user_name"`
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

func CreateTask(task Task) (Task, error) {
	err := global.GVA_DB.Create(&task).Error
	return task, err
}

func (m *Task) UpdateTask(task Task) error {
	var oldTask Task
	upDateMap := make(map[string]interface{})
	upDateMap["status"] = task.Status
	upDateMap["start"] = task.Start
	upDateMap["end"] = task.End

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", task.ID).Find(&oldTask)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (m *TaskOutput) CreateTaskOutput(output TaskOutput) (TaskOutput, error) {
	err := global.GVA_DB.Create(&output).Error
	return output, err
}

func (m *Task) getTasks(projectID int, templateID *int, info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Preload("User")
	if templateID == nil {
		db = db.Preload("Template", "project_id=?", projectID)
	} else {
		db = db.Preload("Template", "project_id=?", projectID).Where("template_id=?", templateID)
	}
	db.Order("created desc, id desc")
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var Tasks []Task
	var TaskWithTpls []TaskWithTpl
	err = db.Limit(limit).Offset(offset).Find(&Tasks).Error

	for _, task := range Tasks {
		taskWithTpl := TaskWithTpl{
			Task:             task,
			TemplatePlaybook: task.Template.Playbook,
			TemplateAlias:    task.Template.Name,
			TemplateType:     task.Template.Type,
			UserName:         &task.User.Username,
		}
		TaskWithTpls = append(TaskWithTpls, taskWithTpl)
	}
	return err, TaskWithTpls, total
}

func (m *Task) GetTask(projectID int, taskID int) (task Task, err error) {
	err = global.GVA_DB.Preload("Template", "project_id=?", projectID).
		Where("id = ?", taskID).First(&task).Error
	return
}

func (m *Task) GetTemplateTasks(projectID int, templateID int, info request.PageInfo) (err error, list interface{}, total int64) {
	return m.getTasks(projectID, &templateID, info)
}

func (m *Task) GetProjectTasks(projectID int, info request.PageInfo) (err error, list interface{}, total int64) {
	return m.getTasks(projectID, nil, info)
}

func (m *Task) DeleteTaskWithOutputs(projectID int, taskID int) (err error) {
	// check if task exists in the project
	_, err = m.GetTask(projectID, taskID)
	if err != nil {
		return
	}
	var task Task
	var taskOutputs []TaskOutput
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Where("task_id = ?", taskID).Find(&taskOutputs).Delete(&taskOutputs).Error
		if txErr != nil {
			return txErr
		}
		txErr = tx.Where("id = ?", taskID).Find(&task).Delete(&task).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return
}

func (m *TaskOutput) GetTaskOutputs(projectID int, taskID int) (output []TaskOutput, err error) {
	var task *Task
	// check if task exists in the project
	_, err = task.GetTask(projectID, taskID)
	if err != nil {
		return
	}

	err = global.GVA_DB.Where("task_id = ?", taskID).Order("time").Find(&output).Error
	return
}
