package application

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
//@function: AddSystem
//@description: 添加系统
//@param: system model.ApplicationSystem
//@return: error

func (cmdbSystemService *CmdbSystemService) AddSystem(addSystemRequest request2.AddSystem) (err error) {
	if !errors.Is(global.GVA_DB.Where("name = ?", addSystemRequest.System.Name).First(&application.ApplicationSystem{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Create(&addSystemRequest.System).Error; txErr != nil {
			global.GVA_LOG.Error("添加系统失败", zap.Any("err", err))
			return txErr
		}
		if addSystemRequest.SystemAdmin != nil && len(addSystemRequest.SystemAdmin) > 0 {
			for _, admin := range addSystemRequest.SystemAdmin {
				admin.SystemId = int(addSystemRequest.System.ID)
				if err := global.GVA_DB.Create(&admin).Error; err != nil {
					global.GVA_LOG.Error("添加管理员失败", zap.Any("err", err))
				}
			}
		}
		return nil
	})
	return
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
//@param: system model.ApplicationSystem
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
			existAdmins := make([]application.ApplicationSystemAdmin, 0)
			existIds := make([]int, 0, len(existAdmins))
			if txErr = tx.Where("system_id = ?", addSystemRequest.System.ID).Find(&existAdmins).Error; txErr != nil {
				return txErr
			}
			for _, admin := range existAdmins {
				existIds = append(existIds, admin.AdminId)
			}
			requestIds := make([]int, 0, len(addSystemRequest.SystemAdmin))
			for _, admin := range addSystemRequest.SystemAdmin {
				requestIds = append(requestIds, admin.AdminId)
			}
			toAdd := utils.SubInt(requestIds, existIds)
			toDel := utils.SubInt(existIds, requestIds)
			for _, id := range toAdd {
				if txErr = tx.Where("id = ?", id).First(&application.Admin{}).Error; txErr != nil {
					return txErr
				}
				var sysAdmin application.ApplicationSystemAdmin
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.System.ID).Find(&sysAdmin).Error; errors.Is(txErr, gorm.ErrRecordNotFound) {
					if txErr = tx.Create(&application.ApplicationSystemAdmin{
						SystemId: int(addSystemRequest.System.ID),
						AdminId:  id,
					}).Error; txErr != nil {
						return txErr
					}
				} else if txErr == nil {
					if txErr = tx.Unscoped().Where("admin_id = ? and system_id = ?", id, addSystemRequest.System.ID).Find(&sysAdmin).Update("deleted_at", nil).Error; txErr != nil {
						return txErr
					}
				} else {
					return txErr
				}
			}
			for _, id := range toDel {
				admin := application.ApplicationSystemAdmin{}
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.System.ID).Find(&admin).Delete(&admin).Error; txErr != nil {
					return txErr
				}
			}
		}
		return nil
	})
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemById
//@description: 返回当前选中system
//@param: id float64
//@return: err error, system model.ApplicationSystem

