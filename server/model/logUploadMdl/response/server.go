package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
)

type ServersResponse struct {
	Servers []logUploadMdl.Server `json:"servers"`
}

type ServerResponse struct {
	Server logUploadMdl.Server `json:"server"`
}

type SecretsResponse struct {
	Secrets []logUploadMdl.Secret `json:"secrets"`
}

type SecretResponse struct {
	Secret logUploadMdl.Secret `json:"secret"`
}
