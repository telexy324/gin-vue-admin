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

type TaskDashboardInfo struct {
	Date   string `json:"date"`
	Number int64  `json:"number"`
}

type TaskDashboardResponse struct {
	TaskDashboardInfos []TaskDashboardInfo `json:"taskDashboardInfos"`
}

type StopTaskResponse struct {
	FailedIps []string `json:"failedIps"`
}