package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type InventoriesResponse struct {
	Inventories []ansible.Inventory `json:"Inventories"`
}

type InventoryResponse struct {
	Inventory ansible.Inventory `json:"Inventory"`
}