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

type CmdbSystemService struct {
}

var CmdbSystemServiceApp = new(CmdbSystemService)

//@author: [telexy324](https://github.com/telexy324)
//@function: AddServer
//@description: 添加系统
//@param: server model.ApplicationServer
//@return: error

func (cmdbSystemService *CmdbSystemService) AddSystem(addSystemRequest request2.AddSystem) error {
	if !errors.Is(global.GVA_DB.Where("name = ?", addSystemRequest.System.Name).First(&application.ApplicationSystem{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	if addSystemRequest.SystemAdmin != nil && len(addSystemRequest.SystemAdmin) > 0 {
		for _, admin := range addSystemRequest.SystemAdmin {
			if err := global.GVA_DB.Create(&admin).Error; err != nil {
				global.GVA_LOG.Error("添加管理员失败", zap.Any("err", err))
			}
		}
	}
	return global.GVA_DB.Create(&addSystemRequest.System).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSystem
//@description: 删除系统
//@param: id float64
//@return: err error

func (cmdbSystemService *CmdbSystemService) DeleteSystem(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.ApplicationSystem{}).Error
	if err != nil {
		return
	}
	var system application.ApplicationSystem
	if err = global.GVA_DB.Where("id = ?", id).First(&system).Delete(&system).Error; err != nil {
		return err
	}
	var systemAdmins []application.ApplicationSystemAdmin
	err = global.GVA_DB.Where("system_id = ?", id).Find(&systemAdmins).Delete(&systemAdmins).Error
	if err != nil {
		return err
	}
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateSystem
//@description: 更新系统
//@param: server model.ApplicationSystem
//@return: err error

func (cmdbSystemService *CmdbSystemService) UpdateSystem(addSystemRequest request2.AddSystem) (err error) {
	var oldSystem application.ApplicationSystem
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = addSystemRequest.System.Name
	upDateMap["position"] = addSystemRequest.System.Position

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", addSystemRequest.System.ID).Find(&oldSystem)
		if oldSystem.Name != addSystemRequest.System.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", addSystemRequest.System.ID, addSystemRequest.System.Name).First(&application.ApplicationSystem{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}

		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}

		if addSystemRequest.SystemAdmin != nil && len(addSystemRequest.SystemAdmin) > 0 {
			for _, admin := range addSystemRequest.SystemAdmin {

			}
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

func (cmdbSystemService *CmdbSystemService) GetServerById(id float64) (err error, server application.ApplicationServer) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetServerList
//@description: 获取服务器分页
//@return: err error, list interface{}, total int64

func (cmdbSystemService *CmdbSystemService) GetServerList(info request2.ServerSearch) (err error, list interface{}, total int64) {
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
func (cmdbSystemService *CmdbSystemService) GetSystemServers(systemId float64) (err error, serverList []application.ApplicationServer) {
	db := global.GVA_DB.Model(&application.ApplicationServer{})
	err = db.Where("system_id = ?", systemId).Find(&serverList).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddRelation
//@description: 添加联系
//@param: relation model.ServerRelation
//@return: error

func (cmdbSystemService *CmdbSystemService) AddRelation(relation application.ServerRelation) error {
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

func (cmdbSystemService *CmdbSystemService) ServerRelations(id float64) (err error, relations []application.ServerRelation, nodes []application.Node) {
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
					Type: relation.EndServer.Position,
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

func (cmdbSystemService *CmdbSystemService) ParseInfoList2Excel(infoList []application.ApplicationServer, headers []string, filePath string) error {
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

func (cmdbSystemService *CmdbSystemService) ImportExcel2db(file multipart.File, header *multipart.FileHeader) error {
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

func (cmdbSystemService *CmdbSystemService) ExportTemplate() (*excelize.File, error) {
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
