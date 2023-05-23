package logUploadSvr

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl/request"
	"gorm.io/gorm"
	"strings"
)

type ServerService struct {
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddServer
//@description: 添加服务器
//@param: server model.Server
//@return: error

func (serverService *ServerService) AddServer(server logUploadMdl.Server) error {
	if !errors.Is(global.GVA_DB.Where("hostname = ?", server.Hostname).First(&logUploadMdl.Server{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复hostname，请修改name")
	}
	return global.GVA_DB.Create(&server).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteServer
//@description: 删除服务器
//@param: id float64
//@return: err error

func (serverService *ServerService) DeleteServer(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&logUploadMdl.Server{}).Error
	if err != nil {
		return
	}
	var server logUploadMdl.Server
	return global.GVA_DB.Where("id = ?", id).First(&server).Delete(&server).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteServerByIds
//@description: 批量删除服务器
//@param: servers []model.Servers
//@return: err error

func (serverService *ServerService) DeleteServerByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]logUploadMdl.Server{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateServer
//@description: 更新路由
//@param: server model.Server
//@return: err error

func (serverService *ServerService) UpdateServer(server logUploadMdl.Server) (err error) {
	var oldServer logUploadMdl.Server
	upDateMap := make(map[string]interface{})
	upDateMap["hostname"] = server.Hostname
	upDateMap["manage_ip"] = server.ManageIp
	upDateMap["mode"] = server.Mode
	upDateMap["port"] = server.Port

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", server.ID).Find(&oldServer)
		if oldServer.Hostname != server.Hostname {
			if !errors.Is(tx.Where("id <> ? AND hostname = ?", server.ID, server.Hostname).First(&logUploadMdl.Server{}).Error, gorm.ErrRecordNotFound) {
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

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerById
//@description: 返回当前选中server
//@param: id float64
//@return: err error, server model.Server

func (serverService *ServerService) GetServerById(id float64) (err error, server logUploadMdl.Server) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerList
//@description: 获取服务器分页
//@return: err error, list interface{}, total int64

func (serverService *ServerService) GetServerList(info request2.ServerSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var serverList []logUploadMdl.Server
	db := global.GVA_DB.Model(&logUploadMdl.Server{})
	if info.Hostname != "" {
		hostname := strings.Trim(info.Hostname, " ")
		db = db.Where("`hostname` LIKE ?", "%"+hostname+"%")
	}
	if info.ManageIp != "" {
		manageIp := strings.Trim(info.ManageIp, " ")
		db = db.Where("`manage_ip` LIKE ?", "%"+manageIp+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&serverList).Error
	return err, serverList, total
}
