package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	"gorm.io/gorm"
	"time"
)

type ProjectService struct {
}

var ProjectServiceApp = new(ProjectService)

func (projectService *ProjectService) CreateProject(project ansible.Project) (newProject ansible.Project, err error) {
	project.Created = time.Now()
	err = global.GVA_DB.Create(&project).Error
	return project, err
}

func (projectService *ProjectService) GetProjects(userID int,info request.GetById) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var projects []ansible.Project
	db := global.GVA_DB.Model(&ansible.Project{})
	db = db.Joins("inner join ansible_project_users on ansible_projects.id = ansible_project_users.project_id").
		Where("ansible_project_users.user_id=?", userID)
	if info.SortBy != "" {
		order := ""
		if info.SortInverted {
			order = "desc"
		}
		db = db.Order(info.SortBy + " " + order)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&projects).Error
	return err, projects, total
}

func (projectService *ProjectService) GetProject(projectID int) (project ansible.Project, err error) {
	err = global.GVA_DB.Where("id =?", projectID).First(&project).Error
	return
}

func (projectService *ProjectService) DeleteProject(projectID int) error {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if txErr := global.GVA_DB.Where("id = ?", projectID).First(&ansible.Project{}).Error; txErr != nil {
			return txErr
		}
		var project ansible.Project
		if txErr := global.GVA_DB.Where("id = ?", projectID).First(&project).Delete(&project).Error; txErr != nil {
			return txErr
		}
		var templates []ansible.Template
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&templates).Delete(&templates).Error; txErr != nil {
			return txErr
		}
		var inventories []ansible.Inventory
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&inventories).Delete(&inventories).Error; txErr != nil {
			return txErr
		}
		var users []ansible.ProjectUser
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&users).Delete(&users).Error; txErr != nil {
			return txErr
		}
		var keys []ansible.AccessKey
		if txErr := global.GVA_DB.Where("project_id = ? ", projectID).Find(&keys).Delete(&keys).Error; txErr != nil {
			return txErr
		}
		return nil
	})
	return err
}

func (projectService *ProjectService) UpdateProject(project ansible.Project) error {
	var oldProject ansible.Project
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
