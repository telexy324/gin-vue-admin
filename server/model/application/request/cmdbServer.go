package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// Add Server
type AddServer struct {
	Server      *application.ApplicationServer `json:"server"`
	AuthorityId string                         `json:"authorityId"` // 角色ID
}

type ServerSearch struct {
	application.ApplicationServer
	request.PageInfo
}

type ExcelInfo struct {
	FileName string                          `json:"fileName"` // 文件名
	InfoList []application.ApplicationServer `json:"infoList"`
	Header   []string                        `json:"header"`
}
