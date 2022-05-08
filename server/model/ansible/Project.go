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
