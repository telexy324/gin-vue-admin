package application

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname       string `json:"hostname" gorm:"column:hostname"`              // 主机名
	Architecture   int    `json:"architecture" gorm:"column:architecture"`      // 架构
	ManageIp       string `json:"manageIp" gorm:"column:manage_ip"`             // 管理ip
	Os             int    `json:"os" gorm:"column:os"`                          // 系统
	OsVersion      string `json:"osVersion" gorm:"column:os_version"`           // 系统版本
	SystemId       int    `json:"systemId" gorm:"column:system_id"`             // 所属系统id
	ApplicationIds string `json:"applicationIds" gorm:"column:application_ids"` // 安装应用列表
	Applications   []int  `json:"applications"`                                 // 安装应用列表
}

func (m *ApplicationServer) AfterFind(tx *gorm.DB) {
	if m.ApplicationIds != "" {
		if err := json.Unmarshal([]byte(m.ApplicationIds), m.Applications); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
		}
	}
}

type ServerRelation struct {
	global.GVA_MODEL
	StartServerId int               `json:"startServerId" gorm:"column:start_server_id"` // 源节点id
	EndServerId   int               `json:"endServerId" gorm:"column:end_server_id"`     // 目的节点id',
	EndServerUrl  string            `json:"endServerUrl" gorm:"column:end_server_url"`   // 目的节点url',
	Relation      string            `json:"relation" gorm:"column:relation"`             // 调用关系',
	StartServer   ApplicationServer `json:"startServer"`
	EndServer     ApplicationServer `json:"endServer"`
}

func (m *ServerRelation) TableName() string {
	return "application_server_relations"
}

type Node struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"` //0 outer 1 inner
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type SystemRelation struct {
	global.GVA_MODEL
	StartSystemId int               `json:"startSystemId" gorm:"column:start_system_id"` // 源节点id
	EndSystemId   int               `json:"endSystemId" gorm:"column:end_system_id"`     // 目的节点id',
	EndSystemUrl  string            `json:"endSystemUrl" gorm:"column:end_system_url"`   // 目的节点url',
	Relation      string            `json:"relation" gorm:"column:relation"`             // 调用关系',
	StartSystem   ApplicationServer `json:"startSystem"`
	EndSystem     ApplicationServer `json:"endSystem"`
}

func (m *SystemRelation) TableName() string {
	return "application_system_relations"
}

type ApplicationSystem struct {
	global.GVA_MODEL
	Name     string `json:"name" gorm:"column:name"`         // 主机名
	Position int    `json:"position" gorm:"column:position"` // 系统位置
}

type ApplicationSystemAdmin struct {
	global.GVA_MODEL
	SystemId int `json:"systemId" gorm:"column:system_id"` // 系统id
	AdminId  int `json:"adminId" gorm:"column:admin_id"`   // 管理员id
	IsPrime  int `json:"isPrime" gorm:"column:is_prime"`   // 0 非主管 1 主管
}

type Application struct {
	global.GVA_MODEL
	ApplicationType int    `json:"type" gorm:"column:type"`       // 应用类型 0 未定义 1 数据库 2 缓存 3 应用 4 存储 5 负载均衡 6 备份 7 反向代理
	Name            string `json:"name" gorm:"column:name"`       // 应用名称
	Version         string `json:"version" gorm:"column:version"` // 版本
}
