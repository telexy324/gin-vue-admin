package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type TasksResponse struct {
	Tasks []ansible.Task `json:"tasks"`
}

type TaskResponse struct {
	Task ansible.Task `json:"task"`
}

type TaskOutputsResponse struct {
	TaskOutputs []ansible.TaskOutput `json:"taskOutputs"`
}
