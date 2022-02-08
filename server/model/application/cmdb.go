package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname     string `json:"hostname" gorm:"column:hostname"`         // 主机名
	Architecture int    `json:"architecture" gorm:"column:architecture"` // 架构
	ManageIp     string `json:"manageIp" gorm:"column:manage_ip"`        // 管理ip
	Os           int    `json:"os" gorm:"column:os"`                     // 系统
	OsVersion    string `json:"osVersion" gorm:"column:os_version"`      // 系统版本
}

type SystemRelation struct {
	global.GVA_MODEL
	StartServerId       int    `json:"startServerId" gorm:"column:start_server_id"`             // 源节点id
	StartServerName     string `json:"startServerName" gorm:"column:start_server_name"`         // 源节点名称',
	StartServerPosition int    `json:"startServerPosition" gorm:"column:start_server_position"` // 源节点是否系统外 0 系统外, 1 系统内',
	EndServerId         int    `json:"endServerId" gorm:"column:end_server_id"`                 // 目的节点id',
	EndServerName       string `json:"endServerName" gorm:"column:end_server_name"`             // 目的节点名称',
	EndServerPosition   int    `json:"endServerPosition" gorm:"column:end_server_position"`     // 目的节点是否系统外 0 系统外, 1 系统内',
	EndServerIp         string `json:"endServerIp" gorm:"column:end_server_ip"`                 // 目的节点ip',
	EndServerPort       int    `json:"endServerPort" gorm:"column:end_server_port"`             // 目的节点端口',
	Relation            string `json:"relation" gorm:"column:relation"`                         // 调用关系',
}
