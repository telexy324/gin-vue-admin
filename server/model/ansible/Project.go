package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
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
