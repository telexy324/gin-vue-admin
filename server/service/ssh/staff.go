package ssh

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"gorm.io/gorm"
	"strings"
)

type SshService struct {
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddAdmin
//@description: 添加管理员
//@param: server model.Admin
//@return: error

func (staffService *SshService) AddAdmin(admin application.Admin) error {
	if !errors.Is(global.GVA_DB.Where("name = ? or mobile = ?", admin.Name, admin.Mobile).First(&application.Admin{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复姓名或电话，请修改")
	}
	return global.GVA_DB.Create(&admin).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteAdmin
//@description: 删除管理员
//@param: id float64
//@return: err error

func (staffService *SshService) DeleteAdmin(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application.Admin{}).Error
	if err != nil {
		return
	}
	var admin application.Admin
	return global.GVA_DB.Where("id = ?", id).First(&admin).Delete(&admin).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateAdmin
//@description: 更新管理员
//@param: server model.Admin
//@return: err error

func (staffService *SshService) UpdateAdmin(admin application.Admin) (err error) {
	var oldAdmin application.Admin
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = admin.Name
	upDateMap["mobile"] = admin.Mobile
	upDateMap["department_id"] = admin.DepartmentId
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", admin.ID).Find(&oldAdmin)
		if oldAdmin.Name != admin.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", admin.ID, admin.Name).First(&application.Admin{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
		if oldAdmin.Mobile != admin.Mobile {
			if !errors.Is(tx.Where("id <> ? AND mobile = ?", admin.ID, admin.Mobile).First(&application.Admin{}).Error, gorm.ErrRecordNotFound) {
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
//@function: GetAdminById
//@description: 返回当前选中admin
//@param: id float64
//@return: err error, server model.Admin

func (staffService *SshService) GetAdminById(id float64) (err error, admin application.Admin) {
	err = global.GVA_DB.Where("id = ?", id).First(&admin).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetAdminList
//@description: 获取管理员分页
//@return: err error, list interface{}, total int64

func (staffService *SshService) GetAdminList(info request2.AdminSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var adminList []application.Admin
	db := global.GVA_DB.Model(&application.Admin{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	if info.Mobile != "" {
		mobile := strings.Trim(info.Mobile, " ")
		db = db.Where("`mobile` LIKE ?", "%"+mobile+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&adminList).Error
	return err, adminList, total
}

// @author: [telexy324](https://github.com/telexy324)
// @function: GetDepartmentAll
// @description: 获取全部部门
// @return: err error, departmentList []application.Department
func (staffService *SshService) GetDepartmentAll() (err error, departmentList []application.Department) {
	db := global.GVA_DB.Model(&application.Department{})
	err = db.Find(&departmentList).Error
	return
}
