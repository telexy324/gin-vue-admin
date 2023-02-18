package application

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"gorm.io/gorm"
	"strings"
)

type TaskService struct {
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddTaskTemplate
//@description: 添加任务模板
//@param: taskTemplate model.TaskTemplate
//@return: error

func (taskService *TaskService) AddTaskTemplate(taskTemplate application.TaskTemplate) error {
	return global.GVA_DB.Create(&taskTemplate).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteTaskTemplate
//@description: 删除任务模板
//@param: id float64
//@return: err error

func (taskService *TaskService) DeleteTaskTemplate(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.TaskTemplate{}).Error
	if err != nil {
		return
	}
	var taskTemplate application.TaskTemplate
	return global.GVA_DB.Where("id = ?", id).First(&taskTemplate).Delete(&taskTemplate).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateTaskTemplate
//@description: 更新任务模板
//@param: server model.TaskTemplate
//@return: err error

func (taskService *TaskService) UpdateTaskTemplate(taskTemplate application.TaskTemplate) (err error) {
	upDateMap := make(map[string]interface{})
	if len(taskTemplate.TargetServers) > 0 {
		if js, err := json.Marshal(taskTemplate.TargetServers); err != nil {
			return err
		} else {
			upDateMap["task_server_ids"] = string(js)
		}
	}
	upDateMap["name"] = taskTemplate.Name
	upDateMap["description"] = taskTemplate.Description
	upDateMap["mode"] = taskTemplate.Mode
	upDateMap["command"] = taskTemplate.Command
	upDateMap["script_path"] = taskTemplate.ScriptPath
	upDateMap["cron"] = taskTemplate.Cron
	upDateMap["last_task_id"] = taskTemplate.LastTaskId
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetTaskTemplateById
//@description: 返回当前选中taskTemplate
//@param: id float64
//@return: err error, server model.TaskTemplate

func (taskService *TaskService) GetTaskTemplateById(id float64) (err error, taskTemplate application.TaskTemplate) {
	err = global.GVA_DB.Where("id = ?", id).First(&taskTemplate).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetTaskTemplateList
//@description: 获取任务模板分页
//@return: err error, list interface{}, total int64

func (taskService *TaskService) GetTaskTemplateList(info request2.TaskTemplateSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var taskTemplateList []application.TaskTemplate
	db := global.GVA_DB.Model(&application.TaskTemplate{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&taskTemplateList).Error
	return err, taskTemplateList, total
}
