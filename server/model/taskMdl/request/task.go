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
	OrderKey     string `json:"orderKey"` // 排序
	Desc         bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
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
	OrderKey  string `json:"orderKey"` // 排序
	Desc      bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
	SystemIDs []int  `json:"systemIds"`
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
	OrderKey  string `json:"orderKey"` // 排序
	Desc      bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
	SystemIDs []int  `json:"systemIds"`
}

type SetTaskSearch struct {
	taskMdl.SetTask
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type GetTaskBySetTaskIdWithSeq struct {
	SetTaskId    float64 `json:"setTaskId" form:"setTaskId"`
	CurrentSeq   int     `json:"currentSeq" form:"currentSeq"`
	CurrentIndex int     `json:"currentIndex" form:"currentIndex"`
	Redo         bool    `json:"redo" form:"redo"`
	request.PageInfo
}
