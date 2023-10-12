package application

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

//@author: [granty1](https://github.com/granty1)
//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

type ApplicationRecordService struct {
}

func (applicationRecordService *ApplicationRecordService) CreateApplicationRecord(applicationRecord system.SysOperationRecord) (err error) {
	err = global.GVA_DB.Create(&applicationRecord).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids request.IdsReq
//@return: err error

func (applicationRecordService *ApplicationRecordService) DeleteApplicationRecordByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func (applicationRecordService *ApplicationRecordService) DeleteApplicationRecord(applicationRecord system.SysOperationRecord) (err error) {
	err = global.GVA_DB.Delete(&applicationRecord).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@function: DeleteSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: err error, sysOperationRecord model.SysOperationRecord

func (applicationRecordService *ApplicationRecordService) GetApplicationRecord(id uint) (err error, applicationRecord system.SysOperationRecord) {
	err = global.GVA_DB.Where("id = ?", id).First(&applicationRecord).Error
	return
}

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info systemReq.SysOperationRecordSearch
//@return: err error, list interface{}, total int64

func (applicationRecordService *ApplicationRecordService) GetApplicationRecordInfoList(info systemReq.SysOperationRecordSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&system.SysOperationRecord{})
	var sysOperationRecords []system.SysOperationRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return err, sysOperationRecords, total
}
