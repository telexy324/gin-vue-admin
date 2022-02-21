package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AddAdmin struct {
	Server      *application.Admin `json:"admin"`
	AuthorityId string                         `json:"authorityId"` // 角色ID
}

type AdminSearch struct {
	application.Admin
	request.PageInfo
}