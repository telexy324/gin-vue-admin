package task

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/task"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/task/request"
	"gorm.io/gorm"
)

type TaskTemplatesService struct {
}

var TaskTemplatesServiceApp = new(TaskTemplatesService)

func (templateService *TaskTemplatesService) CreateTaskTemplate(template task.TaskTemplate) (task.TaskTemplate, error) {
	targetServersJson, err := json.Marshal(template.TargetServers)
	if err != nil {
		return template, err
	}
	s := string(targetServersJson)
	template.TargetServerIds = s
	err = global.GVA_DB.Create(&template).Error
	return template, err
}

func (templateService *TaskTemplatesService) UpdateTaskTemplate(template task.TaskTemplate) error {
	var oldTaskTemplate task.TaskTemplate
	_, targetServersJson := json.Marshal(template.TargetServers)
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = template.Name
	upDateMap["description"] = template.Description
	upDateMap["target_server_ids"] = targetServersJson
	upDateMap["mode"] = template.Mode
	upDateMap["command"] = template.Command
	upDateMap["script_path"] = template.ScriptPath
	upDateMap["cron"] = template.Cron
	upDateMap["last_task_id"] = template.LastTaskId

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", template.ID).Find(&oldTaskTemplate)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (templateService *TaskTemplatesService) GetTaskTemplates(info request2.GetByProjectId) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var templates []task.TaskTemplate
	db := global.GVA_DB.Model(&task.TaskTemplate{})
	order := ""
	if info.SortInverted {
		order = "desc"
	}
	switch info.SortBy {
	case "name", "playbook":
		db = db.Where("project_id=?", info.ProjectId).
			Order(info.SortBy + " " + order)
	case "inventory":
		db = db.Joins("left join project_inventory on inventory_id = project_inventory.id").
			Where("project_id=?", info.ProjectId).
			Order("project_inventory.name " + order)
	case "environment":
		db = db.Joins("project_environment on environment_id = environment.id)").
			Where("project_id=?", info.ProjectId).
			Order("project_environment.name " + order)
	case "repository":
		db = db.Joins("project_repository on repository_id = repository.id)").
			Where("project_id=?", info.ProjectId).
			Order("project_repository.name " + order)
	default:
		db = db.Where("project_id=?", info.ProjectId).
			Order("name " + order)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&templates).Error
	if err != nil {
		return
	}
	err = templateService.FillTaskTemplates(templates)
	return err, templates, total
}

func (templateService *TaskTemplatesService) GetTaskTemplate(projectID float64, templateID float64) (template task.TaskTemplate, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, templateID).First(&template).Error
	if err != nil {
		return
	}
	err = templateService.FillTaskTemplate(&template)
	return
}

func (templateService *TaskTemplatesService) DeleteTaskTemplate(projectID float64, templateID float64) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&task.TaskTemplate{}).Error
	if err != nil {
		return err
	}
	var template task.TaskTemplate
	return global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&template).Delete(&template).Error
}

//func (d *SqlDb) GetTaskTemplateRefs(projectID int, templateID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.TaskTemplateProps, templateID)
//}

func (templateService *TaskTemplatesService) Validate(tpl task.TaskTemplate) error {
	if tpl.Name == "" {
		return errors.New("template name can not be empty")
	}

	if tpl.Command == "" && tpl.ScriptPath == "" {
		return errors.New("template playbook can not be empty")
	}

	return nil
}

func (templateService *TaskTemplatesService) FillTaskTemplates(templates []task.TaskTemplate) (err error) {
	for i := range templates {
		tpl := &templates[i]
		var tasks []task.TaskWithTpl
		e, iTasks, _ := TaskServiceApp.GetTaskTemplateTasks(tpl.ProjectID, int(tpl.ID), request.PageInfo{
			Page:     1,
			PageSize: 1,
		})
		tasks = iTasks.([]task.TaskWithTpl)
		if e != nil {
			return e
		}
		if len(tasks) > 0 {
			tpl.LastTask = &tasks[0]
		}
	}
	return
}

func (templateService *TaskTemplatesService) FillTaskTemplate(template *task.TaskTemplate) (err error) {
	if template.VaultKeyID != nil {
		template.VaultKey, err = KeyServiceApp.GetAccessKey(float64(template.ProjectID), float64(*template.VaultKeyID))
	}

	if err != nil {
		return
	}

	err = templateService.FillTaskTemplates([]task.TaskTemplate{*template})

	if err != nil {
		return
	}

	if template.SurveyVarsJSON != nil {
		err = json.Unmarshal([]byte(*template.SurveyVarsJSON), &template.SurveyVars)
	}

	return
}
