package application

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getMenuTreeMap
//@description: 获取路由总树map
//@param: authorityId string
//@return: err error, treeMap map[string][]model.SysMenu

type CmdbService struct {
}

var CmdbServiceApp = new(CmdbService)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddServer
//@description: 添加基础路由
//@param: menu model.ApplicationServer
//@return: error

func (cmdbService *CmdbService) AddServer(server application.ApplicationServer) error {
	if !errors.Is(global.GVA_DB.Where("hostname = ?", server.Hostname).First(&application.ApplicationServer{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复hostname，请修改name")
	}
	return global.GVA_DB.Create(&server).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteBaseMenu
//@description: 删除基础路由
//@param: id float64
//@return: err error

func (cmdbService *CmdbService) DeleteServer(id float64) (err error) {
	err = global.GVA_DB.Preload("Parameters").Where("parent_id = ?", id).First(&system.SysBaseMenu{}).Error
	if err != nil {
		var menu system.SysBaseMenu
		db := global.GVA_DB.Preload("SysAuthoritys").Where("id = ?", id).First(&menu).Delete(&menu)
		err = global.GVA_DB.Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
		if err != nil {
			return err
		}
		if len(menu.SysAuthoritys) > 0 {
			err = global.GVA_DB.Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
		} else {
			err = db.Error
			if err != nil {
				return
			}
		}
	} else {
		return errors.New("此菜单存在子菜单不可删除")
	}
	return err
}