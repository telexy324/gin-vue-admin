package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"gorm.io/gorm"
)

type UserService struct {
}

var UserServiceApp = new(UserService)

func (userService *UserService) CreateProjectUser(projectUser ansible.ProjectUser) (newProjectUser ansible.ProjectUser, err error) {
	err = global.GVA_DB.Create(&projectUser).Error
	return projectUser, err
}

func (userService *UserService) GetProjectUser(projectID, userID float64) (projectUser ansible.ProjectUser, err error) {
	err = global.GVA_DB.Where("project_id=? and user_id =?", projectID, userID).First(&projectUser).Error
	return
}

func (userService *UserService) GetProjectUsers(projectID int) (users []ansible.ProjectUser, err error) {
	db := global.GVA_DB.Model(&ansible.ProjectUser{})
	db = db.Where("project_id=?", projectID)
	err = db.Find(&users).Error
	return
}

func (userService *UserService) UpdateProjectUser(projectUser ansible.ProjectUser) error {
	var oldUser ansible.ProjectUser
	upDateMap := make(map[string]interface{})
	upDateMap["project_id"] = projectUser.ProjectId
	upDateMap["user_id"] = projectUser.UserId
	upDateMap["admin"] = projectUser.Admin

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", projectUser.ID, projectUser.ProjectId).Find(&oldUser)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (userService *UserService) DeleteProjectUser(projectID, userID int) error {
	err := global.GVA_DB.Where("user_id = ? and project_id = ?", userID, projectID).First(&ansible.ProjectUser{}).Error
	if err != nil {
		return err
	}
	var user ansible.ProjectUser
	return global.GVA_DB.Where("user_id = ? and project_id = ?", userID, projectID).First(&user).Delete(&user).Error
}
