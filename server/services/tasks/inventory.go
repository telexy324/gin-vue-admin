package tasks

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"io/ioutil"
	"strconv"
)

func (t *TaskRunner) installInventory() (err error) {
	if t.inventory.SSHKeyID != nil {
		err = keyService.Install(&t.inventory.SSHKey, ansible.AccessKeyRoleAnsibleUser)
		if err != nil {
			return
		}
	}

	if t.inventory.BecomeKeyID != nil {
		err = keyService.Install(&t.inventory.BecomeKey, ansible.AccessKeyRoleAnsibleBecomeUser)
		if err != nil {
			return
		}
	}

	if t.inventory.Type == ansible.InventoryStatic {
		err = t.installStaticInventory()
	}

	return
}

func (t *TaskRunner) installStaticInventory() error {
	t.Log("installing static inventory")

	// create inventory file
	return ioutil.WriteFile(global.GVA_CONFIG.Ansible.TmpPath+"/inventory_"+strconv.Itoa(int(t.task.ID)), []byte(t.inventory.Inventory), 0664)
}
