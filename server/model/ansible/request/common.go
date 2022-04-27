package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// Find by id and project_id structure
type GetByProjectId struct {
	ID        float64 `json:"id" form:"id"`                 // 主键ID
	ProjectId float64 `json:"project_id" form:"project_id"` // 主键ID
	request.PageInfo
	SortInverted bool
	SortBy       string
}
