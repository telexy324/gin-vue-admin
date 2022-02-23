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
			existAdmins := make([]application.Admin, 0)
			existIds := make([]int, 0, len(existAdmins))
			if err = db.Where("system_id = ?", addSystemRequest.System.ID).First(&existAdmins).Error; err != nil {
				return err
			}
			for _, admin := range existAdmins {
				existIds = append(existIds, int(admin.ID))
			}
			requestIds := make([]int, 0, len(addSystemRequest.SystemAdmin))
			for _, admin := range addSystemRequest.SystemAdmin {
				requestIds = append(requestIds, admin.AdminId)
			}
			toAdd := utils.SubInt(requestIds, existIds)
			toDel := utils.SubInt(existIds, requestIds)
			for _, id := range toAdd {
				if err = db.Where("id = ?", id).First(&application.Admin{}).Error; err != nil {
					return err
				}
				if err = db.Create(application.ApplicationSystemAdmin{
					SystemId: int(addSystemRequest.System.ID),
					AdminId:  id,
				}).Error; err != nil {
					return err
				}
			}
			for _, id := range toDel {
				admin := application.Admin{}
				if err = db.Where("id = ?", id).Find(&admin).Delete(&admin).Error; err != nil {
					return err
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
//@return: err error, server model.ApplicationSystem

func (cmdbSystemService *CmdbSystemService) GetSystemById(id float64) (err error, system application.ApplicationSystem, admins []application.Admin) {
	if err = global.GVA_DB.Where("id = ?", id).First(&system).Error; err != nil {
		return
	}
	if err = global.GVA_DB.Where("system_id = ?", id).Find(&admins).Error; err != nil {
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
		admins := make([]application.ApplicationSystemAdmin, 0)
		adminInfos := make([]application.Admin, 0)
		if err = global.GVA_DB.Where("system_id = ?", system.ID).Find(&admins).Error; err != nil {
			return
		}
		for _, admin := range admins {
			var adminInfo = application.Admin{}
			if err = global.GVA_DB.Where("id = ?", admin.ID).First(&adminInfo).Error; err != nil {
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
	return err, systemList, total
}
