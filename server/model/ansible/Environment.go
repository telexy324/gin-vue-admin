package ansible

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

// Environment is used to pass additional arguments, in json form to ansible
type Environment struct {
	ID        int     `db:"id" json:"id"`
	Name      string  `db:"name" json:"name" binding:"required"`
	ProjectID int     `db:"project_id" json:"project_id"`
	Password  *string `db:"password" json:"password"`
	JSON      string  `db:"json" json:"json" binding:"required"`
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

func (m *Environment) GetEnvironment(projectID int, environmentID int) (Environment, error) {
	var environment Environment
	err := global.GVA_DB.Where("project_id=? and id =?", projectID, environmentID).First(&environment).Error
	return environment, err
}

//func (d *SqlDb) GetEnvironmentRefs(projectID int, environmentID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.EnvironmentProps, environmentID)
//}

func (m *Environment) GetEnvironments(projectID int, sortInverted bool, sortBy string) ([]Environment, error) {
	var environments []Environment
	db := global.GVA_DB.Model(&Environment{})
	order := ""
	if sortInverted {
		order = "desc"
	}
	db = db.Where("project_id=?", projectID).Order(sortBy + " " + order)
	err := db.Find(&environments).Error
	return environments, err
}

func (m *Environment) UpdateEnvironment(env Environment) error {
	err := env.Validate()
	if err != nil {
		return err
	}
	var oldEnv Environment
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = env.Name
	upDateMap["json"] = env.JSON

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", env.ID, env.ProjectID).Find(&oldEnv)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func CreateEnvironment(env Environment) (newEnv Environment, err error) {
	err = env.Validate()
	if err != nil {
		return
	}
	err = global.GVA_DB.Create(&env).Error
	return env, err
}

func (m *Environment) DeleteEnvironment(projectID int, environmentID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", environmentID, projectID).First(&Environment{}).Error
	if err != nil {
		return err
	}
	var env Environment
	return global.GVA_DB.Where("id = ? and project_id = ?", environmentID, projectID).First(&env).Delete(&env).Error
}
