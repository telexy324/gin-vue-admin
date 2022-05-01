package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type KeysResponse struct {
	Keys []ansible.AccessKey `json:"keys"`
}

type KeyResponse struct {
	Key ansible.AccessKey `json:"key"`
}
