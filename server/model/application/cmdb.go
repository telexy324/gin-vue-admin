package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname     string `json:"parentId" gorm:"column:hostname"`    // 主机名
	Architecture int    `json:"path" gorm:"column:architecture"`    // 架构
	ManageIp     string `json:"name" gorm:"column:manage_ip"`       // 管理ip
	Os           bool   `json:"hidden" gorm:"column:os"`            // 系统
	OsVersion    string `json:"component" gorm:"column:os_version"` // 系统版本
}
