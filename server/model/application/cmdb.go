package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname     string `json:"hostname" gorm:"column:hostname"`    // 主机名
	Architecture int    `json:"architecture" gorm:"column:architecture"`    // 架构
	ManageIp     string `json:"manage_ip" gorm:"column:manage_ip"`       // 管理ip
	Os           bool   `json:"os" gorm:"column:os"`            // 系统
	OsVersion    string `json:"os_version" gorm:"column:os_version"` // 系统版本
}
