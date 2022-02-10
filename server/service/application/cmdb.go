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

//@author: [telexy324](https://github.com/telexy324)
//@function: AddServer
//@description: 添加服务器
//@param: server model.ApplicationServer
//@return: error

func (cmdbService *CmdbService) AddServer(server application.ApplicationServer) error {
	if !errors.Is(global.GVA_DB.Where("hostname = ?", server.Hostname).First(&application.ApplicationServer{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复hostname，请修改name")
	}
	return global.GVA_DB.Create(&server).Error
}

//@author: [telexy324](https://github.com/telexy324)
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

//@author: [telexy324](https://github.com/telexy324)
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

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerById
//@description: 返回当前选中server
//@param: id float64
//@return: err error, server model.ApplicationServer

func (cmdbService *CmdbService) GetServerById(id float64) (err error, server application.ApplicationServer) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerList
//@description: 获取路由分页
//@return: err error, list interface{}, total int64

func (cmdbService *CmdbService) GetServerList() (err error, list interface{}, total int64) {
	var serverList []application.ApplicationServer
	err = global.GVA_DB.Find(&serverList).Error
	return err, serverList, int64(len(serverList))
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddRelation
//@description: 添加联系
//@param: relation model.SystemRelation
//@return: error

func (cmdbService *CmdbService) AddRelation(relation application.SystemRelation) error {
	if !errors.Is(global.GVA_DB.Where("start_server_id = ? and end_server_id = ?", relation.StartServerId, relation.EndServerId).First(&application.SystemRelation{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复关系")
	}

	startServer := application.ApplicationServer{}
	endServer := application.ApplicationServer{}
	err := global.GVA_DB.Where("id = ?", relation.StartServerId).First(&startServer).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Where("id = ?", relation.EndServerId).First(&endServer).Error
	if err != nil {
		return err
	}
	return global.GVA_DB.Create(&relation).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: SystemRelations
//@description: 返回当前选中server的关系路径
//@param: id float64
//@return: err error, server model.ApplicationServer

func (cmdbService *CmdbService) SystemRelations(id float64) (err error, relations []application.SystemRelation, nodes []application.Node) {
	server := application.ApplicationServer{}
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	if err != nil {
		return
	}
	mapNodes := make(map[int]bool)
	nodes = append(nodes, application.Node{
		Id:   int(server.ID),
		Type: server.Position,
		Name: server.Hostname,
	})
	mapNodes[int(server.ID)] = true
	relationOneSrc := make([]application.SystemRelation, 0)
	relationOneDest := make([]application.SystemRelation, 0)
	err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", id).Find(&relationOneSrc).Error
	err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", id).Find(&relationOneDest).Error
	if len(relationOneSrc) > 0 {
		relations = append(relations, relationOneSrc...)
		for _, relation := range relationOneSrc {
			if mapNodes[relation.EndServerId] == false {
				nodes = append(nodes, application.Node{
					Id:   relation.EndServerId,
					Type: relation.EndServer.Position,
					Name: relation.EndServer.Hostname,
				})
				mapNodes[relation.EndServerId] = true
			}
			relationTwoSrc := make([]application.SystemRelation, 0)
			relationTwoDest := make([]application.SystemRelation, 0)
			err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", relation.EndServerId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", relation.EndServerId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndServerId,
							Type: relation.EndServer.Position,
							Name: relation.EndServer.Hostname,
						})
						mapNodes[relation.EndServerId] = true
					}
				}
			}
			if len(relationTwoDest) > 0 {
				relations = append(relations, relationTwoDest...)
				for _, relation := range relationTwoDest {
					if mapNodes[relation.StartServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.StartServerId,
							Type: relation.StartServer.Position,
							Name: relation.StartServer.Hostname,
						})
						mapNodes[relation.StartServerId] = true
					}
				}
			}
		}
	}
	if len(relationOneDest) > 0 {
		relations = append(relations, relationOneDest...)
		for _, relation := range relationOneDest {
			if mapNodes[relation.StartServerId] == false {
				nodes = append(nodes, application.Node{
					Id:   relation.StartServerId,
					Type: relation.StartServer.Position,
					Name: relation.StartServer.Hostname,
				})
				mapNodes[relation.StartServerId] = true
			}
			relationTwoSrc := make([]application.SystemRelation, 0)
			relationTwoDest := make([]application.SystemRelation, 0)
			err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", relation.StartServerId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", relation.StartServerId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndServerId,
							Type: relation.EndServer.Position,
							Name: relation.EndServer.Hostname,
						})
						mapNodes[relation.EndServerId] = true
					}
				}
			}
			if len(relationTwoDest) > 0 {
				relations = append(relations, relationTwoDest...)
				for _, relation := range relationTwoDest {
					if mapNodes[relation.StartServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.StartServerId,
							Type: relation.StartServer.Position,
							Name: relation.StartServer.Hostname,
						})
						mapNodes[relation.StartServerId] = true
					}
				}
			}
		}
	}
	return
}
