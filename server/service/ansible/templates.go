package ansible

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
)

type TemplatesService struct {
}

var TemplatesServiceApp = new(TemplatesService)

func (templateService *TemplatesService) CreateTemplate(template ansible.Template) (ansible.Template, error) {
	err := global.GVA_DB.Create(&template).Error
	return template, err
}

func (templateService *TemplatesService) UpdateTemplate(template ansible.Template) error {
	var oldTemplate ansible.Template
	_, surveyVarsJson := json.Marshal(template.SurveyVars)
	upDateMap := make(map[string]interface{})
	upDateMap["inventory_id"] = template.InventoryID
	upDateMap["environment_id"] = template.EnvironmentID
	upDateMap["name"] = template.Name
	upDateMap["playbook"] = template.Playbook
	upDateMap["arguments"] = template.Arguments
	upDateMap["allow_override_args_in_task"] = template.AllowOverrideArgsInTask
	upDateMap["description"] = template.Description
	upDateMap["vault_key_id"] = template.VaultKeyID
	upDateMap["`type`"] = template.Type
	upDateMap["start_version"] = template.StartVersion
	upDateMap["build_template_id"] = template.BuildTemplateID
	upDateMap["view_id"] = template.ViewID
	upDateMap["autorun"] = template.Autorun
	upDateMap["survey_vars"] = surveyVarsJson
	upDateMap["suppress_success_alerts"] = template.SuppressSuccessAlerts

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", template.ID, template.ProjectID).Find(&oldTemplate)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (templateService *TemplatesService) GetTemplates(info request2.GetByProjectId, filter ansible.TemplateFilter) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var templates []ansible.Template
	db := global.GVA_DB.Model(&ansible.Template{})
	if filter.ViewID != nil {
		db = db.Where("view_id=?", *filter.ViewID)
	}
	if filter.BuildTemplateID != nil {
		db = db.Where("build_template_id=?", *filter.BuildTemplateID)
		if filter.AutorunOnly {
			db = db.Where("autorun=true")
		}
	}
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
	err = templateService.FillTemplates(templates)
	return
}

func (templateService *TemplatesService) GetTemplate(projectID float64, templateID float64) (template ansible.Template, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, templateID).First(&template).Error
	if err != nil {
		return
	}
	err = templateService.FillTemplate(&template)
	return
}

func (templateService *TemplatesService) DeleteTemplate(projectID float64, templateID float64) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&ansible.Template{}).Error
	if err != nil {
		return err
	}
	var template ansible.Template
	return global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&template).Delete(&template).Error
}

//func (d *SqlDb) GetTemplateRefs(projectID int, templateID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.TemplateProps, templateID)
//}

func (templateService *TemplatesService) Validate(tpl ansible.Template) error {
	if tpl.Name == "" {
		return errors.New("template name can not be empty")
	}

	if tpl.Playbook == "" {
		return errors.New("template playbook can not be empty")
	}

	if tpl.Arguments != nil {
		if !json.Valid([]byte(*tpl.Arguments)) {
			return errors.New("template arguments must be valid JSON")
		}
	}

	return nil
}

func (templateService *TemplatesService) FillTemplates(templates []ansible.Template) (err error) {
	for i := range templates {
		tpl := &templates[i]
		var tasks []ansible.TaskWithTpl
		e, iTasks, _ := TaskServiceApp.GetTemplateTasks(tpl.ProjectID, int(tpl.ID), request.PageInfo{
			Page:     1,
			PageSize: 1,
		})
		tasks = iTasks.([]ansible.TaskWithTpl)
		if e != nil {
			return e
		}
		if len(tasks) > 0 {
			tpl.LastTask = &tasks[0]
		}
	}
	return
}

func (templateService *TemplatesService) FillTemplate(template *ansible.Template) (err error) {
	if template.VaultKeyID != nil {
		template.VaultKey, err = KeyServiceApp.GetAccessKey(float64(template.ProjectID), float64(*template.VaultKeyID))
	}

	if err != nil {
		return
	}

	err = templateService.FillTemplates([]ansible.Template{*template})

	if err != nil {
		return
	}

	if template.SurveyVarsJSON != nil {
		err = json.Unmarshal([]byte(*template.SurveyVarsJSON), &template.SurveyVars)
	}

	return
}
