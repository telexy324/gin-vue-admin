package application

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
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
//@description: 添加服务器
//@param: menu model.ApplicationServer
//@return: error

func (cmdbService *CmdbService) AddServer(server application.ApplicationServer) error {
	if !errors.Is(global.GVA_DB.Where("hostname = ?", server.Hostname).First(&application.ApplicationServer{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复hostname，请修改name")
	}
	return global.GVA_DB.Create(&server).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteServer
//@description: 删除服务器
//@param: id float64
//@return: err error

func (cmdbService *CmdbService) DeleteServer(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.ApplicationServer{}).Error
	if err != nil {
		return
	}
	var server application.ApplicationServer
	return global.GVA_DB.Where("id = ?", id).First(&server).Delete(&server).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateServer
//@description: 更新路由
//@param: server model.ApplicationServer
//@return: err error

func (cmdbService *CmdbService) UpdateServer(server application.ApplicationServer) (err error) {
	var oldServer application.ApplicationServer
	upDateMap := make(map[string]interface{})
	upDateMap["hostname"] = server.Hostname
	upDateMap["manage_ip"] = server.ManageIp
	upDateMap["os"] = server.Os
	upDateMap["os_version"] = server.OsVersion
	upDateMap["architecture"] = server.Architecture

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", server.ID).Find(&oldServer)
		if oldServer.Hostname != server.Hostname {
			if !errors.Is(tx.Where("id <> ? AND hostname = ?", server.ID, server.Hostname).First(&application.ApplicationServer{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}

		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}
