package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"time"
)

// Project is the top level structure in Semaphore
type Project struct {
	global.GVA_MODEL
	Name             string    `gorm:"name" json:"name" binding:"required"`
	Created          time.Time `gorm:"created" json:"created"`
	Alert            bool      `gorm:"alert" json:"alert"`
	AlertChat        *string   `gorm:"alert_chat" json:"alert_chat"`
	MaxParallelTasks int       `gorm:"max_parallel_tasks" json:"max_parallel_tasks"`
}

type ProjectUser struct {
	global.GVA_MODEL
	ProjectId int `json:"projectId" gorm:"column:project_id"` // 系统id
	UserId    int `json:"userId" gorm:"column:user_id"`       // 管理员id
	Admin     int `json:"admin" gorm:"column:admin"`          // 0 非主管 1 主管
}

func CreateProject(project Project) (newProject Project, err error) {
	project.Created = time.Now()
	err = global.GVA_DB.Create(&project).Error
	return project, err
}

func (m *Project) GetProjects(userID int) (projects []Project, err error) {
	err = global.GVA_DB.Joins("inner join project_user on id = project_user.project_id").
		Where("project_user.user_id=?", userID).
		Order("name ").Find(&projects).Error
	return
}

func (m *Project) GetProject(projectID int) (project Project, err error) {
	err = global.GVA_DB.Where("id =?", projectID).First(&project).Error
	return
}

func (m *Project) DeleteProject(projectID int) error {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if txErr := global.GVA_DB.Where("id = ?", projectID).First(&Project{}).Error; txErr != nil {
			return txErr
		}
		var project Project
		if txErr := global.GVA_DB.Where("id = ?", projectID).First(&project).Delete(&project).Error; txErr != nil {
			return txErr
		}
		var templates []Template
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&templates).Delete(&templates).Error; txErr != nil {
			return txErr
		}
		var inventories []Inventory
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&inventories).Delete(&inventories).Error; txErr != nil {
			return txErr
		}
		var users []ProjectUser
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&users).Delete(&users).Error; txErr != nil {
			return txErr
		}
		var keys []AccessKey
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&keys).Delete(&keys).Error; txErr != nil {
			return txErr
		}
		return nil
	})
	return err
}

func (m *Project) UpdateProject(project Project) error {
	var oldProject Project
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = project.Name
	upDateMap["alert"] = project.Alert
	upDateMap["alert_chat"] = project.AlertChat
	upDateMap["max_parallel_tasks"] = project.MaxParallelTasks

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", project.ID).Find(&oldProject)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}
