package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type EnvironmentsResponse struct {
	Environments []ansible.Environment `json:"environments"`
}

type EnvironmentResponse struct {
	Environment ansible.Environment `json:"environment"`
}
