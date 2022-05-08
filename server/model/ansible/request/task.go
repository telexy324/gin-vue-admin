package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// Find by id and project_id structure
type AddTaskByProjectId struct {
	ProjectId float64      `json:"project_id" form:"project_id"`
	ansible.Task
}

// Find by id and project_id structure
type GetTaskByTemplateId struct {
	ID         float64 `json:"id" form:"id"` // 主键ID
	ProjectId  float64 `json:"project_id" form:"project_id"`
	TemplateId float64 `json:"template_id" form:"template_id"`
	request.PageInfo
	SortInverted bool
	SortBy       string
}
