package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/task"
	"gorm.io/gorm"
)

type TaskService struct {
}

var TaskServiceApp = new(TaskService)

func (taskService *TaskService) CreateTask(task task.Task) (task.Task, error) {
	err := global.GVA_DB.Create(&task).Error
	return task, err
}

func (taskService *TaskService) UpdateTask(t task.Task) error {
	var oldTask task.Task
	upDateMap := make(map[string]interface{})
	upDateMap["status"] = t.Status
	upDateMap["begin_time"] = t.BeginTime
	upDateMap["end_time"] = t.EndTime

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", t.ID).Find(&oldTask)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (taskService *TaskService) CreateTaskOutput(output task.TaskOutput) (task.TaskOutput, error) {
	err := global.GVA_DB.Create(&output).Error
	return output, err
}

func (taskService *TaskService) getTasks(projectID int, templateID *int, info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&task.Task{}).Preload("User")
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
	var Tasks []task.Task
	//var TaskWithTpls []task.TaskWithTpl
	err = db.Limit(limit).Offset(offset).Find(&Tasks).Error

	//for _, t := range Tasks {
	//	taskWithTpl := task.TaskWithTpl{
	//		Task:             t,
	//		TemplatePlaybook: t.Template.Playbook,
	//		TemplateAlias:    t.Template.Name,
	//		TemplateType:     t.Template.Type,
	//		UserName:         &t.User.Username,
	//	}
	//	TaskWithTpls = append(TaskWithTpls, taskWithTpl)
	//}
	return err, Tasks, total
}

func (taskService *TaskService) GetTask(projectID int, taskID int) (task task.Task, err error) {
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
	var t task.Task
	var taskOutputs []task.TaskOutput
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Where("task_id = ?", taskID).Find(&taskOutputs).Delete(&taskOutputs).Error
		if txErr != nil {
			return txErr
		}
		txErr = tx.Where("id = ?", taskID).Find(&t).Delete(&t).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return
}

func (taskService *TaskService) GetTaskOutputs(projectID int, taskID int) (output []task.TaskOutput, err error) {
	// check if task exists in the project
	_, err = taskService.GetTask(projectID, taskID)
	if err != nil {
		return
	}

	err = global.GVA_DB.Where("task_id = ?", taskID).Order("time").Find(&output).Error
	return
}

//func (taskService *TaskService) GetIncomingVersion(t task.Task) *string {
//	if t.BuildTaskID == nil {
//		return nil
//	}
//
//	buildTask, err := taskService.GetTask(t.ProjectID, *t.BuildTaskID)
//
//	if err != nil {
//		return nil
//	}
//
//	tpl, err := TemplatesServiceApp.GetTemplate(float64(t.ProjectID), float64(buildTask.TemplateID))
//	if err != nil {
//		return nil
//	}
//
//	if tpl.Type == t.TemplateBuild {
//		return buildTask.Version
//	}
//
//	return taskService.GetIncomingVersion(buildTask)
//}

//func (taskService *TaskService) ValidateNewTask(template task.Template) error {
//	switch template.Type {
//	case task.TemplateBuild:
//	case task.TemplateDeploy:
//	case task.TemplateTask:
//	}
//	return nil
//}

//func (taskService *TaskService) Fill(t task.TaskWithTpl) error {
//	if t.Task.BuildTaskID != nil {
//		build, err := taskService.GetTask(t.Task.ProjectID, *t.Task.BuildTaskID)
//		if err == gorm.ErrRecordNotFound {
//			return nil
//		}
//		if err != nil {
//			return err
//		}
//		t.BuildTask = &build
//	}
//	return nil
//}
