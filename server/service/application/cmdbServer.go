package application

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
	"strings"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getMenuTreeMap
//@description: 获取路由总树map
//@param: authorityId string
//@return: err error, treeMap map[string][]model.SysMenu

type CmdbServerService struct {
}

var CmdbServerServiceApp = new(CmdbServerService)

//@author: [telexy324](https://github.com/telexy324)
//@function: AddServer
//@description: 添加服务器
//@param: server model.ApplicationServer
//@return: error

func (cmdbServerService *CmdbServerService) AddServer(server application.ApplicationServer) error {
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

func (cmdbServerService *CmdbServerService) DeleteServer(id float64) (err error) {
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

func (cmdbServerService *CmdbServerService) UpdateServer(server application.ApplicationServer) (err error) {
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

func (cmdbServerService *CmdbServerService) GetServerById(id float64) (err error, server application.ApplicationServer) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerList
//@description: 获取服务器分页
//@return: err error, list interface{}, total int64

func (cmdbServerService *CmdbServerService) GetServerList(info request2.ServerSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var serverList []application.ApplicationServer
	db := global.GVA_DB.Model(&application.ApplicationServer{})
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

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemServers
//@description: 获取系统内全部服务器
//@return: err error, serverList []application.ApplicationServer
func (cmdbServerService *CmdbServerService) GetSystemServers(systemId float64) (err error, serverList []application.ApplicationServer) {
	db := global.GVA_DB.Model(&application.ApplicationServer{})
	err = db.Where("system_id = ?", systemId).Find(&serverList).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddRelation
//@description: 添加联系
//@param: relation model.ServerRelation
//@return: error

func (cmdbServerService *CmdbServerService) AddRelation(relation application.ServerRelation) error {
	if !errors.Is(global.GVA_DB.Where("start_server_id = ? and end_server_id = ?", relation.StartServerId, relation.EndServerId).First(&application.ServerRelation{}).Error, gorm.ErrRecordNotFound) {
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
//@function: ServerRelations
//@description: 返回当前选中server的关系路径
//@param: id float64
//@return: err error, server model.ApplicationServer

func (cmdbServerService *CmdbServerService) ServerRelations(id float64) (err error, relations []application.ServerRelation, nodes []application.Node) {
	server := application.ApplicationServer{}
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	if err != nil {
		return
	}
	mapNodes := make(map[int]bool)
	nodes = append(nodes, application.Node{
		Id:   int(server.ID),
		Name: server.Hostname,
	})
	mapNodes[int(server.ID)] = true
	relationOneSrc := make([]application.ServerRelation, 0)
	relationOneDest := make([]application.ServerRelation, 0)
	err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", id).Find(&relationOneSrc).Error
	err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", id).Find(&relationOneDest).Error
	if len(relationOneSrc) > 0 {
		relations = append(relations, relationOneSrc...)
		for _, relation := range relationOneSrc {
			if mapNodes[relation.EndServerId] == false {
				nodes = append(nodes, application.Node{
					Id:   relation.EndServerId,
					Name: relation.EndServer.Hostname,
				})
				mapNodes[relation.EndServerId] = true
			}
			relationTwoSrc := make([]application.ServerRelation, 0)
			relationTwoDest := make([]application.ServerRelation, 0)
			err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", relation.EndServerId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", relation.EndServerId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndServerId,
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
					Name: relation.StartServer.Hostname,
				})
				mapNodes[relation.StartServerId] = true
			}
			relationTwoSrc := make([]application.ServerRelation, 0)
			relationTwoDest := make([]application.ServerRelation, 0)
			err = global.GVA_DB.Preload("EndServer").Where("start_server_id = ?", relation.StartServerId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartServer").Where("end_server_id = ?", relation.StartServerId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndServerId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndServerId,
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

func (cmdbServerService *CmdbServerService) ParseInfoList2Excel(infoList []application.ApplicationServer, headers []string, filePath string) error {
	excel := excelize.NewFile()
	err := excel.SetSheetRow("Sheet1", "A1", &headers)
	if err != nil {
		return err
	}
	for i, server := range infoList {
		axis := fmt.Sprintf("A%d", i+2)
		err = excel.SetSheetRow("Sheet1", axis, &[]interface{}{
			server.ID,
			server.Hostname,
			server.Architecture,
			server.ManageIp,
			server.Os,
			server.OsVersion,
		})
		if err != nil {
			global.GVA_LOG.Error("转换Excel行失败!，server id为: "+strconv.Itoa(int(server.ID)), zap.Any("err", err))
		}
	}
	err = excel.SaveAs(filePath)
	return err
}

func (cmdbServerService *CmdbServerService) ImportExcel2db(file multipart.File, header *multipart.FileHeader) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	if len(rows) <= 1 {
		return errors.New("数据表内容为空")
	}
	rows = rows[1:]
	for _, row := range rows {
		server := application.ApplicationServer{
			Hostname:     row[0],
			Architecture: int(consts.OsMapReverse[row[1]]),
			ManageIp:     row[2],
			Os:           int(consts.OsMapReverse[row[3]]),
			OsVersion:    row[4],
		}
		if !errors.Is(global.GVA_DB.Where("hostname = ?", server.Hostname).First(&application.ApplicationServer{}).Error, gorm.ErrRecordNotFound) {
			global.GVA_LOG.Error("存在重复hostname，请修改name" + row[0])
			continue
		}
		if err = global.GVA_DB.Create(&server).Error; err != nil {
			global.GVA_LOG.Error("插入失败", zap.Any("err", err))
		}
	}
	return err
}

func (cmdbServerService *CmdbServerService) ExportTemplate() (*excelize.File, error) {
	excel := excelize.NewFile()
	sheetName := "Sheet1"
	headers := []string{"主机名", "架构", "管理IP", "系统类型", "系统版本"}
	architectures := []string{consts.ArchitectureMap[consts.ArchitectureX86], consts.ArchitectureMap[consts.ArchitectureArm]}
	oses := []string{consts.OsMap[consts.OsRedhat], consts.OsMap[consts.OsSuse], consts.OsMap[consts.OsCentos], consts.OsMap[consts.OsKylin]}
	err := excel.SetSheetRow("Sheet1", "A1", &headers)
	if err != nil {
		return nil, err
	}
	dvRangeArchitecture := excelize.NewDataValidation(true)
	dvRangeArchitecture.Sqref = "B2:B255"
	if err = dvRangeArchitecture.SetDropList(architectures); err != nil {
		return nil, err
	}
	if err = excel.AddDataValidation(sheetName, dvRangeArchitecture); err != nil {
		return nil, err
	}
	dvRangeOs := excelize.NewDataValidation(true)
	dvRangeOs.Sqref = "D2:D255"
	if err = dvRangeOs.SetDropList(oses); err != nil {
		return nil, err
	}
	if err = excel.AddDataValidation(sheetName, dvRangeOs); err != nil {
		return nil, err
	}
	return excel, err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddApp
//@description: 添加应用
//@param: server model.App
//@return: error

func (cmdbServerService *CmdbServerService) AddApp(app application.App) error {
	if !errors.Is(global.GVA_DB.Where("name = ?", app.Name).First(&application.App{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.GVA_DB.Create(&app).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteApp
//@description: 删除应用
//@param: id float64
//@return: err error

func (cmdbServerService *CmdbServerService) DeleteApp(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.App{}).Error
	if err != nil {
		return
	}
	var app application.App
	return global.GVA_DB.Where("id = ?", id).First(&app).Delete(&app).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateApp
//@description: 更新应用
//@param: server model.App
//@return: err error

func (cmdbServerService *CmdbServerService) UpdateApp(app application.App) (err error) {
	var oldApp application.App
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = app.Name
	upDateMap["application_type"] = app.ApplicationType
	upDateMap["version"] = app.Version

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", app.ID).Find(&oldApp)
		if oldApp.Name != app.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", app.ID, app.Name).First(&application.App{}).Error, gorm.ErrRecordNotFound) {
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
//@function: GetAppById
//@description: 返回当前选中app
//@param: id float64
//@return: err error, server model.App

func (cmdbServerService *CmdbServerService) GetAppById(id float64) (err error, app application.App) {
	err = global.GVA_DB.Where("id = ?", id).First(&app).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetAppList
//@description: 获取应用分页
//@return: err error, list interface{}, total int64

func (cmdbServerService *CmdbServerService) GetAppList(info request2.AppSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var appList []application.App
	db := global.GVA_DB.Model(&application.App{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`hostname` LIKE ?", "%"+name+"%")
	}
	if info.ApplicationType > 0 {
		db = db.Where("type = ?", info.ApplicationType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&appList).Error
	return err, appList, total
}