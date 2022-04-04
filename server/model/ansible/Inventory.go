package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

const (
	InventoryStatic = "static"
	InventoryFile   = "file"
)

// Inventory is the model of an ansible inventory file
type Inventory struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name" binding:"required"`
	ProjectID int    `db:"project_id" json:"project_id"`
	Inventory string `db:"inventory" json:"inventory"`

	// accesses hosts in inventory
	SSHKeyID *int      `db:"ssh_key_id" json:"ssh_key_id"`
	SSHKey   AccessKey `db:"-" json:"-"`

	BecomeKeyID *int      `db:"become_key_id" json:"become_key_id"`
	BecomeKey   AccessKey `db:"-" json:"-"`

	// static/file
	Type string `db:"type" json:"type"`
}

func FillInventory(inventory *Inventory) (err error) {
	k := &AccessKey{}
	if inventory.SSHKeyID != nil {
		inventory.SSHKey, err = k.GetAccessKey(inventory.ProjectID, *inventory.SSHKeyID)
	}

	if err != nil {
		return
	}

	if inventory.BecomeKeyID != nil {
		inventory.BecomeKey, err = k.GetAccessKey(inventory.ProjectID, *inventory.BecomeKeyID)
	}

	return
}

func (m *Inventory) GetInventory(projectID int, inventoryID int) (inventory Inventory, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, inventoryID).First(&inventory).Error
	if err != nil {
		return
	}
	err = FillInventory(&inventory)
	return
}

func (m *Inventory) GetInventories(projectID int, sortInverted bool, sortBy string) ([]Inventory, error) {
	var inventories []Inventory
	db := global.GVA_DB.Model(&Inventory{})
	order := ""
	if sortInverted {
		order = "desc"
	}
	db = db.Where("project_id=?", projectID).Order(sortBy + " " + order)
	err := db.Find(&inventories).Error
	return inventories, err
}

//func (m *Inventory) GetInventoryRefs(projectID int, inventoryID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.InventoryProps, inventoryID)
//}

func (m *Inventory) DeleteInventory(projectID int, inventoryID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", inventoryID, projectID).First(&Inventory{}).Error
	if err != nil {
		return err
	}
	var inventory Inventory
	return global.GVA_DB.Where("id = ? and project_id = ?", inventoryID, projectID).First(&inventory).Delete(&inventory).Error
}

func (m *Inventory) UpdateInventory(inventory Inventory) error {
	var oldInventory Inventory
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = inventory.Name
	upDateMap["secret"] = inventory.SSHKeyID
	upDateMap["ssh_key_id"] = inventory.Type
	upDateMap["inventory"] = inventory.Inventory
	upDateMap["become_key_id"] = inventory.BecomeKeyID

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", inventory.ID, inventory.ProjectID).Find(&oldInventory)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func CreateInventory(inventory Inventory) (newInventory Inventory, err error) {
	err = global.GVA_DB.Create(&inventory).Error
	return inventory, err
}
