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

	ProjectID     int  `gorm:"project_id" json:"project_id"`
	InventoryID   int  `gorm:"inventory_id" json:"inventory_id"`
	EnvironmentID *int `gorm:"environment_id" json:"environment_id"`

	// Name as described in https://github.com/ansible-semaphore/semaphore/issues/188
	Name string `gorm:"name" json:"name"`
	// playbook name in the form of "some_play.yml"
	Playbook string `gorm:"playbook" json:"playbook"`
	// to fit into []string
	Arguments *string `gorm:"arguments" json:"arguments"`
	// if true, semaphore will not prepend any arguments to `arguments` like inventory, etc
	AllowOverrideArgsInTask bool `gorm:"allow_override_args_in_task" json:"allow_override_args_in_task"`

	Description *string `gorm:"description" json:"description"`

	VaultKeyID *int      `gorm:"vault_key_id" json:"vault_key_id"`
	VaultKey   AccessKey `gorm:"-" json:"-"`

	Type            TemplateType `gorm:"type" json:"type"`
	StartVersion    *string      `gorm:"start_version" json:"start_version"`
	BuildTemplateID *int         `gorm:"build_template_id" json:"build_template_id"`

	ViewID *int `gorm:"view_id" json:"view_id"`

	LastTask *TaskWithTpl `gorm:"-" json:"last_task"`

	Autorun bool `gorm:"autorun" json:"autorun"`

	// SurveyVarsJSON used internally for read from database.
	// It is not used for store survey vars to database.
	// Do not use it in your code. Use SurveyVars instead.
	SurveyVarsJSON *string     `gorm:"survey_vars" json:"-"`
	SurveyVars     []SurveyVar `gorm:"-" json:"survey_vars"`

	SuppressSuccessAlerts bool `gorm:"suppress_success_alerts" json:"suppress_success_alerts"`
}
