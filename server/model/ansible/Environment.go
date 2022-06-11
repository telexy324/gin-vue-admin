package ansible

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Environment is used to pass additional arguments, in json form to ansible
type Environment struct {
	global.GVA_MODEL
	Name      string  `gorm:"column:name" json:"name" binding:"required"`
	ProjectID int     `gorm:"column:project_id" json:"projectId"`
	Password  *string `gorm:"column:password" json:"password"`
	JSON      string  `gorm:"column:json" json:"json" binding:"required"`
}

func (m *Environment) TableName() string {
	return "ansible_environments"
}

func (env *Environment) Validate() error {
	if env.Name == "" {
		return errors.New("environment name can not be empty")
	}

	if !json.Valid([]byte(env.JSON)) {
		return errors.New("environment must be valid JSON")
	}

	return nil
}
