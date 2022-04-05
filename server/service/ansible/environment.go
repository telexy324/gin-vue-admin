package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"gorm.io/gorm"
)

type EnvironmentService struct {
}

var EnvironmentServiceApp = new(EnvironmentService)

func (environmentService *EnvironmentService) GetEnvironment(projectID int, environmentID int) (ansible.Environment, error) {
	var environment ansible.Environment
	err := global.GVA_DB.Where("project_id=? and id =?", projectID, environmentID).First(&environment).Error
	return environment, err
}

//func (d *SqlDb) GetEnvironmentRefs(projectID int, environmentID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.EnvironmentProps, environmentID)
//}

func (environmentService *EnvironmentService) GetEnvironments(projectID int, sortInverted bool, sortBy string) ([]ansible.Environment, error) {
	var environments []ansible.Environment
	db := global.GVA_DB.Model(&ansible.Environment{})
	order := ""
	if sortInverted {
		order = "desc"
	}
	db = db.Where("project_id=?", projectID).Order(sortBy + " " + order)
	err := db.Find(&environments).Error
	return environments, err
}

func (environmentService *EnvironmentService) UpdateEnvironment(env ansible.Environment) error {
	err := env.Validate()
	if err != nil {
		return err
	}
	var oldEnv ansible.Environment
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

func (environmentService *EnvironmentService) CreateEnvironment(env ansible.Environment) (newEnv ansible.Environment, err error) {
	err = env.Validate()
	if err != nil {
		return
	}
	err = global.GVA_DB.Create(&env).Error
	return env, err
}

func (environmentService *EnvironmentService) DeleteEnvironment(projectID int, environmentID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", environmentID, projectID).First(&ansible.Environment{}).Error
	if err != nil {
		return err
	}
	var env ansible.Environment
	return global.GVA_DB.Where("id = ? and project_id = ?", environmentID, projectID).First(&env).Delete(&env).Error
}
