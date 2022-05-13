package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type UsersResponse struct {
	Users []ansible.ProjectUser `json:"users"`
}

type UserResponse struct {
	User ansible.ProjectUser `json:"user"`
}
