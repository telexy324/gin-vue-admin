package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/application"

type GetToken struct {
	Servers []application.ApplicationServer `json:"servers" form:"servers"` // 主键ID
	Client  string                          `json:"client" form:"client"`
}

type Token struct {
	Token string `json:"token"` // token
}
