package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"gorm.io/gorm"
)

type InventoryService struct {
}

var InventoryServiceApp = new(InventoryService)

func (inventoryService *InventoryService) FillInventory(inventory *ansible.Inventory) (err error) {
	if inventory.SSHKeyID != nil {
		inventory.SSHKey, err = KeyServiceApp.GetAccessKey(inventory.ProjectID, *inventory.SSHKeyID)
	}

	if err != nil {
		return
	}

	if inventory.BecomeKeyID != nil {
		inventory.BecomeKey, err = KeyServiceApp.GetAccessKey(inventory.ProjectID, *inventory.BecomeKeyID)
	}

	return
}

func (inventoryService *InventoryService) GetInventory(projectID int, inventoryID int) (inventory ansible.Inventory, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, inventoryID).First(&inventory).Error
	if err != nil {
		return
	}
	err = inventoryService.FillInventory(&inventory)
	return
}

func (inventoryService *InventoryService) GetInventories(projectID int, sortInverted bool, sortBy string) ([]ansible.Inventory, error) {
	var inventories []ansible.Inventory
	db := global.GVA_DB.Model(&ansible.Inventory{})
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

func (inventoryService *InventoryService) DeleteInventory(projectID int, inventoryID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", inventoryID, projectID).First(&ansible.Inventory{}).Error
	if err != nil {
		return err
	}
	var inventory ansible.Inventory
	return global.GVA_DB.Where("id = ? and project_id = ?", inventoryID, projectID).First(&inventory).Delete(&inventory).Error
}

func (inventoryService *InventoryService) UpdateInventory(inventory ansible.Inventory) error {
	var oldInventory ansible.Inventory
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

func (inventoryService *InventoryService) CreateInventory(inventory ansible.Inventory) (newInventory ansible.Inventory, err error) {
	err = global.GVA_DB.Create(&inventory).Error
	return inventory, err
}
