package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AddServer struct {
	Server      *application.ApplicationServer `json:"server"`
	Apps        []int                          `json:"apps"`
	AuthorityId string                         `json:"authorityId"` // 角色ID
}

type UpdateServer struct {
	Server      *application.ApplicationServer `json:"server"`
	Apps        []int                          `json:"apps"`
	AuthorityId string                         `json:"authorityId"` // 角色ID
}

type ServerSearch struct {
	application.ApplicationServer
	request.PageInfo
	SystemIDs []int
}

type ExcelInfo struct {
	FileName string                          `json:"fileName"` // 文件名
	InfoList []application.ApplicationServer `json:"infoList"`
	Header   []string                        `json:"header"`
}

type AddApp struct {
	App         *application.App `json:"app"`
	AuthorityId string           `json:"authorityId"` // 角色ID
}

type AppSearch struct {
	application.App
	request.PageInfo
}
