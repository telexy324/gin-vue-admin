package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AddSystem struct {
	System      *application.ApplicationSystem       `json:"system"`
	SystemAdmin []application.ApplicationSystemAdmin `json:"systemAdmin"`
	AuthorityId string                               `json:"authorityId"` // 角色ID
}

type SystemSearch struct {
	application.ApplicationServer
	request.PageInfo
}

type ExcelInfoSystem struct {
	FileName string                          `json:"fileName"` // 文件名
	InfoList []application.ApplicationServer `json:"infoList"`
	Header   []string                        `json:"header"`
}
