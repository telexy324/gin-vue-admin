package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

const (
	NotAdmin = iota
	IsAdmin
)

// Project is the top level structure in Semaphore
type Project struct {
	global.GVA_MODEL
	Name             string    `gorm:"column:name" json:"name" binding:"required"`
	Created          time.Time `gorm:"column:created" json:"created"`
	Alert            bool      `gorm:"column:alert" json:"alert"`
	AlertChat        *string   `gorm:"column:alert_chat" json:"alertChat"`
	MaxParallelTasks int       `gorm:"column:max_parallel_tasks" json:"maxParallelTasks"`
}

type ProjectUser struct {
	global.GVA_MODEL
	ProjectId int `json:"projectId" gorm:"column:project_id"` // 系统id
	UserId    int `json:"userId" gorm:"column:user_id"`       // 管理员id
	Admin     int `json:"admin" gorm:"column:admin"`          // 0 非主管 1 主管
}

func (m *Project) TableName() string {
	return "ansible_projects"
}

func (m *ProjectUser) TableName() string {
	return "ansible_project_users"
}
