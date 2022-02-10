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
	Position     int    `json:"position" gorm:"column:position"`         // 系统位置
}

type SystemRelation struct {
	global.GVA_MODEL
	StartServerId int               `json:"startServerId" gorm:"column:start_server_id"` // 源节点id
	EndServerId   int               `json:"endServerId" gorm:"column:end_server_id"`     // 目的节点id',
	EndServerUrl  string            `json:"endServerUrl" gorm:"column:end_server_url"`   // 目的节点url',
	Relation      string            `json:"relation" gorm:"column:relation"`             // 调用关系',
	StartServer   ApplicationServer `json:"startServer"`
	EndServer     ApplicationServer `json:"endServer"`
}

func (SystemRelation) TableName() string {
	return "application_server_relations"
}

type Node struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"` //0 outer 1 inner
	Name  string `json:"name"`
	Value int    `json:"value"`
}
