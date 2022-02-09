package utils

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

func convertStruct(a interface{}, b interface{}) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	return nil
}

func ConvertStruct(a interface{}, b interface{}) error {
	err := convertStruct(a, b)
	if err != nil {
		global.GVA_LOG.Error("转换失败", zap.Any("err", err))
	}
	return err
}
