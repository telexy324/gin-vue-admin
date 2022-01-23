package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname     string `json:"hostname" gorm:"column:hostname"`         // 主机名
	Architecture int    `json:"architecture" gorm:"column:architecture"` // 架构
	ManageIp     string `json:"manageIp" gorm:"column:manage_ip"`       // 管理ip
	Os           int    `json:"os" gorm:"column:os"`                     // 系统
	OsVersion    string `json:"osVersion" gorm:"column:os_version"`     // 系统版本
}
