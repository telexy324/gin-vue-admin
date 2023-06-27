package application

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplicationServer struct {
	global.GVA_MODEL
	Hostname     string `json:"hostname" gorm:"type:varchar(100);not null;default:'';column:hostname"`       // 主机名
	Architecture int    `json:"architecture" gorm:"type:tinyint(2);not null;default:0;column:architecture"`  // 架构
	ManageIp     string `json:"manageIp" gorm:"type:varchar(15);not null;default:'';column:manage_ip"`       // 管理ip
	Os           int    `json:"os" gorm:"type:tinyint(2);not null;default:0;column:os"`                      // 系统
	OsVersion    string `json:"osVersion" gorm:"type:varchar(50);not null;default:'';column:os_version"`     // 系统版本
	SystemId     int    `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id"`             // 所属系统id
	AppIds       string `json:"appIds" gorm:"type:varchar(100);not null;default:'';column:app_ids"`          // 安装应用列表
	Apps         []App  `json:"apps" gorm:"-"`                                                               // 安装应用列表
	SshPort      int    `json:"sshPort" gorm:"type:int(5);not null;default:0;column:ssh_port"`               // ssh端口、
	DisplayName  string `json:"displayName" gorm:"type:varchar(50);not null;default:'';column:display_name"` // 展示名
}

func (m *ApplicationServer) AfterFind(tx *gorm.DB) (err error) {
	appIds := make([]int, 0)
	if m.AppIds != "" {
		if err = json.Unmarshal([]byte(m.AppIds), &appIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if err = tx.Model(&App{}).Where("id = ?", appIds).Find(&m.Apps).Error; err != nil {
		global.GVA_LOG.Error("转换失败", zap.Any("err", err))
		return
	}
	return nil
}

type ServerRelation struct {
	global.GVA_MODEL
	StartServerId int               `json:"startServerId" gorm:"type:bigint;not null;default:0;column:start_server_id"`      // 源节点id
	EndServerId   int               `json:"endServerId" gorm:"type:bigint;not null;default:0;column:end_server_id"`          // 目的节点id',
	EndServerUrl  string            `json:"endServerUrl" gorm:"type:varchar(255);not null;default:'';column:end_server_url"` // 目的节点url',
	Relation      string            `json:"relation" gorm:"type:varchar(100);not null;default:'';column:relation"`           // 调用关系',
	StartServer   ApplicationServer `json:"startServer"`
	EndServer     ApplicationServer `json:"endServer"`
}

func (m *ServerRelation) TableName() string {
	return "application_server_relations"
}

type Node struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type SystemRelation struct {
	global.GVA_MODEL
	StartSystemId int               `json:"startSystemId" gorm:"type:bigint;not null;default:0;column:start_system_id"`      // 源节点id
	EndSystemId   int               `json:"endSystemId" gorm:"type:bigint;not null;default:0;column:end_system_id"`          // 目的节点id',
	EndSystemUrl  string            `json:"endSystemUrl" gorm:"type:varchar(255);not null;default:'';column:end_system_url"` // 目的节点url',
	Relation      string            `json:"relation" gorm:"type:varchar(100);not null;default:'';column:relation"`           // 调用关系',
	StartSystem   ApplicationSystem `json:"startSystem" gorm:"-"`
	EndSystem     ApplicationSystem `json:"endSystem" gorm:"-"`
}

func (m *SystemRelation) TableName() string {
	return "application_system_relations"
}

type ApplicationSystem struct {
	global.GVA_MODEL
	Name     string   `json:"name" gorm:"type:varchar(100);not null;default:'';column:name"`      // 系统名
	Position int      `json:"position" gorm:"type:tinyint(2);not null;default:0;column:position"` // 系统位置
	Network  string   `json:"network" gorm:"type:text;column:network"`                            // 系统位置
	SshUser  string   `json:"sshUser" gorm:"type:text;column:ssh_user"`                           // 系统位置
	SshUsers []string `json:"sshUsers" gorm:"-"`
}

func (m *ApplicationSystem) AfterFind(tx *gorm.DB) (err error) {
	users := make([]string, 0)
	if m.SshUser != "" {
		if err = json.Unmarshal([]byte(m.SshUser), &users); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
		m.SshUsers = users
	}
	return nil
}

type ApplicationSystemAdmin struct {
	global.GVA_MODEL
	SystemId int `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id"` // 系统id
	AdminId  int `json:"adminId" gorm:"type:bigint;not null;default:0;column:admin_id"`   // 管理员id
	//IsPrime  int `json:"isPrime" gorm:"column:is_prime"`   // 0 非主管 1 主管
}

type App struct {
	global.GVA_MODEL
	ApplicationType int    `json:"type" gorm:"type:tinyint(2);not null;default:0;column:type"`          // 应用类型 0 未定义 1 数据库 2 缓存 3 应用 4 存储 5 负载均衡 6 备份 7 反向代理
	Name            string `json:"name" gorm:"type:varchar(100);not null;default:'';column:name"`       // 应用名称
	Version         string `json:"version" gorm:"type:varchar(100);not null;default:'';column:version"` // 版本
}

func (m *App) TableName() string {
	return "application_apps"
}

type ApplicationSystemSysAdmin struct {
	global.GVA_MODEL
	SystemId int `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id"` // 系统id
	AdminId  int `json:"adminId" gorm:"type:bigint;not null;default:0;column:admin_id"`   // 管理员id
}

type ApplicationSystemEditRelation struct {
	global.GVA_MODEL
	SystemId int    `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id"` // 系统id
	Relation string `json:"relation" gorm:"type:text;column:relation"`                       // 调用关系',
}
