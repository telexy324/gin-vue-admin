package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
)

// Add Server
type AddServer struct {
	Server      *application.ApplicationServer `json:"server"`
	AuthorityId string                         `json:"authorityId"` // 角色ID
}
