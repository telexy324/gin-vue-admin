package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
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

type TaskTemplateSearch struct {
	taskMdl.TaskTemplate
	request.PageInfo
	SystemIDs []int
}

type AddTaskTemplate struct {
	TaskTemplate *taskMdl.TaskTemplate `json:"taskTemplate"`
	ServerIds    []int                 `json:"serverIds"`
	AuthorityId  string                `json:"authorityId"` // 角色ID
}

type UpdateTaskTemplate struct {
	TaskTemplate *taskMdl.TaskTemplate `json:"taskTemplate"`
	ServerIds    []int                 `json:"serverIds"`
	AuthorityId  string                `json:"authorityId"` // 角色ID
}

type TaskTemplateSetSearch struct {
	taskMdl.TaskTemplateSet
	request.PageInfo
}

type SetTaskSearch struct {
	taskMdl.SetTask
	request.PageInfo
}