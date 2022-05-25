package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

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
	global.GVA_MODEL

	ProjectID     int  `gorm:"column:project_id" json:"project_id"`
	InventoryID   int  `gorm:"column:inventory_id" json:"inventory_id"`
	EnvironmentID *int `gorm:"column:environment_id" json:"environment_id"`

	// Name as described in https://github.com/ansible-semaphore/semaphore/issues/188
	Name string `gorm:"column:name" json:"name"`
	// playbook name in the form of "some_play.yml"
	Playbook string `gorm:"column:playbook" json:"playbook"`
	// to fit into []string
	Arguments *string `gorm:"column:arguments" json:"arguments"`
	// if true, semaphore will not prepend any arguments to `arguments` like inventory, etc
	AllowOverrideArgsInTask bool `gorm:"column:allow_override_args_in_task" json:"allow_override_args_in_task"`

	Description *string `gorm:"column:description" json:"description"`

	VaultKeyID *int      `gorm:"column:vault_key_id" json:"vault_key_id"`
	VaultKey   AccessKey `gorm:"-" json:"-"`

	Type            TemplateType `gorm:"column:type" json:"type"`
	StartVersion    *string      `gorm:"column:start_version" json:"start_version"`
	BuildTemplateID *int         `gorm:"column:build_template_id" json:"build_template_id"`

	ViewID *int `gorm:"column:view_id" json:"view_id"`

	LastTask *TaskWithTpl `gorm:"-" json:"last_task"`

	Autorun bool `gorm:"column:autorun" json:"autorun"`

	// SurveyVarsJSON used internally for read from database.
	// It is not used for store survey vars to database.
	// Do not use it in your code. Use SurveyVars instead.
	SurveyVarsJSON *string     `gorm:"column:survey_vars" json:"-"`
	SurveyVars     []SurveyVar `gorm:"-" json:"survey_vars"`

	SuppressSuccessAlerts bool `gorm:"column:suppress_success_alerts" json:"suppress_success_alerts"`
}

func (m *Template) TableName() string {
	return "ansible_templates"
}
