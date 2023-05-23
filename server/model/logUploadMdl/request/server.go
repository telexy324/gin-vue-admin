package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
)

type ServerSearch struct {
	logUploadMdl.Server
	request.PageInfo
}
