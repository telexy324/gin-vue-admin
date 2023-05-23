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

type SecretService struct {
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddSecret
//@description: 添加密钥
//@param: secret model.Secret
//@return: error

func (secretService *SecretService) AddSecret(secret logUploadMdl.Secret) error {
	if !errors.Is(global.GVA_DB.Where("name = ?", secret.Name).First(&logUploadMdl.Secret{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.GVA_DB.Create(&secret).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSecret
//@description: 删除密钥
//@param: id float64
//@return: err error

func (secretService *SecretService) DeleteSecret(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&logUploadMdl.Secret{}).Error
	if err != nil {
		return
	}
	var secret logUploadMdl.Secret
	return global.GVA_DB.Where("id = ?", id).First(&secret).Delete(&secret).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSecretByIds
//@description: 批量删除密钥
//@param: secrets []model.Secrets
//@return: err error

func (secretService *SecretService) DeleteSecretByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]logUploadMdl.Secret{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateSecret
//@description: 更新路由
//@param: secret model.Secret
//@return: err error

func (secretService *SecretService) UpdateSecret(secret logUploadMdl.Secret) (err error) {
	var oldSecret logUploadMdl.Secret
	upDateMap := make(map[string]interface{})
	upDateMap["server_id"] = secret.ServerId
	upDateMap["name"] = secret.Name
	upDateMap["password"] = secret.Password

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", secret.ID).Find(&oldSecret)
		if oldSecret.Name != secret.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", secret.ID, secret.Name).First(&logUploadMdl.Secret{}).Error, gorm.ErrRecordNotFound) {
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
//@function: GetSecretById
//@description: 返回当前选中secret
//@param: id float64
//@return: err error, secret model.Secret

func (secretService *SecretService) GetSecretById(id float64) (err error, secret logUploadMdl.Secret) {
	err = global.GVA_DB.Where("id = ?", id).First(&secret).Error
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSecretList
//@description: 获取密钥分页
//@return: err error, list interface{}, total int64

func (secretService *SecretService) GetSecretList(info request2.SecretSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var secretList []logUploadMdl.Secret
	db := global.GVA_DB.Model(&logUploadMdl.Secret{})
	if info.Name != "" {
		hostname := strings.Trim(info.Name, " ")
		db = db.Where("`hostname` LIKE ?", "%"+hostname+"%")
	}
	if info.ServerId > 0 {
		db = db.Where("server_id = ?", info.ServerId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&secretList).Error
	return err, secretList, total
}
