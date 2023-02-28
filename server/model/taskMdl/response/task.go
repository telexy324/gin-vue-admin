package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
)

type TasksResponse struct {
	Tasks []taskMdl.Task `json:"tasks"`
}

type TaskResponse struct {
	Task taskMdl.Task `json:"task"`
}

type TaskOutputsResponse struct {
	TaskOutputs []taskMdl.TaskOutput `json:"taskOutputs"`
}
