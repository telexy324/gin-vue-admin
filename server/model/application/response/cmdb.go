package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
)

type ApplicationServerResponse struct {
	Servers []application.ApplicationServer `json:"servers"`
}
