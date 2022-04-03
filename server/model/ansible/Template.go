package ansible

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
)

type TemplateType string

const (
	TemplateTask   TemplateType = ""
	TemplateBuild  TemplateType = "build"
	TemplateDeploy TemplateType = "deploy"
)

type SurveyVarType string

const (
	SurveyVarStr TemplateType = ""
	SurveyVarInt TemplateType = "int"
)

type SurveyVar struct {
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Required    bool          `json:"required"`
	Type        SurveyVarType `json:"type"`
	Description string        `json:"description"`
}

type TemplateFilter struct {
	ViewID          *int
	BuildTemplateID *int
	AutorunOnly     bool
}

// Template is a user defined model that is used to run a task
type Template struct {
	ID int `db:"id" json:"id"`

	ProjectID     int  `db:"project_id" json:"project_id"`
	InventoryID   int  `db:"inventory_id" json:"inventory_id"`
	RepositoryID  int  `db:"repository_id" json:"repository_id"`
	EnvironmentID *int `db:"environment_id" json:"environment_id"`

	// Name as described in https://github.com/ansible-semaphore/semaphore/issues/188
	Name string `db:"name" json:"name"`
	// playbook name in the form of "some_play.yml"
	Playbook string `db:"playbook" json:"playbook"`
	// to fit into []string
	Arguments *string `db:"arguments" json:"arguments"`
	// if true, semaphore will not prepend any arguments to `arguments` like inventory, etc
	AllowOverrideArgsInTask bool `db:"allow_override_args_in_task" json:"allow_override_args_in_task"`

	Description *string `db:"description" json:"description"`

	VaultKeyID *int      `db:"vault_key_id" json:"vault_key_id"`
	VaultKey   AccessKey `db:"-" json:"-"`

	Type            TemplateType `db:"type" json:"type"`
	StartVersion    *string      `db:"start_version" json:"start_version"`
	BuildTemplateID *int         `db:"build_template_id" json:"build_template_id"`

	ViewID *int `db:"view_id" json:"view_id"`

	LastTask *TaskWithTpl `db:"-" json:"last_task"`

	Autorun bool `db:"autorun" json:"autorun"`

	// SurveyVarsJSON used internally for read from database.
	// It is not used for store survey vars to database.
	// Do not use it in your code. Use SurveyVars instead.
	SurveyVarsJSON *string     `db:"survey_vars" json:"-"`
	SurveyVars     []SurveyVar `db:"-" json:"survey_vars"`

	SuppressSuccessAlerts bool `db:"suppress_success_alerts" json:"suppress_success_alerts"`
}

func CreateTemplate(template Template) (Template, error) {
	err := global.GVA_DB.Create(&template).Error
	return template, err
}

func (m *Template) UpdateTemplate(template Template) error {
	var oldTemplate Template
	_, surveyVarsJson := json.Marshal(template.SurveyVars)
	upDateMap := make(map[string]interface{})
	upDateMap["inventory_id"] = template.InventoryID
	upDateMap["repository_id"] = template.RepositoryID
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

func (m *Template) GetTemplates(projectID int, filter TemplateFilter, sortInverted bool, sortBy string) (templates []Template, err error) {
	db := global.GVA_DB.Model(&Template{})
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
	if sortInverted {
		order = "desc"
	}
	switch sortBy {
	case "name", "playbook":
		db = db.Where("project_id=?", projectID).
			Order(sortBy + " " + order)
	case "inventory":
		db = db.Joins("left join project_inventory on inventory_id = project_inventory.id").
			Where("project_id=?", projectID).
			Order("project_inventory.name " + order)
	case "environment":
		db = db.Joins("project_environment on environment_id = environment.id)").
			Where("project_id=?", projectID).
			Order("project_environment.name " + order)
	case "repository":
		db = db.Joins("project_repository on repository_id = repository.id)").
			Where("project_id=?", projectID).
			Order("project_repository.name " + order)
	default:
		db = db.Where("project_id=?", projectID).
			Order("name " + order)
	}
	err = db.Find(&templates).Error
	if err != nil {
		return
	}
	err = FillTemplates(templates)
	return
}

func (m *Template) GetTemplate(projectID int, templateID int) (template Template, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, templateID).First(&template).Error
	if err != nil {
		return
	}
	err = FillTemplate(&template)
	return
}

func (m *Template) DeleteTemplate(projectID int, templateID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&Template{}).Error
	if err != nil {
		return err
	}
	var template Template
	return global.GVA_DB.Where("id = ? and project_id = ?", templateID, projectID).First(&template).Delete(&template).Error
}

//func (d *SqlDb) GetTemplateRefs(projectID int, templateID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.TemplateProps, templateID)
//}

func (tpl *Template) Validate() error {
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

func FillTemplates(templates []Template) (err error) {
	t := &Task{}
	for i := range templates {
		tpl := &templates[i]
		var tasks []TaskWithTpl
		err, iTasks, _ := t.GetTemplateTasks(tpl.ProjectID, tpl.ID, request.PageInfo{
			Page:     1,
			PageSize: 1,
		})
		tasks = iTasks.([]TaskWithTpl)
		if err != nil {
			return
		}
		if len(tasks) > 0 {
			tpl.LastTask = &tasks[0]
		}
	}
	return
}

func FillTemplate(template *Template) (err error) {
	k := &AccessKey{}
	if template.VaultKeyID != nil {
		template.VaultKey, err = k.GetAccessKey(template.ProjectID, *template.VaultKeyID)
	}

	if err != nil {
		return
	}

	err = FillTemplates([]Template{*template})

	if err != nil {
		return
	}

	if template.SurveyVarsJSON != nil {
		err = json.Unmarshal([]byte(*template.SurveyVarsJSON), &template.SurveyVars)
	}

	return
}
