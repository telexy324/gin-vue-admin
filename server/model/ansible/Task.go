package ansible

import (
	"database/sql"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"strings"
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

func (m *Task) GetTask(projectID int, taskID int) (task Task, err error) {
	err = global.GVA_DB.Preload("Template", "project_id=?", projectID).
		Where("id = ?", taskID).First(&task).Error
	return
}

func (m *Task) CreateTask(task Task) (Task, error) {
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

func (m *TaskOutput) getTasks(projectID int, templateID *int, info request.PageInfo, tasks *[]db.TaskWithTpl) (err error) {
	fields := "task.*"
	fields += ", tpl.playbook as tpl_playbook" +
		", `user`.name as user_name" +
		", tpl.name as tpl_alias" +
		", tpl.type as tpl_type"

	q := squirrel.Select(fields).
		From("task").
		Join("project__template as tpl on task.template_id=tpl.id").
		LeftJoin("`user` on task.user_id=`user`.id").
		OrderBy("task.created desc, id desc")

	if templateID == nil {
		q = q.Where("tpl.project_id=?", projectID)
	} else {
		q = q.Where("tpl.project_id=? AND task.template_id=?", projectID, templateID)
	}

	if params.Count > 0 {
		q = q.Limit(uint64(params.Count))
	}

	query, args, _ := q.ToSql()

	_, err = d.selectAll(tasks, query, args...)

	for i := range *tasks {
		err = (*tasks)[i].Fill(d)
		if err != nil {
			return
		}
	}

	return
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var tasks []Task
	err = global.GVA_DB.Preload("Template", "project_id=?", projectID).
		Where("id = ?", taskID).First(&task).Error
	db := global.GVA_DB.Model(&application.ApplicationSystem{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&systemList).Error

	systemInfoList := make([]applicationRes.ApplicationSystemResponse, 0, len(systemList))
	for _, system := range systemList {
		sysAdmins := make([]application.ApplicationSystemAdmin, 0)
		adminInfos := make([]application.Admin, 0)
		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&sysAdmins).Error; err != nil {
			return
		}
		for _, sysAdmin := range sysAdmins {
			var adminInfo = application.Admin{}
			if err = global.GVA_DB.Where("id = ?", sysAdmin.AdminId).First(&adminInfo).Error; err != nil {
				return
			} else {
				adminInfos = append(adminInfos, adminInfo)
			}
		}
		systemInfoList = append(systemInfoList, applicationRes.ApplicationSystemResponse{
			System: system,
			Admins: adminInfos,
		})
	}
	return err, systemInfoList, total
}

func (d *SqlDb) GetTask(projectID int, taskID int) (task db.Task, err error) {
	q := squirrel.Select("task.*").
		From("task").
		Join("project__template as tpl on task.template_id=tpl.id").
		Where("tpl.project_id=? AND task.id=?", projectID, taskID)

	query, args, err := q.ToSql()

	if err != nil {
		return
	}

	err = d.selectOne(&task, query, args...)

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
		return
	}

	if err != nil {
		return
	}

	return
}

func (d *SqlDb) GetTemplateTasks(projectID int, templateID int, params db.RetrieveQueryParams) (tasks []db.TaskWithTpl, err error) {
	err = d.getTasks(projectID, &templateID, params, &tasks)
	return
}

func (d *SqlDb) GetProjectTasks(projectID int, params db.RetrieveQueryParams) (tasks []db.TaskWithTpl, err error) {
	err = d.getTasks(projectID, nil, params, &tasks)
	return
}

func (d *SqlDb) DeleteTaskWithOutputs(projectID int, taskID int) (err error) {
	// check if task exists in the project
	_, err = d.GetTask(projectID, taskID)

	if err != nil {
		return
	}

	_, err = d.exec("delete from task__output where task_id=?", taskID)

	if err != nil {
		return
	}

	_, err = d.exec("delete from task where id=?", taskID)
	return
}

func (d *SqlDb) GetTaskOutputs(projectID int, taskID int) (output []db.TaskOutput, err error) {
	// check if task exists in the project
	_, err = d.GetTask(projectID, taskID)

	if err != nil {
		return
	}

	_, err = d.selectAll(&output,
		"select task_id, task, time, output from task__output where task_id=? order by time asc",
		taskID)
	return
}
