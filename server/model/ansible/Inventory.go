package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

const (
	InventoryStatic = "static"
	InventoryFile   = "file"
)

// Inventory is the model of an ansible inventory file
type Inventory struct {
	global.GVA_MODEL
	Name      string `gorm:"name" json:"name" binding:"required"`
	ProjectID int    `gorm:"project_id" json:"project_id"`
	Inventory string `gorm:"inventory" json:"inventory"`

	// accesses hosts in inventory
	SSHKeyID *int      `gorm:"ssh_key_id" json:"ssh_key_id"`
	SSHKey   AccessKey `gorm:"-" json:"-"`

	BecomeKeyID *int      `gorm:"become_key_id" json:"become_key_id"`
	BecomeKey   AccessKey `gorm:"-" json:"-"`

	// static/file
	Type string `gorm:"type" json:"type"`
}

func (m *Inventory) TableName() string {
	return "ansible_inventorys"
}
