package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

const (
	InventoryStatic = "static"
	InventoryFile   = "file"
)

// Inventory is the model of an ansible inventory file
type Inventory struct {
	global.GVA_MODEL
	Name      string `gorm:"column:name" json:"name" binding:"required"`
	ProjectID int    `gorm:"column:project_id" json:"projectId"`
	Inventory string `gorm:"column:inventory" json:"inventory"`

	// accesses hosts in inventory
	SSHKeyID *int      `gorm:"column:ssh_key_id" json:"sshKeyId"`
	SSHKey   AccessKey `gorm:"-" json:"-"`

	BecomeKeyID *int      `gorm:"column:become_key_id" json:"becomeKeyId"`
	BecomeKey   AccessKey `gorm:"-" json:"-"`

	// static/file
	Type string `gorm:"column:type" json:"type"`
}

func (m *Inventory) TableName() string {
	return "ansible_inventorys"
}