func (cmdbSystemService *CmdbSystemService) GetSystemById(id float64) (err error, system application.ApplicationSystem, admins []application.Admin) {
	if err = global.GVA_DB.Where("id = ?", id).First(&system).Error; err != nil {
		return
	}
	var systemAdmins []application.ApplicationSystemAdmin
	if err = global.GVA_DB.Where("system_id = ?", id).Find(&systemAdmins).Error; err != nil {
		return
	}
	adminIds := make([]int, 0, len(systemAdmins))
	for _, sysAdmin := range systemAdmins {
		adminIds = append(adminIds, sysAdmin.AdminId)
	}
	if err = global.GVA_DB.Where("id in ?", adminIds).Find(&admins).Error; err != nil {
		return
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemList
//@description: 获取系统分页
//@return: err error, list interface{}, total int64

func (cmdbSystemService *CmdbSystemService) GetSystemList(info request2.SystemSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var systemList []application.ApplicationSystem
	db := global.GVA_DB.Model(&application.ApplicationSystem{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&systemList).Error

	systemInfoList := make([]applicationRes.ApplicationSystemResponse, 0, len(systemList))
	for _, system := range systemList {
		sysAdmins := make([]application.ApplicationSystemAdmin, 0)
		adminInfos := make([]application.Admin, 0)
		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&sysAdmins).Error; err != nil {
			return
		}
		for _, sysAdmin := range sysAdmins {
			var adminInfo = application.Admin{}
			if err = global.GVA_DB.Where("id = ?", sysAdmin.AdminId).First(&adminInfo).Error; err != nil {
				return
			} else {
				adminInfos = append(adminInfos, adminInfo)
			}
		}
		systemInfoList = append(systemInfoList, applicationRes.ApplicationSystemResponse{
			System: system,
			Admins: adminInfos,
		})
	}
	return err, systemInfoList, total
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddRelation
//@description: 添加联系
//@param: relation model.SystemRelation
//@return: error

func (cmdbSystemService *CmdbSystemService) AddRelation(relation application.SystemRelation) error {
	if !errors.Is(global.GVA_DB.Where("start_system_id = ? and end_system_id = ?", relation.StartSystemId, relation.EndSystemId).First(&application.SystemRelation{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复关系")
	}

	startSystem := application.ApplicationSystem{}
	endSystem := application.ApplicationSystem{}
	err := global.GVA_DB.Where("id = ?", relation.StartSystemId).First(&startSystem).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Where("id = ?", relation.EndSystemId).First(&endSystem).Error
	if err != nil {
		return err
	}
	return global.GVA_DB.Create(&relation).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: SystemRelations
//@description: 返回当前选中System的关系路径
//@param: id float64
//@return: err error, system model.SystemRelation

func (cmdbSystemService *CmdbSystemService) SystemRelations(id float64) (err error, relations []application.SystemRelation, nodes []application.Node) {
	system := application.ApplicationSystem{}
	err = global.GVA_DB.Where("id = ?", id).First(&system).Error
	if err != nil {
		return
	}
	mapNodes := make(map[int]bool)
	nodes = append(nodes, application.Node{
		Id:   int(system.ID),
		Name: system.Name,
	})
	mapNodes[int(system.ID)] = true
	relationOneSrc := make([]application.SystemRelation, 0)
	relationOneDest := make([]application.SystemRelation, 0)
	err = global.GVA_DB.Preload("EndSystem").Where("start_system_id = ?", id).Find(&relationOneSrc).Error
	err = global.GVA_DB.Preload("StartSystem").Where("end_system_id = ?", id).Find(&relationOneDest).Error
	if len(relationOneSrc) > 0 {
		relations = append(relations, relationOneSrc...)
		for _, relation := range relationOneSrc {
			if mapNodes[relation.EndSystemId] == false {
				nodes = append(nodes, application.Node{
					Id:   relation.EndSystemId,
					Name: relation.EndSystem.Name,
				})
				mapNodes[relation.EndSystemId] = true
			}
			relationTwoSrc := make([]application.SystemRelation, 0)
			relationTwoDest := make([]application.SystemRelation, 0)
			err = global.GVA_DB.Preload("EndSystem").Where("start_system_id = ?", relation.EndSystemId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartSystem").Where("end_system_id = ?", relation.EndSystemId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndSystemId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndSystemId,
							Name: relation.EndSystem.Name,
						})
						mapNodes[relation.EndSystemId] = true
					}
				}
			}
			if len(relationTwoDest) > 0 {
				relations = append(relations, relationTwoDest...)
				for _, relation := range relationTwoDest {
					if mapNodes[relation.StartSystemId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.StartSystemId,
							Name: relation.StartSystem.Name,
						})
						mapNodes[relation.StartSystemId] = true
					}
				}
			}
		}
	}
	if len(relationOneDest) > 0 {
		relations = append(relations, relationOneDest...)
		for _, relation := range relationOneDest {
			if mapNodes[relation.StartSystemId] == false {
				nodes = append(nodes, application.Node{
					Id:   relation.StartSystemId,
					Name: relation.StartSystem.Name,
				})
				mapNodes[relation.StartSystemId] = true
			}
			relationTwoSrc := make([]application.SystemRelation, 0)
			relationTwoDest := make([]application.SystemRelation, 0)
			err = global.GVA_DB.Preload("EndSystem").Where("start_system_id = ?", relation.StartSystemId).Find(&relationTwoSrc).Error
			err = global.GVA_DB.Preload("StartSystem").Where("end_system_id = ?", relation.StartSystemId).Find(&relationTwoDest).Error
			if len(relationTwoSrc) > 0 {
				relations = append(relations, relationTwoSrc...)
				for _, relation := range relationTwoSrc {
					if mapNodes[relation.EndSystemId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.EndSystemId,
							Name: relation.EndSystem.Name,
						})
						mapNodes[relation.EndSystemId] = true
					}
				}
			}
			if len(relationTwoDest) > 0 {
				relations = append(relations, relationTwoDest...)
				for _, relation := range relationTwoDest {
					if mapNodes[relation.StartSystemId] == false {
						nodes = append(nodes, application.Node{
							Id:   relation.StartSystemId,
							Name: relation.StartSystem.Name,
						})
						mapNodes[relation.StartSystemId] = true
					}
				}
			}
		}
	}
	return
}
