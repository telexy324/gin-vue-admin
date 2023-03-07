package taskSvr

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"gorm.io/gorm"
)

type TaskTemplatesService struct {
}

var TaskTemplatesServiceApp = new(TaskTemplatesService)

func (templateService *TaskTemplatesService) CreateTaskTemplate(template taskMdl.TaskTemplate) (taskMdl.TaskTemplate, error) {
	targetServersJson, err := json.Marshal(template.TargetIds)
	if err != nil {
		return template, err
	}
	s := string(targetServersJson)
	template.TargetServerIds = s
	err = global.GVA_DB.Create(&template).Error
	return template, err
}

func (templateService *TaskTemplatesService) UpdateTaskTemplate(template taskMdl.TaskTemplate) error {
	var oldTaskTemplate taskMdl.TaskTemplate
	targetServersJson, _ := json.Marshal(template.TargetIds)
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = template.Name
	upDateMap["description"] = template.Description
	upDateMap["target_server_ids"] = targetServersJson
	upDateMap["mode"] = template.Mode
	upDateMap["command"] = template.Command
	upDateMap["script_path"] = template.ScriptPath
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

func (templateService *TaskTemplatesService) GetTaskTemplates(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var templates []taskMdl.TaskTemplate
	db := global.GVA_DB.Model(&taskMdl.TaskTemplate{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&templates).Error
	if err != nil {
		return
	}
	//err = templateService.FillTaskTemplates(templates)
	return err, templates, total
}

func (templateService *TaskTemplatesService) GetTaskTemplate(templateID float64) (template taskMdl.TaskTemplate, err error) {
	err = global.GVA_DB.Where("id =?", templateID).First(&template).Error
	if err != nil {
		return
	}
	//err = templateService.FillTaskTemplate(&template)
	return
}

func (templateService *TaskTemplatesService) DeleteTaskTemplate(templateID float64) error {
	err := global.GVA_DB.Where("id = ?", templateID).First(&taskMdl.TaskTemplate{}).Error
	if err != nil {
		return err
	}
	var template taskMdl.TaskTemplate
	return global.GVA_DB.Where("id = ?", templateID).First(&template).Delete(&template).Error
}

//func (d *SqlDb) GetTaskTemplateRefs(projectID int, templateID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.TaskTemplateProps, templateID)
//}

func (templateService *TaskTemplatesService) Validate(tpl taskMdl.TaskTemplate) error {
	if tpl.Name == "" {
		return errors.New("template name can not be empty")
	}

	if tpl.Command == "" && tpl.ScriptPath == "" {
		return errors.New("template playbook can not be empty")
	}

	return nil
}

//func (templateService *TaskTemplatesService) FillTaskTemplates(templates []task.TaskTemplate) (err error) {
//	for i := range templates {
//		tpl := &templates[i]
//		var tasks []task.TaskWithTpl
//		e, iTasks, _ := TaskServiceApp.GetTaskTemplateTasks(tpl.ProjectID, int(tpl.ID), request.PageInfo{
//			Page:     1,
//			PageSize: 1,
//		})
//		tasks = iTasks.([]task.TaskWithTpl)
//		if e != nil {
//			return e
//		}
//		if len(tasks) > 0 {
//			tpl.LastTask = &tasks[0]
//		}
//	}
//	return
//}
//
//func (templateService *TaskTemplatesService) FillTaskTemplate(template *task.TaskTemplate) (err error) {
//	if template.VaultKeyID != nil {
//		template.VaultKey, err = KeyServiceApp.GetAccessKey(float64(template.ProjectID), float64(*template.VaultKeyID))
//	}
//
//	if err != nil {
//		return
//	}
//
//	err = templateService.FillTaskTemplates([]task.TaskTemplate{*template})
//
//	if err != nil {
//		return
//	}
//
//	if template.SurveyVarsJSON != nil {
//		err = json.Unmarshal([]byte(*template.SurveyVarsJSON), &template.SurveyVars)
//	}
//
//	return
//}
