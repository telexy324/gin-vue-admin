package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

//// Find by id and project_id structure
//type AddTaskByProjectId struct {
//	ProjectId float64 `json:"projectId" form:"projectId"`
//	task.Task
//}

// Find by id and project_id structure
type GetTaskByTemplateId struct {
	ID float64 `json:"id" form:"id"` // 主键ID
	// ProjectId  float64 `json:"projectId" form:"projectId"`
	TemplateId float64 `json:"templateId" form:"templateId"`
	request.PageInfo
	SortInverted bool
	SortBy       string
}

// Find by id and project_id structure
type GetTaskOutputsByTaskId struct {
	ID           float64 `json:"id" form:"id"` // 主键ID
	TaskId       float64 `json:"taskId" form:"taskId"`
	SortInverted bool
	SortBy       string
}