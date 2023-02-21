package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
)

type TaskService struct {
}

var TaskServiceApp = new(TaskService)

func (taskService *TaskService) CreateTask(task ansible.Task) (ansible.Task, error) {
	err := global.GVA_DB.Create(&task).Error
	return task, err
}

func (taskService *TaskService) UpdateTask(task ansible.Task) error {
	var oldTask ansible.Task
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

func (taskService *TaskService) CreateTaskOutput(output ansible.TaskOutput) (ansible.TaskOutput, error) {
	err := global.GVA_DB.Create(&output).Error
	return output, err
}

func (taskService *TaskService) getTasks(projectID int, templateID *int, info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ansible.Task{}).Preload("User")
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
	var Tasks []ansible.Task
	var TaskWithTpls []ansible.TaskWithTpl
	err = db.Limit(limit).Offset(offset).Find(&Tasks).Error

	for _, task := range Tasks {
		taskWithTpl := ansible.TaskWithTpl{
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

func (taskService *TaskService) GetTask(projectID int, taskID int) (task ansible.Task, err error) {
	err = global.GVA_DB.Preload("Template", "project_id=?", projectID).
		Where("id = ?", taskID).First(&task).Error
	return
}

func (taskService *TaskService) GetTemplateTasks(projectID int, templateID int, info request.PageInfo) (err error, list interface{}, total int64) {
	return taskService.getTasks(projectID, &templateID, info)
}

func (taskService *TaskService) GetProjectTasks(projectID int, info request.PageInfo) (err error, list interface{}, total int64) {
	return taskService.getTasks(projectID, nil, info)
}

func (taskService *TaskService) DeleteTaskWithOutputs(projectID int, taskID int) (err error) {
	// check if task exists in the project
	_, err = taskService.GetTask(projectID, taskID)
	if err != nil {
		return
	}
	var task ansible.Task
	var taskOutputs []ansible.TaskOutput
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

func (taskService *TaskService) GetTaskOutputs(projectID int, taskID int, task *ansible.Task) (output []ansible.TaskOutput, err error) {
	// check if task exists in the project
	_, err = taskService.GetTask(projectID, taskID)
	if err != nil {
		return
	}

	err = global.GVA_DB.Where("task_id = ?", taskID).Order("time").Find(&output).Error
	return
}

func (taskService *TaskService) GetIncomingVersion(task ansible.Task) *string {
	if task.BuildTaskID == nil {
		return nil
	}

	buildTask, err := taskService.GetTask(task.ProjectID, *task.BuildTaskID)

	if err != nil {
		return nil
	}

	tpl, err := TemplatesServiceApp.GetTemplate(float64(task.ProjectID), float64(buildTask.TemplateID))
	if err != nil {
		return nil
	}

	if tpl.Type == ansible.TemplateBuild {
		return buildTask.Version
	}

	return taskService.GetIncomingVersion(buildTask)
}

func (taskService *TaskService) ValidateNewTask(template ansible.Template) error {
	switch template.Type {
	case ansible.TemplateBuild:
	case ansible.TemplateDeploy:
	case ansible.TemplateTask:
	}
	return nil
}

func (taskService *TaskService) Fill(task ansible.TaskWithTpl) error {
	if task.Task.BuildTaskID != nil {
		build, err := taskService.GetTask(task.Task.ProjectID, *task.Task.BuildTaskID)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		task.BuildTask = &build
	}
	return nil
}
