package logUploadMdl

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	global.GVA_MODEL
	Hostname string   `json:"hostname" gorm:"type:varchar(100);not null;default:'';column:hostname"` // 主机名
	ManageIp string   `json:"manageIp" gorm:"type:varchar(15);not null;default:'';column:manage_ip"` // 管理ip
	Mode     int      `json:"mode" gorm:"type:tinyint(2);not null;default:0;column:mode"`            // 方式 1 ftp 2 sftp
	Port     int      `json:"port" gorm:"type:int(5);not null;default:0;column:port"`                // 上传端口
	Secrets  []Secret `json:"secrets" gorm:"-"`
}

func (m *Server) TableName() string {
	return "application_upload_servers"
}

func (m *Server) AfterFind(tx *gorm.DB) (err error) {
	secrets := make([]Secret, 0)
	if err = tx.Model(&Secret{}).Where("server_id = ?", m.ID).Find(&m.Secrets).Error; err != nil {
		global.GVA_LOG.Error("转换失败", zap.Any("err", err))
		return
	}
	m.Secrets = secrets
	return nil
}

type Secret struct {
	global.GVA_MODEL
	ServerId int    `json:"serverId" gorm:"type:bigint;not null;default:0;column:server_id"`       // 主机ID
	Name     string `json:"name" gorm:"type:varchar(100);not null;default:'';column:name"`         // 用户名
	Password string `json:"password" gorm:"type:varchar(100);not null;default:'';column:password"` // 密码
}

func (m *Secret) TableName() string {
	return "application_upload_secrets"
}
