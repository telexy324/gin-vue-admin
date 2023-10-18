package application

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	applicationRes "github.com/flipped-aurora/gin-vue-admin/server/model/application/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
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
	if !errors.Is(global.GVA_DB.Where("name = ?", addSystemRequest.Name).First(&application.ApplicationSystem{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Create(&addSystemRequest.ApplicationSystem).Error; txErr != nil {
			global.GVA_LOG.Error("添加系统失败", zap.Any("err", err))
			return txErr
		}
		if addSystemRequest.SystemAdmin != nil && len(addSystemRequest.SystemAdmin) > 0 {
			for _, admin := range addSystemRequest.SystemAdmin {
				admin.SystemId = int(addSystemRequest.ApplicationSystem.ID)
				if err = global.GVA_DB.Create(&admin).Error; err != nil {
					global.GVA_LOG.Error("添加管理员失败", zap.Any("err", err))
				}
			}
		}
		if addSystemRequest.AdminIds != nil && len(addSystemRequest.AdminIds) > 0 {
			for _, id := range addSystemRequest.AdminIds {
				user := &system.SysUser{}
				user.ID = uint(id)
				if err = global.GVA_DB.Find(user).Error; err != nil {
					global.GVA_LOG.Error("管理员不存在", zap.Any("err", err))
					continue
				}
				admin := &application.ApplicationSystemSysAdmin{
					SystemId: int(addSystemRequest.ApplicationSystem.ID),
					AdminId:  id,
				}
				if err = global.GVA_DB.Create(&admin).Error; err != nil {
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
	var existSystem application.ApplicationSystem
	if err = global.GVA_DB.Where("id = ?", id).First(&existSystem).Delete(&existSystem).Error; err != nil {
		return err
	}
	var systemAdmins []application.ApplicationSystemAdmin
	err = global.GVA_DB.Where("system_id = ?", id).Find(&systemAdmins).Delete(&systemAdmins).Error
	if err != nil {
		return err
	}
	var admins []application.ApplicationSystemSysAdmin
	err = global.GVA_DB.Where("system_id = ?", id).Find(&admins).Delete(&admins).Error
	if err != nil {
		return err
	}
	var editRelation application.ApplicationSystemEditRelation
	err = global.GVA_DB.Where("system_id = ?", id).First(&editRelation).Delete(&editRelation).Error
	if err != nil {
		return err
	}
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSystemByIds
//@description: 批量删除系统
//@param: applicationSystems []model.applicationSystems
//@return: err error

func (cmdbSystemService *CmdbSystemService) DeleteSystemByIds(ids request.IdsReq) (err error) {
	if ids.Ids == nil || len(ids.Ids) <= 0 {
		return
	}
	for _, id := range ids.Ids {
		if err = cmdbSystemService.DeleteSystem(float64(id)); err != nil {
			return
		}
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateSystem
//@description: 更新系统
//@param: system model.ApplicationSystem
//@return: err error

func (cmdbSystemService *CmdbSystemService) UpdateSystem(addSystemRequest request2.AddSystem) (err error) {
	var oldSystem application.ApplicationSystem
	upDateMap := make(map[string]interface{})
	userJson, err := json.Marshal(addSystemRequest.SshUsers)
	if err != nil {
		return err
	}
	upDateMap["name"] = addSystemRequest.Name
	upDateMap["position"] = addSystemRequest.Position
	upDateMap["network"] = addSystemRequest.Network
	upDateMap["ssh_user"] = userJson

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", addSystemRequest.ID).Find(&oldSystem)
		if oldSystem.Name != addSystemRequest.Name {
			if err = tx.Where("id <> ? AND name = ?", addSystemRequest.ID, addSystemRequest.Name).First(&application.ApplicationSystem{}).Error; err != nil && err != gorm.ErrRecordNotFound {
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
			if txErr = tx.Where("system_id = ?", addSystemRequest.ID).Find(&existAdmins).Error; txErr != nil {
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
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&sysAdmin).Error; txErr != nil {
					return txErr
				}
				if sysAdmin.ID > 0 {
					if txErr = tx.Unscoped().Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&sysAdmin).Update("deleted_at", nil).Error; txErr != nil {
						return txErr
					}
				} else {
					if txErr = tx.Create(&application.ApplicationSystemAdmin{
						SystemId: int(addSystemRequest.ID),
						AdminId:  id,
					}).Error; txErr != nil {
						return txErr
					}
				}
			}
			for _, id := range toDel {
				admin := application.ApplicationSystemAdmin{}
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&admin).Delete(&admin).Error; txErr != nil {
					return txErr
				}
			}
		}
		if addSystemRequest.AdminIds != nil && len(addSystemRequest.AdminIds) > 0 {
			existAdmins := make([]application.ApplicationSystemSysAdmin, 0)
			existIds := make([]int, 0, len(existAdmins))
			if txErr = tx.Where("system_id = ?", addSystemRequest.ID).Find(&existAdmins).Error; txErr != nil {
				return txErr
			}
			for _, admin := range existAdmins {
				existIds = append(existIds, admin.AdminId)
			}
			requestIds := make([]int, 0, len(addSystemRequest.AdminIds))
			for _, id := range addSystemRequest.AdminIds {
				requestIds = append(requestIds, id)
			}
			toAdd := utils.SubInt(requestIds, existIds)
			toDel := utils.SubInt(existIds, requestIds)
			for _, id := range toAdd {
				if txErr = tx.Where("id = ?", id).First(&system.SysUser{}).Error; txErr != nil {
					return txErr
				}
				var sysAdmin application.ApplicationSystemSysAdmin
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&sysAdmin).Error; txErr != nil {
					return txErr
				}
				if sysAdmin.ID > 0 {
					if txErr = tx.Unscoped().Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&sysAdmin).Update("deleted_at", nil).Error; txErr != nil {
						return txErr
					}
				} else {
					if txErr = tx.Create(&application.ApplicationSystemSysAdmin{
						SystemId: int(addSystemRequest.ID),
						AdminId:  id,
					}).Error; txErr != nil {
						return txErr
					}
				}
			}
			for _, id := range toDel {
				admin := application.ApplicationSystemSysAdmin{}
				if txErr = tx.Where("admin_id = ? and system_id = ?", id, addSystemRequest.ID).Find(&admin).Delete(&admin).Error; txErr != nil {
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

func (cmdbSystemService *CmdbSystemService) GetSystemById(id float64) (err error, system application.ApplicationSystem, admins []application.Admin, adminIds []int) {
	if err = global.GVA_DB.Where("id = ?", id).First(&system).Error; err != nil {
		return
	}
	var systemAdmins []application.ApplicationSystemAdmin
	if err = global.GVA_DB.Where("system_id = ?", id).Find(&systemAdmins).Error; err != nil {
		return
	}
	sysAdminIds := make([]int, 0, len(systemAdmins))
	for _, sysAdmin := range systemAdmins {
		sysAdminIds = append(sysAdminIds, sysAdmin.AdminId)
	}
	if err = global.GVA_DB.Where("id in ?", sysAdminIds).Find(&admins).Error; err != nil {
		return
	}
	ads := make([]application.ApplicationSystemSysAdmin, 0)
	adminIds = make([]int, 0)
	if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&ads).Error; err != nil {
		return
	}
	for _, admin := range ads {
		adminIds = append(adminIds, admin.AdminId)
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemList
//@description: 获取系统分页
//@return: err error, list interface{}, total int64

func (cmdbSystemService *CmdbSystemService) GetSystemList(info request2.SystemSearch, adminID uint) (err error, list interface{}, total int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var systemList []application.ApplicationSystem
	//db := global.GVA_DB.Model(&application.ApplicationSystem{})
	//if info.Name != "" {
	//	name := strings.Trim(info.Name, " ")
	//	db = db.Where("`name` LIKE ?", "%"+name+"%")
	//}
	sqlRaw := `select a.* from application_systems a, application_system_sys_admins ass where a.id=ass.system_id and ass.admin_id = ?`
	db := global.GVA_DB
	if info.Name != "" {
		sqlRaw = sqlRaw + ` and a.name LIKE ?`
		db = db.Raw(sqlRaw, adminID, "%"+strings.Trim(info.Name, " ")+"%")
	} else {
		db = db.Raw(sqlRaw, adminID)
	}
	//err = db.Count(&total).Error
	//if err != nil {
	//	return
	//}
	err = db.Limit(limit).Offset(offset).Scan(&systemList).Error
	if err != nil {
		return
	}

	systemInfoList := make([]applicationRes.ApplicationSystemResponse, 0, len(systemList))
	for _, system := range systemList {
		sysAdmins := make([]application.ApplicationSystemAdmin, 0)
		adminInfos := make([]application.Admin, 0)
		admins := make([]application.ApplicationSystemSysAdmin, 0)
		adminIds := make([]int, 0)
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
		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&admins).Error; err != nil {
			return
		}
		for _, admin := range admins {
			adminIds = append(adminIds, int(admin.ID))
		}
		systemInfoList = append(systemInfoList, applicationRes.ApplicationSystemResponse{
			System:   system,
			Admins:   adminInfos,
			AdminIds: adminIds,
		})
	}
	return err, systemInfoList, len(systemList)
}

//func (cmdbSystemService *CmdbSystemService) GetSystemList(info request2.SystemSearch, adminID uint) (err error, list interface{}, total int64) {
//	limit := info.PageSize
//	offset := info.PageSize * (info.Page - 1)
//	var systemList []application.ApplicationSystem
//	db := global.GVA_DB.Model(&application.ApplicationSystem{})
//	if info.Name != "" {
//		name := strings.Trim(info.Name, " ")
//		db = db.Where("`name` LIKE ?", "%"+name+"%")
//	}
//	err = db.Count(&total).Error
//	if err != nil {
//		return
//	}
//	err = db.Limit(limit).Offset(offset).Find(&systemList).Error
//
//	systemInfoList := make([]applicationRes.ApplicationSystemResponse, 0, len(systemList))
//	for _, system := range systemList {
//		sysAdmins := make([]application.ApplicationSystemAdmin, 0)
//		adminInfos := make([]application.Admin, 0)
//		admins := make([]application.ApplicationSystemSysAdmin, 0)
//		adminIds := make([]int, 0)
//		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&sysAdmins).Error; err != nil {
//			return
//		}
//		for _, sysAdmin := range sysAdmins {
//			var adminInfo = application.Admin{}
//			if err = global.GVA_DB.Where("id = ?", sysAdmin.AdminId).First(&adminInfo).Error; err != nil {
//				return
//			} else {
//				adminInfos = append(adminInfos, adminInfo)
//			}
//		}
//		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&admins).Error; err != nil {
//			return
//		}
//		var hasSystem bool
//		for _, admin := range admins {
//			if admin.AdminId == int(adminID) {
//				hasSystem = true
//			}
//			adminIds = append(adminIds, int(admin.ID))
//		}
//		if !hasSystem {
//			continue
//		}
//		systemInfoList = append(systemInfoList, applicationRes.ApplicationSystemResponse{
//			System:   system,
//			Admins:   adminInfos,
//			AdminIds: adminIds,
//		})
//	}
//	return err, systemInfoList, int64(len(systemInfoList))
//}

// @author: [telexy324](https://github.com/telexy324)
// @function: GetSystemServers
// @description: 获取系统内全部服务器
// @return: err error, systemList []application.ApplicationSystem
func (cmdbSystemService *CmdbSystemService) GetAdminSystems(adminID uint) (err error, systemList []application.ApplicationSystem) {
	sysAdmins := make([]application.ApplicationSystemSysAdmin, 0)
	if err = global.GVA_DB.Where("admin_id = ?", adminID).Find(&sysAdmins).Error; err != nil {
		return
	}
	for _, admin := range sysAdmins {
		system := application.ApplicationSystem{}
		if err = global.GVA_DB.Model(&application.ApplicationSystem{}).Where("id = ?", admin.SystemId).Find(&system).Error; err != nil {
			return
		}
		systemList = append(systemList, system)
	}
	return
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

//@author: [telexy324](https://github.com/telexy324)
//@function: AddEditRelations
//@description: 添加编辑器关系图
//@param: relation model.ApplicationSystemEditRelation
//@return: error

func (cmdbSystemService *CmdbSystemService) AddEditRelations(relation application.ApplicationSystemEditRelation) error {
	return global.GVA_DB.Create(&relation).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteEditRelations
//@description: 删除编辑器关系图
//@param: id float64
//@return: err error

func (cmdbSystemService *CmdbSystemService) DeleteEditRelations(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.ApplicationSystemEditRelation{}).Error
	if err != nil {
		return
	}
	var relation application.ApplicationSystemEditRelation
	return global.GVA_DB.Where("id = ?", id).First(&relation).Delete(&relation).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateEditRelations
//@description: 更新编辑器关系图
//@param: relation model.ApplicationSystemEditRelation
//@return: err error

func (cmdbSystemService *CmdbSystemService) UpdateEditRelations(relation application.ApplicationSystemEditRelation) (err error) {
	upDateMap := make(map[string]interface{})
	upDateMap["system_id"] = relation.SystemId
	upDateMap["relation"] = relation.Relation
	err = global.GVA_DB.Model(&relation).Transaction(func(tx *gorm.DB) error {
		txErr := tx.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetEditRelations
//@description: 返回编辑器关系图
//@param: id float64
//@return: err error, server model.ApplicationSystemEditRelation

func (cmdbSystemService *CmdbSystemService) GetSystemEditRelations(id float64) (err error, relation application.ApplicationSystemEditRelation) {
	err = global.GVA_DB.Where("system_id = ?", id).First(&relation).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}
