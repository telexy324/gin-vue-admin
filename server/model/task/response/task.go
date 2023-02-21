package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/task"
)

type TasksResponse struct {
	Tasks []task.Task `json:"tasks"`
}

type TaskResponse struct {
	Task task.Task `json:"task"`
}

type TaskOutputsResponse struct {
	TaskOutputs []task.TaskOutput `json:"taskOutputs"`
}
