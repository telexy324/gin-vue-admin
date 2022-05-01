package ansible

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

type InventoryService struct {
}

var InventoryServiceApp = new(InventoryService)

func (inventoryService *InventoryService) FillInventory(inventory *ansible.Inventory) (err error) {
	if inventory.SSHKeyID != nil {
		inventory.SSHKey, err = KeyServiceApp.GetAccessKey(float64(inventory.ProjectID), float64(*inventory.SSHKeyID))
	}

	if err != nil {
		return
	}

	if inventory.BecomeKeyID != nil {
		inventory.BecomeKey, err = KeyServiceApp.GetAccessKey(float64(inventory.ProjectID), float64(*inventory.BecomeKeyID))
	}

	return
}

func (inventoryService *InventoryService) GetInventory(projectID float64, inventoryID float64) (inventory ansible.Inventory, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, inventoryID).First(&inventory).Error
	if err != nil {
		return
	}
	err = inventoryService.FillInventory(&inventory)
	return
}

func (inventoryService *InventoryService) GetInventories(info request.GetByProjectId) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var inventories []ansible.Inventory
	db := global.GVA_DB.Model(&ansible.Inventory{})
	order := ""
	if info.SortInverted {
		order = "desc"
	}
	db = db.Where("project_id=?", info.ProjectId).Order(info.SortBy + " " + order)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&inventories).Error
	return err, inventories, total
}

//func (m *Inventory) GetInventoryRefs(projectID int, inventoryID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.InventoryProps, inventoryID)
//}

func (inventoryService *InventoryService) DeleteInventory(projectID float64, inventoryID float64) error {
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

	switch inventory.Type {
	case ansible.InventoryStatic:
		break
	case ansible.InventoryFile:
		if !IsValidInventoryPath(inventory.Inventory) {
			return errors.New("Inventory path is not valid!")
		}
	default:
		return errors.New("Inventory type is not valid!")
	}

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

// IsValidInventoryPath tests a path to ensure it is below the cwd
func IsValidInventoryPath(path string) bool {

	currentPath, err := os.Getwd()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(currentPath, absPath)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(relPath, "..")
}