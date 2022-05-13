package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type ProjectsResponse struct {
	Projects []ansible.Project `json:"projects"`
}

type ProjectResponse struct {
	Project ansible.Project `json:"project"`
}
